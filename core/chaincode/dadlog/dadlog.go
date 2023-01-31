package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	cst "github.com/mycreditchain/chaincode/dadlog/constants"
	mgmt "github.com/mycreditchain/chaincode/dadlog/utils/messages"
	"github.com/mycreditchain/chaincode/dadlog/utils/tooltip"

	"strconv"
	"strings"
	"time"
)

var logger = shim.NewLogger("dadlog")

func main() {
	err := shim.Start(new(STRT_DAD_MAIN))
	if err != nil {
		logger.Errorf(mgmt.GetCCErrMsg("CC_COM_E12"), err)
	}
}

// Init - chaincode instantioation
func (s *STRT_DAD_MAIN) Init(stub shim.ChaincodeStubInterface) pb.Response {
	
	// 관리 데이터 최신으로 반영
	return s.initAdmMgmtData(stub, []string{"0"})
	//return shim.Success(nil)
}


// Invoke - choose function to interact with ledger
// args[0] : function_name
func (s *STRT_DAD_MAIN) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	useYn, errStr := s.checkFncUsing(stub, []string{function})
	if errStr != "" {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E99", errStr))
	} else {
		if  useYn == false {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E11", function))
		}
	}

	if function == "put" {
		return s.putDadLog(stub, args)
	} else if function == "makeDADCLogInfo" {
		return s.makeDADCLogInfo(stub, args)
	} else if function == "getDADCLogInfo" {
		return s.getDADCLogInfo(stub, args)
	} else if function == "delDADCLogInfo" {
		return s.delDADCLogInfo(stub, args)
	} else if function == "initAdmMgmtData" {
		return s.initAdmMgmtData(stub, args)
	} 

	return shim.Error(mgmt.GetCCErrMsg("CC_COM_E02"))
}


/*
함수 사용 유무 처리
args[0] : 함수이름
*/
func (s *STRT_DAD_MAIN) checkFncUsing(stub shim.ChaincodeStubInterface, args []string) (bool, string) {

	useYn := true
	var err error

	// 1. 파라메터 갯수 체크
	if len(args) != 1 {
		return useYn, mgmt.GetCCErrMsg("CC_COM_E01", "1")
	}

	fncNm := args[0]
	if boolStr, exsit := mgmt.GetFncData(fncNm); exsit {

		useYn, err = strconv.ParseBool(boolStr)
		if err != nil {
			return useYn, mgmt.GetCCErrMsg("CC_COM_E10", "fncNm bool", boolStr)
		}
	}  
	 
	return useYn, ""
}

/*
씨앗 선물 데이터 취득 YY~MM~DD~sender~receiver~refID~utcTime~seed~txID"
args[0] : 년 - 로컬
args[1] : 월 - 로컬
args[2] : 일 - 로컬
args[3] : 보낸지갑
args[4] : 받은지갑
args[5] : refID
args[6] : 처리일자 - UTC
*/
func (s *STRT_DAD_MAIN) getDADCLogInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var resMsg []RES_CC 
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.DADIndex, args)
	
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

		resMsg = append(resMsg, RES_CC{keyParts[0]+keyParts[1]+keyParts[2], keyParts[3], keyParts[4], keyParts[5], keyParts[6], keyParts[7], keyParts[8]})
	}

	coreChkInfo := RES_CC_LST{CoreChkInfo: resMsg}
	coreChkInfoAsBytes, _ := json.Marshal(coreChkInfo)

	return shim.Success(coreChkInfoAsBytes)

}	


/*
씨앗 선물 
args[0] : refID
args[1] : senderWalletAddr
args[2] : receiverWalletAddr
args[3] : paramSeedBase
args[4] : paramSeedPenalty
args[5] : procDate
args[6] : user timezone offset seconds
*/
func (s *STRT_DAD_MAIN) putDadLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수
	if len(args) != 7 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "7"))
	}

	// 2. 변수 선언
	refID := args[0] // refID 
	senderWalletAddr := args[1]	// 보낸이 지갑주소
	receiverWalletAddr := args[2]	// 받는이 지갑 주소
	paramSeedBaseStr := args[3] // 개인별 씨앗 선물 기준 갯수
	paramSeedPenaltyStr := args[4]	// 씨앗양 - 페널티포함
	procDate := args[5] // 처리날짜
	timeZoneOffsetStr := args[6] // 지갑 timezone off set

	
	timeZoneOffsetInt, _ := strconv.Atoi(timeZoneOffsetStr)
	paramSeedPenaltyFloat, err := strconv.ParseFloat(paramSeedPenaltyStr, 64)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", paramSeedPenaltyStr, err.Error()))

	}
	coreSeedPenaltyFloat := paramSeedPenaltyFloat

	formatedDateLTZ := tooltip.FormatDateUTC2LTZ(procDate, "2006-01-02 15:04:05", "2006-01-02", timeZoneOffsetInt)
	keyArgs := []string{formatedDateLTZ[0:4], formatedDateLTZ[5:7], formatedDateLTZ[8:10], senderWalletAddr}


	// 코어 체크 여부
	coreSeedPenaltyYnBool, err := strconv.ParseBool(mgmt.GetVarData(cst.CORE_SEED_PENALTY_YN))
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", cst.CORE_SEED_PENALTY_YN, err.Error()))
	}


	logger.Info("coreSeedPenaltyYnBool = "+fmt.Sprint(coreSeedPenaltyYnBool))
	logger.Info("refID = "+refID)
	logger.Info("resenderWalletAddrfID = "+senderWalletAddr)
	logger.Info("receiverWalletAddr = "+receiverWalletAddr)
	logger.Info("paramSeedBaseStr = "+paramSeedBaseStr)
	logger.Info("paramSeedPenaltyStr = "+paramSeedPenaltyStr)
	logger.Info("procDate = "+procDate)
	logger.Info("timeZoneOffsetStr = "+timeZoneOffsetStr)

	if coreSeedPenaltyYnBool == true {

		// 3. 유효성 체크
		// ---------------일 씨앗 선물 제한 체크-------------------------------
		keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.DADIndex, keyArgs)
		if msgErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
		}
		defer keyResultsIterator.Close()

		var sendingCnt int
		for sendingCnt = 0; keyResultsIterator.HasNext(); sendingCnt++ {
			responseRange, nextErr := keyResultsIterator.Next()
			if nextErr != nil {
				return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
			}

			_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
			if splitKeyErr != nil {
				return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", splitKeyErr.Error()))
			}

			if receiverWalletAddr == keyParts[4] { // 이미 받은 유저
				logger.Info(receiverWalletAddr)
				logger.Info(mgmt.GetCCErrMsg("CC_SSG_E01", receiverWalletAddr))

				return shim.Error(mgmt.GetCCErrMsg("CC_SSG_E01", receiverWalletAddr)) 
			}
		}


		// 일일 씨앗 선물 최대 갯수
		seedLimitCntInt, err := strconv.Atoi(mgmt.GetVarData(cst.SEED_LIMIT_CNT)) 
		if err != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", cst.SEED_LIMIT_CNT, err.Error()))
		}

		if sendingCnt+1 > seedLimitCntInt { // 기존갯수+현재요청건수 > 최대건수 : 이미 그날의 모든 씨앗 선물 완료
			return shim.Error(mgmt.GetCCErrMsg("CC_SSG_E02", senderWalletAddr))
		}
	
		// ---------------원장 기준 씨앗 페널티 계산(UTC TIME 기준)----------------------
		layout := "2006-01-02"
		nowDate := procDate[0:4]+"-"+procDate[5:7]+"-"+procDate[8:10]
		tmpTime, err := time.Parse(layout, nowDate)
		if err != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "procDate", err.Error()))
		}
	
		// 최대 9일까지의 데이터 체크
		var penaltySeedInfo []string
		var targetDate string
		for i := 0 ; i<10; i++ {
	
			diffDaty := tmpTime.AddDate(0,0,-i).Format(layout)
			targetDate = string(diffDaty)

			penaltykeyArgs := []string{targetDate[0:4], targetDate[5:7], targetDate[8:10], senderWalletAddr, receiverWalletAddr}
			//keyArgs := []string{formatedDateLTZ[0:4], formatedDateLTZ[5:7], formatedDateLTZ[8:10], senderWalletAddr, receiverWalletAddr}
			keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.DADIndex, penaltykeyArgs)
			if msgErr != nil {
				return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
			}
	
			// 없으면 다음 loop
			if !keyResultsIterator.HasNext() {
				continue
			}
	
			// 존재하면 해당 날짜의 UTC 및 penaltySeed 취득
			for i = 0; keyResultsIterator.HasNext(); i++ {

				responseRange, nextErr := keyResultsIterator.Next()
				if nextErr != nil {
					return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
				}
	
				_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
				if splitKeyErr != nil {
					return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", splitKeyErr.Error()))
				}

				penaltySeedInfo = append(penaltySeedInfo, keyParts[3], keyParts[4], keyParts[6], keyParts[7]) // 보낸이, 받는이, 날짜, 페널티
				break;
	
			}
	
			if len(penaltySeedInfo) > 0 { 
				break;
			}
		}

		// 코어 페널티 정보 변수 초기화
		coreSeedPenaltyFloat = 1
		if len(penaltySeedInfo) > 0 {
	
			lastPenalty , _ := strconv.ParseFloat(penaltySeedInfo[3], 64)
			seedDate := penaltySeedInfo[2]
			distanceOf := tooltip.CalculateDaysBetween(nowDate, seedDate[0:10], "2006-01-02")
	

			if distanceOf > 1 {
				coreSeedPenaltyFloat = lastPenalty + float64(distanceOf-1)*float64(0.1)
	
			} else if distanceOf == 1 {
				coreSeedPenaltyFloat = lastPenalty - float64(0.1)
			} else {
				coreSeedPenaltyFloat = lastPenalty
			}
	
			if coreSeedPenaltyFloat > 1.0 {
				coreSeedPenaltyFloat = 1.0
			} else if coreSeedPenaltyFloat < 0.1 {
				coreSeedPenaltyFloat = 0.1
			}
		}

		// param 과 core 페널티 계산 비교 - 
		if paramSeedPenaltyFloat != coreSeedPenaltyFloat {
			return shim.Error(mgmt.GetCCErrMsg("CC_SSG_E03", paramSeedPenaltyStr, fmt.Sprint(coreSeedPenaltyFloat)))
		}
	}
	

	// ---------------응답 구성----------------------
	txId := stub.GetTxID()
	walletSendSeedInfo := RES_SSG{senderWalletAddr, receiverWalletAddr, paramSeedBaseStr, paramSeedPenaltyStr, fmt.Sprint(coreSeedPenaltyFloat), procDate, txId}
	walletqueryJSON, _ := json.Marshal(walletSendSeedInfo)

	// -------------- DAD penalty 정보 컴포짓 키 생성--------------------- 
	keyArgs = append(keyArgs, receiverWalletAddr, refID, procDate, paramSeedPenaltyStr, txId)
	res := s.makeCompositeKey(stub, cst.DADIndex, keyArgs)
	if res.Status != shim.OK {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", res.Message))
	}

	logger.Info("walletSendSeedInfo = "+fmt.Sprint(walletSendSeedInfo))

	return shim.Success(walletqueryJSON)
}