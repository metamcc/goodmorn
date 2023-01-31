package main

import (
	"encoding/json"
//	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	cst "github.com/mycreditchain/chaincode/confcc/constants"
	msg "github.com/mycreditchain/chaincode/confcc/utils/messages"
)

var logger = shim.NewLogger("confcc")

// Init - chaincode instantioation
func (s *MCCConfigChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//return shim.Success(nil)
	return s.initDefaultCoreData(stub)
}


// Invoke method
func (s *MCCConfigChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "initDefaultCoreData" { 
		return s.initDefaultCoreData(stub)
	} else if function == "setCoreMsg" { 
		return s.setCoreMsg(stub, args)
	} else if function == "getCoreMsg" { 
		return s.getCoreMsg(stub)
	} else if function == "setCoreVar" { 
		return s.setCoreVar(stub, args)
	} else if function == "getCoreVar" {
		return s.getCoreVar(stub)
	} else if function == "setFncEnable" { 
		return s.setFncEnable(stub, args)
	} else if function == "getFncEnable" {
		return s.getFncEnable(stub)
	} else if function == "setCoreEvt" { 
		return s.setCoreEvt(stub, args)
	} else if function == "getCoreEvt" {
		return s.getCoreEvt(stub, args)
	} else if function == "setAdmPwd" { 
		return s.setAdmPwd(stub, args)
	} else if function == "getAdmPwd" {
		return s.getAdmPwd(stub)
	}  

	return shim.Error(msg.GetCCErrMsg("CC_COM_E02"))
}


/*
초기 메세지 데이터 생성
*/
func (s *MCCConfigChaincode) initDefaultCoreData(stub shim.ChaincodeStubInterface) pb.Response {
	
	// 1. 관리자비번 
	admStrt := STRT_PWD{cst.AdmPwd}
	admStrtAsBytes, _ := json.Marshal(admStrt)

	err := stub.PutState(cst.AdmKey, admStrtAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.AdmKey))
	}

	// 1. 메세지 
	msgDataMap := msg.GetDefaultMsgMap()
	admMgmtConfigStrt := STRT_ADM{msgDataMap, len(msgDataMap)}
	admMgmtConfigAsBytes, _ := json.Marshal(admMgmtConfigStrt)

	err = stub.PutState(cst.MsgKey, admMgmtConfigAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.MsgKey))
	}

	// 2. 변수 
	varDataMap := msg.GetDefaultVarMap()
	admMgmtConfigStrt = STRT_ADM{varDataMap, len(varDataMap)}
	admMgmtConfigAsBytes, _ = json.Marshal(admMgmtConfigStrt)

	err = stub.PutState(cst.VarKey, admMgmtConfigAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.VarKey))
	}


	// 3. 함수 
	fncDataMap := msg.GetDefaultFncMap()
	admMgmtConfigStrt = STRT_ADM{fncDataMap, len(fncDataMap)}
	admMgmtConfigAsBytes, _ = json.Marshal(admMgmtConfigStrt)

	err = stub.PutState(cst.FncKey, admMgmtConfigAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.FncKey))
	}


	// 4. 재처리대상코드 
	retryCodeMap := msg.GetDefaultRetryMap()
	admMgmtConfigStrt = STRT_ADM{retryCodeMap, len(retryCodeMap)}
	admMgmtConfigAsBytes, _ = json.Marshal(admMgmtConfigStrt)

	err = stub.PutState(cst.RetryKey, admMgmtConfigAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.RetryKey))
	}


	// 4. 이벤트 
	evtDataMap := msg.GetDefaultEvtMap()
	evnCnt, _ := strconv.Atoi(evtDataMap["evtCnt"])
	eventStrt := STRT_EVT{evtDataMap["evtCd"], evtDataMap["evtDt"], evtDataMap["evtAmnt"], evnCnt, evtDataMap["etc"]}
	eventStrtAsBytes, _ := json.Marshal(eventStrt)

	err = stub.PutState(cst.EvtKey+"_"+evtDataMap["evtCd"], eventStrtAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.EvtKey+"_"+evtDataMap["evtCd"]))
	}

	return shim.Success(eventStrtAsBytes)
}


/*
관리자 비번 설정
agrg[0] : 이전 비번
agrg[1] : 새로운 비번
*/
func (s *MCCConfigChaincode) setAdmPwd(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
    if len(args) != 2 {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E01", "2"))
	}

	// 파라메터 
	admPwdBefroe := args[0]
	admPwdBeAfter := args[1]

	// 기존 내용 취득
	admMgmtConfigAsBytes, _ := stub.GetState(cst.AdmKey)
	if admMgmtConfigAsBytes == nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E04"))
	}
	admMgmtConfigStrt := STRT_PWD{}
	json.Unmarshal(admMgmtConfigAsBytes, &admMgmtConfigStrt)

	// 기존 비번 비교
	if admPwdBefroe != admMgmtConfigStrt.Pwd {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E00"))
	}

	// 신규 비번 설정
	admMgmtConfigStrt.Pwd = admPwdBeAfter
	admMgmtConfigAsBytes, _ = json.Marshal(admMgmtConfigStrt)
	err := stub.PutState(cst.AdmKey, admMgmtConfigAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.AdmKey))
	}

	return shim.Success(admMgmtConfigAsBytes)
}

/*
코어 관리 데이터 설정
agrg[0] : 처리구분 -  1: 설정(defalut), 0- 삭제
agrg[1] : 설정키
agrg[2] : 설정값
agrg[3] : 데이터 구분키
*/
func (s *MCCConfigChaincode) setCoreAdmData(stub shim.ChaincodeStubInterface, args []string, dataKey string) pb.Response {
	
	if len(args) != 3 {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E01", "3"))
	}

	actType := args[0] // 처리구분 -  1: 설정(defalut), 0- 삭제
	codeName := args[1]   // 설정키
	codevalue := args[2] // 설정값
	//dataKey := args[3] // 데이터구분키
	
	var mgmtDataMap map[string]string
	admMgmtConfigStrt := STRT_ADM{}


	// 설정 데이터 맵 취득
	admMgmtConfigAsBytes, _ := stub.GetState(dataKey)
	if admMgmtConfigAsBytes == nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E04"))
	} else {

		json.Unmarshal(admMgmtConfigAsBytes, &admMgmtConfigStrt)
		mgmtDataMap = admMgmtConfigStrt.MgmtDataMap
	}

	if actType == "1" { // 설정
		mgmtDataMap[codeName] = codevalue
	} else { // 삭제
		delete(mgmtDataMap, codeName)
	}
	
	// 기록
	admMgmtConfigStrt.MgmtDataMap = mgmtDataMap
	admMgmtConfigStrt.Cnt = len(mgmtDataMap)
	admMgmtConfigAsBytes, _ = json.Marshal(admMgmtConfigStrt)

	err := stub.PutState(dataKey, admMgmtConfigAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", dataKey))
	}

	return shim.Success(admMgmtConfigAsBytes)
}


/*
관리자 비밀번호 취득
*/
func (s *MCCConfigChaincode) getAdmPwd(stub shim.ChaincodeStubInterface) pb.Response {
	return s.getCoreAdmData(stub, []string{cst.AdmKey})
}

/*
코어 관리 데이터 취득
agrg[0] : 데이터 구분키
*/
func (s *MCCConfigChaincode) getCoreAdmData(stub shim.ChaincodeStubInterface, args []string ) pb.Response {

	dataKey := args[0] // 데이터구분키

	if len(args) != 1 {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E01", "1"))
	}

	admMgmtConfigAsBytes, _ := stub.GetState(dataKey)
	if admMgmtConfigAsBytes == nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E04"))
	}

	return shim.Success(admMgmtConfigAsBytes)
}

/*
코어 메세지 설정
agrg[0] : 처리구분 -  1: 설정(defalut), 0- 삭제
agrg[1] : 설정키
agrg[2] : 설정값
*/
func (s *MCCConfigChaincode) setCoreMsg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return s.setCoreAdmData(stub, args, cst.MsgKey)
}

/*
전체 메세지 취득
*/
func (s *MCCConfigChaincode) getCoreMsg(stub shim.ChaincodeStubInterface) pb.Response {
	return s.getCoreAdmData(stub, []string{cst.MsgKey})
}

/*
코어 관리 변수 설정
agrg[0] : 처리구분 -  1: 설정(defalut), 0- 삭제
agrg[1] : 설정키
agrg[2] : 설정값
*/
func (s *MCCConfigChaincode) setCoreVar(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return s.setCoreAdmData(stub, args, cst.VarKey)
}


/*
코어 관리 변수 취득
agrg[0] : 설정 키
*/
func (s *MCCConfigChaincode) getCoreVar(stub shim.ChaincodeStubInterface) pb.Response {
	return s.getCoreAdmData(stub, []string{cst.VarKey})
}


/*
코어 관리 변수 취득 - 맵타입
agrg[0] :  설정 키
*/
func (s *MCCConfigChaincode) getCoreVarMap(stub shim.ChaincodeStubInterface) map[string]string {

	var mgmtDataMap map[string]string

	admMgmtConfigAsBytes, _ := stub.GetState(cst.VarKey)
	if admMgmtConfigAsBytes == nil {
		mgmtDataMap = make(map[string]string)
	}

	admMgmtConfigStrt := STRT_ADM{}
	json.Unmarshal(admMgmtConfigAsBytes, &admMgmtConfigStrt)

	mgmtDataMap = admMgmtConfigStrt.MgmtDataMap

	return mgmtDataMap
}


/*
코어 함수 사용 유무 설정
agrg[0] : 처리구분 -  1: 설정(defalut), 0- 삭제
agrg[1] : 설정키
agrg[2] : 설정값
*/
func (s *MCCConfigChaincode) setFncEnable(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return s.setCoreAdmData(stub, args, cst.FncKey)
}



/*
코어 함수 사용 유무 설정
*/
func (s *MCCConfigChaincode) getFncEnable(stub shim.ChaincodeStubInterface) pb.Response {
	return s.getCoreAdmData(stub, []string{cst.FncKey})
}

/*
이벤트 정보 설정
agrg[0] : 이벤트코드
agrg[1] : 이벤트기준일자
agrg[2] : 이벤트열매량
agrg[3] : 이벤트참가횟수 - 0 - 매번 허용 , 1 : 유일 1회 허용,  2 - 2회 허용....
agrg[4] : 기타
*/
func (s *MCCConfigChaincode) setCoreEvt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	

    if len(args) != 5 {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E01", "5"))
	}

	evtCd := args[0]
	evtDt := args[1]
	evtAmnt := args[2]
	evtCnt,_ := strconv.Atoi(args[3])
	etc := args[4]

	eventConfig := STRT_EVT{evtCd, evtDt, evtAmnt, evtCnt, etc}
	eventConfigAsBytes, _ := json.Marshal(eventConfig)

	
	err := stub.PutState(cst.EvtKey+"_"+evtCd, eventConfigAsBytes)
	if err != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E03", cst.EvtKey+"_"+evtCd))
	}

	return shim.Success(eventConfigAsBytes)
}

/*
이벤트 정보 취득
agrg[0] : 이벤트코드
*/
func (s *MCCConfigChaincode) getCoreEvt(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error(msg.GetCCErrMsg("CC_COM_E01", "1"))
	}

	eventConfigAsBytes, _ := stub.GetState(cst.EvtKey+"_"+args[0])

	if eventConfigAsBytes == nil {
		
		return shim.Error(msg.GetCCErrMsg("CC_COM_E04"))

	}

	return shim.Success(eventConfigAsBytes)
}

//
func main() {
	err := shim.Start(new(MCCConfigChaincode))
	if err != nil {
		logger.Errorf(msg.GetCCErrMsg("CC_COM_E12"), err)
	}
}