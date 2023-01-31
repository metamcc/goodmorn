package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	cst "github.com/mycreditchain/chaincode/dadlog/constants"
	mgmt "github.com/mycreditchain/chaincode/dadlog/utils/messages"
	"strconv"
)

/*
관리 데이터 초기화 함수
*/
func (s *STRT_DAD_MAIN) initAdmMgmtData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var msgDataMap map[string]string
	var varDataMap map[string]string
	var fncDataMap map[string]string

	var errStr string
	var msgCnt, varCnt, fncCnt int 

	typ,_ := strconv.Atoi(args[0])
	
	if typ == 1 || typ == 0 {
		// 메세지 데이터 세팅
		msgDataMap, errStr = s.callAdmMgmtData(stub, []string{"getCoreMsg"})
		if errStr != "" {
			return shim.Error(errStr)
		}

		msgCnt = len(msgDataMap)
		mgmt.SetMgmtData(msgDataMap, "msg")

	}

	if typ == 2 || typ == 0 { 
		// 관리 변수 세팅
		varDataMap, errStr = s.callAdmMgmtData(stub, []string{"getCoreVar"})
		if errStr != "" {
			return shim.Error(errStr)
		}

		varCnt = len(varDataMap)
		mgmt.SetMgmtData(varDataMap, "var")
	}


	if typ == 3 || typ == 0 { 
		// 함수 사용 여부 세팅
		fncDataMap, errStr = s.callAdmMgmtData(stub, []string{"getFncEnable"})
		if errStr != "" {
			fncCnt = 0
		} else {
			fncCnt = len(fncDataMap)
			mgmt.SetMgmtData(fncDataMap, "fnc")
		}

	}
	
	mgmtDataInfo := STRT_ADM2{msgDataMap, varDataMap, fncDataMap, msgCnt, varCnt, fncCnt}
	mgmtDataJSON, _ := json.Marshal(mgmtDataInfo)

	return shim.Success(mgmtDataJSON)
}

/*
관리 데이터 취득용 공용 함수
*/
func (s *STRT_DAD_MAIN) callAdmMgmtData(stub shim.ChaincodeStubInterface, args []string) (map[string]string, string) {

	if len(args) != 1 {
		return nil, mgmt.GetCCErrMsg("CC_COM_E01", "1")
	}

	funcNm := args[0] // 함수명
	chaincodeName := cst.Config_CHAINCODE

	ccArgs := []string{funcNm}
	invokeArgs := util.ArrayToChaincodeArgs(ccArgs)

	ccRes := stub.InvokeChaincode(chaincodeName, invokeArgs, stub.GetChannelID())
	if ccRes.Status != shim.OK {
		return nil, ccRes.Message
	}

	ccPayload := STRT_ADM{}
	json.Unmarshal(ccRes.Payload, &ccPayload)

	return ccPayload.MgmtDataMap, ""
}

/*
 컴포지트 키 생성
*/
func (s *STRT_DAD_MAIN) makeCompositeKey(stub shim.ChaincodeStubInterface, keyName string, args []string) pb.Response {

	compositeKey, compositeErr := stub.CreateCompositeKey(keyName, args)
	if compositeErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositeErr.Error()))
	}

	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositePutErr.Error()))
	}

	return shim.Success(nil)
}