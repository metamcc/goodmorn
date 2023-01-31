package client

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"github.com/mycreditchain/kafka-agent/kafka-consumer-agent/util"
)

const (
	SEPARATOR     = "!"
	SEED_GIFT_SVC = 2
	MESSAGE_TOPIC = "message-topic"
)

type KafkaAgentClient struct {
	AgentConfig map[string]interface{}
	SvcList     map[string]*Svc
	Producer    sarama.SyncProducer
	//ConsumerGroup    sarama.ConsumerGroup
}

type Svc struct {
	MessageTimer     *time.Timer
	RedisClient      *redis.Client
	RestfulApiMethod string
	RestfulApiUrl    string
	TimerPipe        chan struct{}
}

func (ka *KafkaAgentClient) Setup(_ sarama.ConsumerGroupSession) error {
	ka.SvcList = make(map[string]*Svc)
	MsgLogger.Println("beginning", time.Now())

	//defer ka.MessageTimer.Stop()
	for svc, svcConfig := range ka.AgentConfig["batch_svc_list"].(map[string]interface{}) {
		createTimer(ka, svc, svcConfig.(map[string]interface{}))
	}

	// db에 미처 전송 못한 메시지가 있으면 리스트를 돌면서 처리
	for svc, svcConfig := range ka.AgentConfig["batch_svc_list"].(map[string]interface{}) {
		svcNum, err := strconv.ParseFloat(svc, 64)
		if err != nil {
			MsgLogger.Println(svc, "Failed to convert svc string to number", err)
		}
		req := map[string]interface{}{
			"svc":                     svcNum,
			PRODUCER_FIELD_CNT_INFO:   map[string]interface{}{},
			MESSAGE_TOPIC:             "batch-resp",
			REST_REQUEST_FIELD_METHOD: svcConfig.(map[string]interface{})["btcMethod"].(string),
			REST_REQUEST_FIELD_URL:    svcConfig.(map[string]interface{})["btcUrl"].(string),
		} // 해당 svc 종류에 따라 req 형식 변경 작업 필요
		MsgLogger.Println("setup svc config", req)
		iter := ka.SvcList[svc].RedisClient.Scan(0, "svc_"+svc+"_*", ka.SvcList[svc].RedisClient.DBSize().Val()).Iterator()
		ka.flushMessages(req, iter)
	}
	return nil
}

func (*KafkaAgentClient) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim: 들어오는 메시지마다 go routine을 생성하여
// 메시지를 파싱하고 restful을 호출하여 응답을 받아
// 카프카에 응답 또는 재시도 토픽에 맞추어 입력하는 함수
func (ka *KafkaAgentClient) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		MsgLogger.Printf("Message topic:%q partition:%d offset:%d\n", claim.Topic(), claim.Partition(), message.Offset)
		sess.MarkMessage(message, "")

		// 1. 해당 토픽에 들어온 메시지 파싱 단계
		reqData, err := parseMessage(message.Value)
		if err != nil {
			MsgLogger.Println("Failed to parse queue message", err, reqData)
			// 파싱 에러시 카프카에 에러로 카프카 브로커에 응답 입력
			continue
		}

		svc := fmt.Sprintf("%d", util.TypeToInt(reqData["svc"]))
		batchSvcConfig, ok := ka.AgentConfig["batch_svc_list"].(map[string]interface{})[svc]
		go func() {
			reqData[MESSAGE_TOPIC] = message.Topic
			if strings.Contains(message.Topic, "batch-req") && ok {
				ka.ManageMultipleMessages(reqData, batchSvcConfig.(map[string]interface{}))
			} else if strings.Contains(message.Topic, "config") {
				ka.ManageConfig(reqData)
			} else {
				ka.ManageMessage(reqData)
			}
		}()
		if ok {
			ka.SvcList[svc].TimerPipe <- struct{}{}
		}
	}
	MsgLogger.Printf("Processed all messages")
	return nil
}

func (ka *KafkaAgentClient) ManageMultipleMessages(req, batchSvcConfig map[string]interface{}) {
	svc := fmt.Sprintf("%d", util.TypeToInt(req["svc"]))
	reqJson, err := json.Marshal(req)
	if err != nil {
		return
	}
	msgCount := int(ka.SvcList[svc].RedisClient.DBSize().Val())
	err = ka.SvcList[svc].RedisClient.Set(fmt.Sprintf("svc_%s_%d", svc, msgCount), string(reqJson), 0).Err()
	if err != nil {
		MsgLogger.Printf("Seting key, value in redis DB is failed")
		return
	}
	msgCount += 1
	if msgCount == util.TypeToInt(batchSvcConfig["btcMCnt"]) {
		iter := ka.SvcList[svc].RedisClient.Scan(0, "svc_"+svc+"_*", ka.SvcList[svc].RedisClient.DBSize().Val()).Iterator()
		ka.flushMessages(req, iter)
	}
	if msgCount == 1 || msgCount == util.TypeToInt(batchSvcConfig["btcMCnt"]) {
		ka.SvcList[svc].MessageTimer.Reset(time.Second * time.Duration(util.TypeToInt(batchSvcConfig["btcCycle"])))
		MsgLogger.Println("***", svc, "Timer is reseted by messaging term", time.Now())
	}
	<-ka.SvcList[svc].TimerPipe
}

func (ka *KafkaAgentClient) flushMessages(req map[string]interface{}, iter *redis.ScanIterator) {
	svc := fmt.Sprintf("%d", util.TypeToInt(req["svc"]))
	var keys, msgs []string
	for iter.Next() {
		keys = append(keys, iter.Val())
	}
	sort.Strings(keys)
	for _, key := range keys {
		val, err := ka.SvcList[svc].RedisClient.Get(key).Result()
		if err != nil {
			MsgLogger.Println("redis get error:", err)
		}
		msgs = append(msgs, val)
	}
	_, err := ka.SvcList[fmt.Sprintf("%d", util.TypeToInt(req["svc"]))].RedisClient.FlushDB().Result()
	if err != nil {
		MsgLogger.Println("redis flush err:", err)
	}
	if len(msgs) == 0 {
		MsgLogger.Println("no messages to create multiple message")
		return
	}
	reqData, err := createMultipleRequest(msgs)
	if err != nil {
		// telegram message
		MsgLogger.Println("creating multiple message is failed:", err)
		return
	}
	req["id"] = fmt.Sprintf("%v", time.Now().UnixNano())
	req[REST_REQUEST_FIELD_PARAMETER] = reqData
	req["method"] = ka.SvcList[svc].RestfulApiMethod
	req["url"] = ka.SvcList[svc].RestfulApiUrl
	ka.ManageMessage(req)
}

func (ka *KafkaAgentClient) ManageMessage(reqData map[string]interface{}) {
	// 2. 메시지를 파싱하여 얻은 정보를 토대로 restful server를 호출하는 단계
	restResp, err, handler := getRestServerResp(reqData, ka.AgentConfig, reqData[MESSAGE_TOPIC].(string))
	if err != nil {
		ka.RespToKafka(handler)
		return
	}

	// 3. restful server에서 온 응답을 처리하는 단계
	handler = handleRestServerResp(reqData, restResp, ka.AgentConfig)
	if handler != nil {
		ka.RespToKafka(handler)
	}
}

func (ka *KafkaAgentClient) ManageConfig(reqData map[string]interface{}) {
	param := reqData[REST_REQUEST_FIELD_PARAMETER].(map[string]interface{})
	if util.TypeToInt(reqData["svc"]) == 81 {
		configProcessor(ka.AgentConfig, param, "retry_config", "rtyCode", "rtyCycle", "rtyMCnt")
	} else if util.TypeToInt(reqData["svc"]) == 82 {
		configProcessor(ka.AgentConfig, param, "batch_svc_list", "btcCode", "btcCycle", "btcMCnt", "btcMethod", "btcUrl")
		if _, ok := ka.SvcList[param["btcCode"].(string)]; !ok {
			createTimer(ka, param["btcCode"].(string), param)
		}
	}
	ka.ManageMessage(reqData)
	MsgLogger.Println("message configuration is updated and sent kafka with", reqData)
}

// RespToKafka: KafkaAgentClient에서 설정된 producer를 통해 kafka에 메시지를 입력하는 함수
func (ka *KafkaAgentClient) RespToKafka(hr HandlerReturn) error {
	topic, value := hr()
	valueJson, err := json.Marshal(value)
	if err != nil {
		MsgLogger.Println("Failed to make json", err)
		return err
	}
	if topic == "Eventlistener Process" {
		MsgLogger.Println("Eventlistener processed - request:", value)
		return nil
	}

	// retry 인 것만 슬립 적용, 나머지 에러 또는 실패 케이스는 슬립없이 진행
	if infoVal, ok := value["retry-info"]; ok {
		if codeVal, ok := infoVal.(map[string]interface{})["resp_code"]; ok {
			resp_code := codeVal.(string)
			config := ka.AgentConfig["retry_config"].(map[string]interface{})
			if _, ok := config[resp_code]; ok {
				time.Sleep(time.Duration(config[resp_code].(map[string]interface{})["rtyCycle"].(float64)) * time.Second)
			}
		}
	}

	partition, offset, err := ka.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(valueJson),
		})
	if err != nil {
		MsgLogger.Println("Failed to put the message to kafka", err)
		return err
	}

	MsgLogger.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n\n", topic, partition, offset)

	return nil
}

/*
func (ka *KafkaAgentClient) Close() {
	if err := ka.Producer.Close(); err != nil {
		MsgLogger.Fatalln("Failed to close producer", err)
	}
	if err := ka.ConsumerGroup.Close(); err != nil {
		MsgLogger.Fatalln("Failed to close consumer group", err)
	}
}
*/
