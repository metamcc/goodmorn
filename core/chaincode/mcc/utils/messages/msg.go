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


func GetCCErrMsg(msgParam ...string) string {
	return getMsg(msgParam)
}

var msgDataMap map[string]string

var varDataMap map[string]string

var fncDataMap map[string]string

func SetMgmtData(data map[string]string, dataKey string) {

	if dataKey == "msg" {
		msgDataMap = data
	} else if dataKey == "var" {
		varDataMap =  data
	} else if dataKey == "fnc" {
		fncDataMap =  data
	}
}

func GetVarData(valuekey string) string {

	return varDataMap[valuekey]
}

func GetMsgData(valuekey string) string {

	return msgDataMap[valuekey]
}


func GetFncData(valuekey string) (string, bool) {

	 value, ok :=  fncDataMap[valuekey]

	 return value, ok
}

func getMsg(msgParam []string) string {

	msgMap := msgDataMap

	var msg string
	msgLen := len(msgParam)

	if msgLen > 0 {

		msgBody := msgMap[msgParam[0]]

		//fmt.Println(msgBody)
		if msgBody == "" {
			return msgMap["RS_MSG_E01"]
		}

		cnt := strings.Count(msgBody, "%")
		if msgLen == 1 {
			if cnt > 0 {
				return msgMap["RS_MSG_E02"]
			}
			msg = msgBody
		} else if msgLen > 1 {

			if cnt != msgLen-1 {
				return msgMap["RS_MSG_E02"]
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
				return msgMap["RS_MSG_E02"]
			}

		}

		msg = "[" + msgParam[0] + "] " + msg
	}


	return msg
}
