package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/Shopify/sarama"

	"github.com/mycreditchain/kafka-agent/kafka-consumer-agent/client"
)

func main() {
	fmt.Println("num of CPU:", runtime.NumCPU())

	// CPU 개수를 구한 뒤 사용할 최대 CPU 개수 설정
	// go routine 생성시 이용
	runtime.GOMAXPROCS(runtime.NumCPU())

	// flag 설정. batch flag: -mode=batch/normal, -type=request/retry
	// mode가 normal인 경우 mode flag 생략 가능
	consumerMode := flag.String("mode", "batch", "mode can be normal or batch")
	consumerName := flag.String("name", "master", "name can be a part of log file name")
	flag.Parse()

	var err error
	client.MsgLogger, err = client.NewMsgLogger("./kafka-"+*consumerMode+"-agent-"+*consumerName+".log", "MESSAGE: ", false)
	if err != nil {
		log.Println("Failed to create message logger", err)
		return
	}

	// flag에서 설정한 mode에 따라
	// ~/workspace/go/src/github.com/mycreditchain/kafka-agent/kafka-consumer-agent/client/config.json
	// 에 설정된 config 값을 로드
	agentConfig := client.GetConfiguration("./client/config.json", *consumerMode)

	// kafka producer 설정
	pConfig := sarama.NewConfig()
	pConfig.Producer.RequiredAcks = sarama.WaitForLocal // Wait for all in-sync replicas to ack the message
	pConfig.Producer.Retry.Max = 10                     // Retry up to 10 times to produce the message
	pConfig.Producer.Return.Successes = true
	//pConfig.Producer.Partitioner = sarama.NewManualPartitioner
	/*
		if tlsConfig != nil {
			config.Net.TLS.Config = tlsConfig
			config.Net.TLS.Enable = true
		}
	*/

	producer, err := client.NewSyncProducer([]string{client.BROKER_URL_01, client.BROKER_URL_02, client.BROKER_URL_03}, pConfig)
	if err != nil {
		client.MsgLogger.Println("Failed to create new producer", err)
		return
	}

	// kafka consumer 설정
	cConfig := sarama.NewConfig()
	cConfig.ClientID = client.BATCH_CONSUMER_CLIENT_ID
	cConfig.Version = sarama.V0_10_2_0
	cConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	cConfig.Consumer.Return.Errors = true
	/*
		if tlsConfig != nil {
			config.Net.TLS.Config = tlsConfig
			config.Net.TLS.Enable = true
		}
	*/

	consumerGroup, err := client.NewConsumerGroup(client.BATCH_CONSUMER_GROUP_NAME, []string{client.BROKER_URL_01, client.BROKER_URL_02, client.BROKER_URL_03}, cConfig)
	if err != nil {
		client.MsgLogger.Println("Failed to create new consumer group", err)
		return
	}

	topics := []string{agentConfig[client.TOPIC_REQUEST].(string), agentConfig[client.TOPIC_RETRY].(string), "config", "sync"}
	kafkaRequestAgentClient := &client.KafkaAgentClient{
		AgentConfig: agentConfig,
		Producer:    producer,
		//ConsumerGroup: consumerGroup,
	}

	client.MsgLogger.Printf("consumer group: %+v\n", consumerGroup)
	client.ConsumeMessages(topics, consumerGroup, kafkaRequestAgentClient)
}
