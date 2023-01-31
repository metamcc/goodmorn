package msg

import (
	"fmt"
	"strings"
)

type CommonMsg struct {
	Code    string
	Desc    string
	Message string
}

type ResponseData struct {
	CommonMsg
	Data interface{}
}

var sysCodeMap = map[string]string{

	"RS_RES_N01": "ok",
}

var evtDataMap = map[string]string{

	"evtCd": "ev0001",
	"evtDt": "2019-01-07 00:00:00",
	"evtAmnt": "20000000000000000000",
	"evtCnt": "1",
	"etc": "마케팅열매선물이벤트",
}

var varDataMap = map[string]string{

	"CORE_SEED_PENALTY_YN": "false",
	"DAD_LIMIT_AMNT": "10000000000000000000000",
	"SEED_LIMIT_CNT": "5",
	"SND_DATA_SYNC_CNT": "9",
	"RCV_DATA_SYNC_CNT": "99",
}

var fncDataMap = map[string]string{

	"getWalletInfo": "true",
}


var retryCodeMap = map[string]string{

	"CC_SSG_E05": "3",                              // 당일 씨앗 선물 횟수 충족 체크
	"MVCC_READ_CONFLICT": "10",
	"PHANTOM_READ_CONFLICT": "10",
	"CONNECTION_FAILED": "10",
}


var msgDataMap = map[string]string{



	// 공통 메세지
	"CC_COM_E00": "Wrong Admin Password",
	"CC_COM_E01": "Incorrect number of arguments. Expecting %s",
	"CC_COM_E02": "Invalid function name",
	"CC_COM_E03": "Failed to execute PutState, %s",
	"CC_COM_E04": "No data Found - GetState",

	"CC_COM_E08": "Sub Calk Error , (%s) must be bigger than (%s)",
	"CC_COM_E09": "Error calling other chaincode",
	"CC_COM_E10": "Failed to convert value, %s : %s",
	"CC_COM_E11": "Execution of function %s is temporary stopped",
	"CC_COM_E12": "Failed to start chaincode, %s",
	"CC_COM_E13": "Failed to execute DelState, %s",

	"CC_COM_E99": "Unkonwn Error, %s",
	
	// signing 관련 메세지
	"CC_SGN_E01": "R \"%s\" S \"%s\" verification parameters not numeric",
	"CC_SGN_E02": "Wallet address %s and key %s from X %s Y %s does not match",
	"CC_SGN_E03": "Wallet verification failed for parameters: %s ",


	// 컴포짓 키 관련 메세지
	"CC_COP_E01": "Failed to make composite key, %s" ,
	"CC_COP_E02": "Failed to get composite key, %s",


	// 지갑 관련 메세지
	"CC_WLT_E01": "Wallet does not exist, %s",
	"CC_WLT_E02": "Same Wallet Error, %s : %s",
	"CC_WLT_E03": "Wrong wallet type, %s - %s : %s",
	"CC_WLT_E04": "eco unique wallet already exists",
	"CC_WLT_E05": "User identity is duplicated, %s",


	// 지갑 동기화 관련 메세지
	"CC_SNC_E01": "PaymentData does not exist, %s",
	"CC_SNC_E02": "SyncRefID does not match, wallet : %s, refID : %s - %s, syncCnt : %s, keyCheckVal : %s",
	

	// 열매 관련 메세지


	// 씨앗 관련 메세지
	"CC_SSG_E01": "You have already sent seeds to that wallet, %s",
	"CC_SSG_E02": "All seeds per day were sent, %s",
	"CC_SSG_E03": "Penalty seed calculation is defferent - %s : %s",
	"CC_SSG_E04": "Max seeding count Error, %s : %s",
	"CC_SSG_E05": "User %s seeding count is not match for MAX - %s : %s",

	// 코어 관련 메세지
	"CC_COR_E01": "CORE Reference ID already existed, %s",
	"CC_COR_E02": "CORE Reference ID does not exist, %s",

	// 이벤트 관련 메세지
	"CC_EVT_E01": "SetEvent Error, %s : %s",
	"CC_EVT_E02": "%s Event amount is defferent, %s : %s",
	
	// 메세지 에러 관련
	"CC_MSG_E01": "incorrect message code",
	"CC_MSG_E02": "incorrect number of message arguments",


}

func GetDefaultMsgMap() map[string]string {
	return msgDataMap
}

func GetDefaultVarMap() map[string]string {
	return varDataMap
}

func GetDefaultFncMap() map[string]string {
	return fncDataMap
}

func GetDefaultRetryMap() map[string]string {
	return retryCodeMap
}

func GetDefaultEvtMap() map[string]string {
	return evtDataMap
}

func GetMsg(msgParam ...string) string {
	return getMsg(1, msgParam)
}

func GetErrMsg(msgParam ...string) string {
	return getMsg(2, msgParam)
}

func GetCCErrMsg(msgParam ...string) string {
	return getMsg(3, msgParam)
}

func getMsg(msgType int, msgParam []string) string {

	msgMap := msgDataMap
	if msgType == 1 {
		msgMap = sysCodeMap
	}

	var msg string
	msgLen := len(msgParam)

	if msgLen > 0 {

		msgBody := msgMap[msgParam[0]]

		//fmt.Println(msgBody)
		if msgBody == "" {
			return msgDataMap["CC_MSG_E01"]
		}

		cnt := strings.Count(msgBody, "%")
		if msgLen == 1 {
			if cnt > 0 {
				return msgDataMap["CC_MSG_E02"]
			}
			msg = msgBody
		} else if msgLen > 1 {

			if cnt != msgLen-1 {
				return msgDataMap["CC_MSG_E02"]
			}

			if msgLen == 2 {
				msg = fmt.Sprintf(msgBody, msgParam[1])
			} else if msgLen == 3 {
				msg = fmt.Sprintf(msgBody, msgParam[1], msgParam[2])
			} else if msgLen == 4 {
				msg = fmt.Sprintf(msgBody, msgParam[1], msgParam[2], msgParam[3])
			} else if msgLen == 5 {
				msg = fmt.Sprintf(msgBody, msgParam[1], msgParam[2], msgParam[3], msgParam[4])
			} else {
				return msgDataMap["CC_MSG_E02"]
			}

		}

	}

	if msgType == 3 {
		msg = "[" + msgParam[0] + "] " + msg
	}
	return msg
}
