package main


import (
	"fmt"
	"encoding/json"
	"strings"

	cst "github.com/mycreditchain/chaincode/mcc/constants"
	mgmt "github.com/mycreditchain/chaincode/mcc/utils/messages"
	mw "github.com/mycreditchain/chaincode/mcc/utils/wallet"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)


/*
관리 데이터 초기화 함수
*/
func (s *MCCWalletChaincode) initAdmMgmtData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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
func (s *MCCWalletChaincode) callAdmMgmtData(stub shim.ChaincodeStubInterface, args []string) (map[string]string, string) {

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
func (s *MCCWalletChaincode) makeCompositeKey(stub shim.ChaincodeStubInterface, keyName string, args []string) pb.Response {

	compositeKey, compositeErr := stub.CreateCompositeKey(keyName, args)
	if compositeErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositeErr.Error()))
	}

	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", compositeKey))
	}

	return shim.Success(nil)
}



// signing 검증
// args[0] : public Key
// args[1] : r
// args[2] : s
// args[3] : plainTxt for hash
func (s *MCCWalletChaincode) verifyArguments(stub shim.ChaincodeStubInterface, args []string) (bool, string) {
	if len(args) != 4 {
		return false, mgmt.GetCCErrMsg("CC_COM_E01", "4")
	}

	hash := mw.GetHash(args[3])

	// Args: pubKey, hash, r, s
	isVerified := mw.Verify(args[0], hash, args[1], args[2])

	return isVerified, hash
}


// 컴포짓 키 생성  
func (s *MCCWalletChaincode) createCompKey(stub shim.ChaincodeStubInterface, indexKey string, keyArgs []string) error {

	// ---------------코어 검증용 컴포짓 키 구성----------------------
//	keyArgs := []string{svcPrefix, refID, txId}
fmt.Println("indexKey ======="+indexKey)
fmt.Println(keyArgs)
	compositeKey, compositeErr := stub.CreateCompositeKey(indexKey, keyArgs)
	if compositeErr != nil {
		return compositeErr
	}

	fmt.Println("compositeKey ======="+compositeKey)

	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return compositePutErr
	}

	return nil

}


// 컴포짓 키 체크용 -  
func (s *MCCWalletChaincode) createCoreCompKey(stub shim.ChaincodeStubInterface, keyArgs []string) error {
	fmt.Println("cst.IndexCoreCheck ======="+cst.IndexCoreCheck)
	return s.createCompKey(stub, cst.IndexCoreCheck, keyArgs)
	/*
		// ---------------코어 검증용 컴포짓 키 구성----------------------
	//	keyArgs := []string{svcPrefix, refID, txId}

		compositeKey, compositeErr := stub.CreateCompositeKey(cst.IndexCoreCheck, keyArgs)
		if compositeErr != nil {
			return compositeErr
		}
		compositePutErr := stub.PutState(compositeKey, []byte{0x00})
		if compositePutErr != nil {
			return compositePutErr
		}

		return nil
*/
}

/*
	코어 refID 검색
*/
func (s *MCCWalletChaincode) getCoreCompKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var resMsg []RES_CC 

	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexCoreCheck, args)
	
	if msgErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
	}
	defer keyResultsIterator.Close()

	// Check the variable existed
	if !keyResultsIterator.HasNext() {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E02", strings.Join(args[:],",")))
	}

	// Iterate through result set and compute final value
	var i int
	for i = 0; keyResultsIterator.HasNext(); i++ {
		responseRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", splitKeyErr.Error()))
		}

		resMsg = append(resMsg, RES_CC{keyParts[0]+keyParts[1]+keyParts[2], keyParts[3], keyParts[4], keyParts[5]})
	}

	coreChkInfo := RES_CC_LST{CoreChkInfo: resMsg}
	coreChkInfoAsBytes, _ := json.Marshal(coreChkInfo)

	return shim.Success(coreChkInfoAsBytes)

}	

/*
	코어 refID 검색 - 갯수만 반환
*/
func (s *MCCWalletChaincode) getCoreCompData(stub shim.ChaincodeStubInterface, args []string) (cnt int, err error) {

	resultCnt := 0
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexCoreCheck, args)

	if msgErr != nil {
		return resultCnt, err
	}
	defer keyResultsIterator.Close()

	// Check the variable existed
	if !keyResultsIterator.HasNext() {
		return resultCnt, nil
	}

	// Iterate through result set and compute final value
	var i int
	for i = resultCnt; keyResultsIterator.HasNext(); i++ {
		responseRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			return i, err
		}
	
		_, _, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return i, err
		}

	}



	return i, err

}




// args[0] : function name
func (s *MCCWalletChaincode) callConfigurationCC(stub shim.ChaincodeStubInterface, args []string) string {
	if len(args) != 1 {
		return mgmt.GetCCErrMsg("CC_COM_E01", "1")
	}

	chaincodeName := cst.Config_CHAINCODE
	channelID := cst.MCCChannelID

	ccArgs := []string{args[0]}
	invokeArgs := util.ArrayToChaincodeArgs(ccArgs)

	ccRes := stub.InvokeChaincode(chaincodeName, invokeArgs, channelID)
	if ccRes.Status != shim.OK {
		if ccRes.Message == "" {
			ccRes.Message = fmt.Sprintf("Chaincode %s is not found or returns no value", chaincodeName)
		}
		return fmt.Sprintf(mgmt.GetCCErrMsg("CC_COM_E09") + " : " + ccRes.Message)
	}

	ccPayload := MCCConfig{}
	json.Unmarshal(ccRes.Payload, &ccPayload)

	return ccPayload.Configuration
}



func (s *MCCWalletChaincode) callConfigurationData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	chaincodeName := cst.Config_CHAINCODE
	channelID := cst.MCCChannelID

	invokeArgs := util.ArrayToChaincodeArgs(args)
	ccRes := stub.InvokeChaincode(chaincodeName, invokeArgs, channelID)


	return ccRes
}


func (s *MCCWalletChaincode) getMgmtConfData(stub shim.ChaincodeStubInterface, args []string) (MCCConfig, string) {

	ccRes := s.callConfigurationData(stub, args)

	errMsg := ""

	confStrt := MCCConfig{}

	if ccRes.Status == shim.OK {

		json.Unmarshal(ccRes.Payload, &confStrt)
	} else {
		errMsg = ccRes.Message
	}

	return confStrt, errMsg
}


func (s *MCCWalletChaincode) getEventConfData(stub shim.ChaincodeStubInterface, args []string) (STRT_EVT, string) {

	ccRes := s.callConfigurationData(stub, args)

	fmt.Println(ccRes)

	errMsg := ""
	eventStrt := STRT_EVT{}

	if ccRes.Status == shim.OK {

		json.Unmarshal(ccRes.Payload, &eventStrt)
	} else {

		errMsg = ccRes.Message
	}

	return eventStrt, errMsg
}


func (s *MCCWalletChaincode) callAdminMgmtData(stub shim.ChaincodeStubInterface, args []string)  (map[string]string, string) {

	ccRes := s.callConfigurationData(stub, args)

	errMsg := ""

	confStrt := STRT_ADM{}

	if ccRes.Status == shim.OK {

		json.Unmarshal(ccRes.Payload, &confStrt)
		
	} else {
		errMsg = ccRes.Message
	}

	return confStrt.MgmtDataMap, errMsg

}


func (s *MCCWalletChaincode) setErrStructDataWidthSSG(seedingSuccessMsg []RES_SSG, senderWalletAddr string,receiverWalletAddr string,paramSeedBaseStr string,paramSeedPenaltyStr string , procDate string, refID string, msg string)  []RES_SSG {

	//seedStrt := RES_SSG{senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,"",procDate,"", refID, msg}
	fmt.Println("--------------------------")
	fmt.Println(msg)
	fmt.Println("--------------------------")
	seedingSuccessMsg = append(seedingSuccessMsg, RES_SSG{senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,"",procDate,"", refID, msg})

	return seedingSuccessMsg

} 