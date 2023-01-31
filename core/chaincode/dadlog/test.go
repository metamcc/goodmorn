package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	cst "github.com/mycreditchain/chaincode/dadlog/constants"
	msg "github.com/mycreditchain/chaincode/dadlog/utils/messages"
	"github.com/mycreditchain/chaincode/dadlog/utils/tooltip"

	"strconv"
	"time"
	"strings"
)


// 테스트용
func (s *STRT_DAD_MAIN) makeDADCLogInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	refID := args[0] // refID 
	senderWalletAddr := args[1]	// 보낸이 지갑주소
	receiverWalletAddr := args[2]	// 받는이 지갑 주소
	paramSeedPenaltyStr := args[3]	// 씨앗양 - 페널티포함
	procDate := args[4] // 처리날짜
	timeZoneOffsetStr := args[5] // 지갑 timezone off set
	timeZoneOffsetInt, _ := strconv.Atoi(timeZoneOffsetStr)
	txId := stub.GetTxID()

	formatedDateLTZ := tooltip.FormatDateUTC2LTZ(procDate, "2006-01-02 15:04:05", "2006-01-02", timeZoneOffsetInt)
	keyArgs := []string{formatedDateLTZ[0:4], formatedDateLTZ[5:7], formatedDateLTZ[8:10], senderWalletAddr}
	keyArgs = append(keyArgs, receiverWalletAddr, refID, procDate, paramSeedPenaltyStr, txId)
	//keyArgs = append(keyArgs, []string{receiverWalletAddr, refID, procDate, paramSeedPenaltyStr, txId})
	res := s.makeCompositeKey(stub, cst.DADIndex, keyArgs)
	if res.Status != shim.OK {
		return shim.Error(msg.GetCCErrMsg("CC_COP_E01"))
	}

	
	coreChkInfo := RES_CC{formatedDateLTZ[0:4]+formatedDateLTZ[5:7]+formatedDateLTZ[8:10], senderWalletAddr, receiverWalletAddr, refID, procDate, paramSeedPenaltyStr, txId}
	coreChkInfoAsBytes, _ := json.Marshal(coreChkInfo)

	return shim.Success(coreChkInfoAsBytes)

}


func (s *STRT_DAD_MAIN) delDADCLogInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.DADIndex, args)
	
	if msgErr != nil {
		return shim.Error(msg.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
	}
	defer keyResultsIterator.Close()

	// Check the variable existed
	if !keyResultsIterator.HasNext() {
		return shim.Error(msg.GetCCErrMsg("CC_COR_E02", strings.Join(args[:],",")))
	}


	for keyResultsIterator.HasNext() {
		keyRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(msg.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		stub.DelState(keyRange.Key)
	}

	currentDate := time.Now().UTC() 	// 생성일
	currentDateStr := string(currentDate.Format("2006-01-02 15:04:05"))
	var resMsg = RES_NOR{"Ok", currentDateStr, stub.GetTxID()}

	resMsgAsBytes, _ := json.Marshal(resMsg)

	return shim.Success(resMsgAsBytes)


}
