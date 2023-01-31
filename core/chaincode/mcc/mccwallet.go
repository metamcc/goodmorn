package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	cst "github.com/mycreditchain/chaincode/mcc/constants"
	"github.com/mycreditchain/chaincode/mcc/utils/tooltip"
	mw "github.com/mycreditchain/chaincode/mcc/utils/wallet"


	mgmt "github.com/mycreditchain/chaincode/mcc/utils/messages"

)

var logger = shim.NewLogger("mcc")

type MCCWalletChaincode struct {

	ReceivedPaymentDataMap  map[string]int  `json:"receivedPaymentDataMap"` // 받은 지갑 누적 갯수 
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(MCCWalletChaincode))
	if err != nil {
		logger.Errorf(mgmt.GetCCErrMsg("CC_COM_E12"), err)
	}
}

func (s *MCCWalletChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	if s.ReceivedPaymentDataMap == nil {
		logger.Info("-----------------map init-----------------------") //  redis 이용해서 리스타트 시에도 정보 유지 필요
		s.ReceivedPaymentDataMap = make(map[string]int)
	}


	// 관리 데이터 최신으로 반영
	return s.initAdmMgmtData(stub, []string{"0"})
	//return shim.Success(nil)
}


// Invoke handles all future requests
// args[0] : function_name
func (s *MCCWalletChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	useYn, errStr := s.checkFncUsing(stub, []string{function})
	if errStr != "" {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E99", errStr))
	} else {
		if  useYn == false {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E11", function))
		}
	}


	if function == "createMCCWallet" {
		return s.createMCCWallet(stub, args)
	} else if function == "updateTimezoneOffset" {
		return s.updateTimezoneOffset(stub, args)
	} else if function == "getWalletInfo" {
		return s.getWalletInfo(stub, args)
	} else if function == "aliveStatusQuery" {
		return s.aliveStatusQuery(stub, args)
	} else if function == "F2F" {
		return s.procF2F(stub, args)
	} else if function == "F2M" {
		return s.procF2M(stub, args)
	} else if function == "M2F" {
		return s.procM2F(stub, args)
	} else if function == "DDW" {
		return s.calcDDW(stub, args)
	} else if function == "DDE" {
		return s.calcDDE(stub, args)
	} else if function == "SSG" {
		return s.procSSG(stub, args)
	} else if function == "SSGM" {
		return s.procSSGM(stub, args)
	} else if function == "EV_F2F" {
		return s.procEvF2F(stub, args)
	} else if function == "restructWallet" {
		return s.restructWallet(stub, args)
	} else if function == "getCoreSeedCompKey" {
		return s.getCoreSeedCompKey(stub, args)
	} else if function == "initAdmMgmtData" {
		return s.initAdmMgmtData(stub, args)
	} else if function == "payFruit" {
		return s.payFruit(stub, args)
	} else if function == "syncWallet" {
		return s.syncWallet(stub, args)
	} else if function == "makeData" {
		return s.makeData(stub, args)
	} else if function == "getCoreCompKey" {
		return s.getCoreCompKey(stub, args)
	} else if function == "delCoreCompKey" {
		return s.delCoreCompKey(stub, args)
	} 
	
	/*
	else if function == "payFruit2" {
		return s.payFruit2(stub, args)
	} 
	else if function == "getFruitCompKey" {
		return s.getFruitCompKey(stub, args)
	} else if function == "payFruit2" {
		return s.payFruit2(stub, args)
	}  else if function == "payFruit3" {
		return s.payFruit3(stub, args)
	} 
	*/
	  



	return shim.Error(mgmt.GetCCErrMsg("CC_COM_E02"))
}

// 일반 열매 지급 - 기본형
// args[0] : R - sign
// args[1] : S - sign
// args[2] : X - sign
// args[3] : Y - sign
// args[4] : refID - 검증 및 sign
// args[5] : date [2018-11-12 00:00:00] - sign
// args[6] : 서비스 코드

// args[7] : 보내는 지갑
// args[8] : 받는 지갑
// args[9] : 결제 양
func (s *MCCWalletChaincode) gFruit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) < 10 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "10"))
	}

	// 2. 파라메터 변수 할당
	var isOk bool	// 공통 성공 실패 체크 변수
	var orgRefID string // 결제취소 - 원본 큐키 기록용

	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svcPrefix := args[6] // 서비스코드

	sender := args[7]
	receiver := args[8]
	reqAmntStr := args[9]

	if len(args) > 10 {
		orgRefID = args[10] 
	}
	
	txId := stub.GetTxID()
	reqAmntInt := new(big.Int)
	reqAmntInt, isOk = reqAmntInt.SetString(reqAmntStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "reqAmnt", reqAmntStr))
	}

	// 3. 싸인 검증 
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)

	if sender != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", sender, pubKeyToCompare, X, Y))
	}

	// hash: date + senderWal + receiverWal + amountToSend
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID+procDate + sender + receiver + reqAmntStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+
			"Y : "+Y+" R : "+R+" S : "+S+" procDate : "+procDate+" sender : "+sender+" receiver : "+receiver+" amount : "+reqAmntStr)))
	}

	// 4. 기본 유효성 체크


	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}

	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	// 동일 지갑
	if sender == receiver {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E02", sender, receiver))
	}

	// 5. 비지니스 로직
	// ---------------보내는이----------------------
	// 보내는 이 지갑 조회 - 간편조회
	senderWalletSMPInfo, _, errStr := s.getWalletSMPInfo(stub, []string{"address", sender})
	if errStr != "" {
		return shim.Error(errStr)
	}
	
	// 잔액 검증
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(senderWalletSMPInfo.TotalFruit, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderWalletSMPInfo.TotalFruit))
	}

	if senderFruitBefore.Cmp(reqAmntInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderWalletSMPInfo.TotalFruit, reqAmntStr))
	}


	// 기록 구성	
	keyArgs = []string{sender, procDate[0:4], procDate[5:7], procDate[8:10], refID, receiver, "-"+reqAmntStr, txId, procDate[11:19]}
	//logger.Info(keyArgs)
	err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	// ---------------받는이----------------------
	// 기록 구성
	keyArgs = []string{receiver, procDate[0:4], procDate[5:7], procDate[8:10], refID, sender, reqAmntStr, txId, procDate[11:19] }
	err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	// 결제 지갑의 동기화 기준 count 구성 - 추후 redis 로 구성

	receiverSyncData := ""

	// 결제 지갑 토탈 갯수 조회
	totalCnt := s.ReceivedPaymentDataMap[receiver]

	logger.Info("------------------누적 받는이 cnt-----------------------")
	logger.Info(receiver)
	logger.Info(totalCnt)

	//senderDataSyncCntStr := mgmt.GetVarData(cst.SND_DATA_SYNC_CNT)
	//senderDataSyncCntInt, _ := strconv.Atoi(senderDataSyncCntStr)

	receiverDataSyncCntStr := mgmt.GetVarData(cst.RCV_DATA_SYNC_CNT)
	receiverDataSyncCntInt, _ := strconv.Atoi(receiverDataSyncCntStr)
	logger.Info("------------------동기화 기준 갯수-----------------------")
	logger.Info(receiverDataSyncCntInt)

	totalCnt = totalCnt+1
	logger.Info("------------------증가된 받는이 cnt-----------------------")
	logger.Info(totalCnt)

	// 현재를 포함한 갯수 대비 기준 갯수를 비교
	if totalCnt % (receiverDataSyncCntInt+1) == 0 {
		receiverSyncData = refID
	}

	logger.Info("------------------받는이 대상-----------------------")
	logger.Info(receiverSyncData)

	// ---------------응답 구성----------------------
	resInfo := STRT_FRUIT_PAYMENT{sender, receiver, reqAmntStr, senderWalletSMPInfo.FruitSyncData, receiverSyncData, procDate, txId,  refID, orgRefID}
	resInfoJSON, _ := json.Marshal(resInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	logger.Info("------------------코어 검증키 -----------------------")
	logger.Info(keyArgs)

	keyArgs = append(keyArgs, txId)
	err = s.createCoreCompKey(stub, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	// 받는이 지갑 갯수 증가 반영
	logger.Info("------------------최종 받는이 cnt-----------------------")
	logger.Info(totalCnt)
	s.ReceivedPaymentDataMap[receiver] = totalCnt

	return shim.Success(resInfoJSON)
}

// 일반 열매 결제 - 기본형
// args[0] : R - sign
// args[1] : S - sign
// args[2] : X - sign
// args[3] : Y - sign
// args[4] : refID - 검증 및 sign
// args[5] : date [2018-11-12 00:00:00] - sign
// args[6] : 서비스 코드

// args[7] : 보내는 지갑
// args[8] : 받는 지갑
// args[9] : 결제 양
func (s *MCCWalletChaincode) payFruit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) < 10 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "10"))
	}

	// 2. 파라메터 변수 할당
	var isOk bool	// 공통 성공 실패 체크 변수
	var orgRefID string // 결제취소 - 원본 큐키 기록용

	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svcPrefix := args[6] // 서비스코드

	sender := args[7]
	receiver := args[8]
	reqAmntStr := args[9]

	if len(args) > 10 {
		orgRefID = args[10] 
	}
	
	txId := stub.GetTxID()
	reqAmntInt := new(big.Int)
	reqAmntInt, isOk = reqAmntInt.SetString(reqAmntStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "reqAmnt", reqAmntStr))
	}

	// 3. 싸인 검증 
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)

	if sender != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", sender, pubKeyToCompare, X, Y))
	}

	// hash: date + senderWal + receiverWal + amountToSend
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID+procDate + sender + receiver + reqAmntStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+
			"Y : "+Y+" R : "+R+" S : "+S+" procDate : "+procDate+" sender : "+sender+" receiver : "+receiver+" amount : "+reqAmntStr)))
	}

	// 4. 기본 유효성 체크

	logger.Info("------------------파라메터 정보 -----------------------")
	logger.Info(procDate)
	logger.Info(svcPrefix)
	logger.Info(refID)


	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(resultCnt)


	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	// 동일 지갑
	if sender == receiver {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E02", sender, receiver))
	}

	// 5. 비지니스 로직
	// ---------------보내는이----------------------
	// 보내는 이 지갑 조회 - 간편조회
	senderWalletSMPInfo, _, errStr := s.getWalletSMPInfo(stub, []string{"address", sender})
	if errStr != "" {
		return shim.Error(errStr)
	}
	
	// 잔액 검증
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(senderWalletSMPInfo.TotalFruit, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderWalletSMPInfo.TotalFruit))
	}

	if senderFruitBefore.Cmp(reqAmntInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderWalletSMPInfo.TotalFruit, reqAmntStr))
	}


	// 기록 구성	
	keyArgs = []string{sender, procDate[0:4], procDate[5:7], procDate[8:10], refID, receiver, "-"+reqAmntStr, txId, procDate[11:19]}
	//logger.Info(keyArgs)
	err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	// ---------------받는이----------------------
	// 기록 구성
	keyArgs = []string{receiver, procDate[0:4], procDate[5:7], procDate[8:10], refID, sender, reqAmntStr, txId, procDate[11:19] }
	err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	// 결제 지갑의 동기화 기준 count 구성 - 추후 redis 로 구성

	receiverSyncData := ""

	// 결제 지갑 토탈 갯수 조회
	totalCnt := s.ReceivedPaymentDataMap[receiver]

	logger.Info("------------------누적 받는이 cnt-----------------------")
	logger.Info(receiver)
	logger.Info(totalCnt)

	//senderDataSyncCntStr := mgmt.GetVarData(cst.SND_DATA_SYNC_CNT)
	//senderDataSyncCntInt, _ := strconv.Atoi(senderDataSyncCntStr)

	receiverDataSyncCntStr := mgmt.GetVarData(cst.RCV_DATA_SYNC_CNT)
	receiverDataSyncCntInt, _ := strconv.Atoi(receiverDataSyncCntStr)
	logger.Info("------------------동기화 기준 갯수-----------------------")
	logger.Info(receiverDataSyncCntInt)

	totalCnt = totalCnt+1
	logger.Info("------------------증가된 받는이 cnt-----------------------")
	logger.Info(totalCnt)

	// 현재를 포함한 갯수 대비 기준 갯수를 비교
	if totalCnt % (receiverDataSyncCntInt+1) == 0 {
		receiverSyncData = refID
	}

	logger.Info("------------------받는이 대상-----------------------")
	logger.Info(receiverSyncData)

	// ---------------응답 구성----------------------
	resInfo := STRT_FRUIT_PAYMENT{sender, receiver, reqAmntStr, senderWalletSMPInfo.FruitSyncData, receiverSyncData, procDate, txId,  refID, orgRefID}
	resInfoJSON, _ := json.Marshal(resInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID, txId}
	err = s.createCoreCompKey(stub, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	// 받는이 지갑 갯수 증가 반영
	logger.Info("------------------최종 받는이 cnt-----------------------")
	logger.Info(totalCnt)
	s.ReceivedPaymentDataMap[receiver] = totalCnt

	return shim.Success(resInfoJSON)
}


// 열매 결제 데이터 동기화
// 열매 결제 syncData와 동일 배열 위치로 구성해서 카프카 스트림즈를 svc, syncData 유무로만 체크

// args[0] : refID - 검증 및 sign
// args[1] : date [2018-11-12 00:00:00] - sign
// args[2] : 서비스 코드

// args[3] : 동기화 지갑
// args[4] : 동기화 시작키
// args[5] : 동기화 종료키
func (s *MCCWalletChaincode) syncWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) < 6 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "6"))
	}

	// 2. 파라메터 변수 할당
	var isOk bool	// 공통 성공 실패 체크 변수
	refID := args[0]
	procDate := args[1]
	svcPrefix := args[2] // 서비스코드

	walletAddr := args[3]
	firtFruitPaymentRefID := args[4]
	lastFruitPaymentRefID := args[5]
	
	var paymentWalletYn string // 결제 지갑 구분 코드
	if len(args) > 6 {
		paymentWalletYn = args[6] 
	}
	
	txId := stub.GetTxID()

	// 4. 기본 유효성 체크
	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}


	// 5. 비지니스 로직
	// 지갑 존재 체크
	walletAsBytes, _ := stub.GetState(walletAddr)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletAddr))
	}

	// 지갑 갱신
	// 결제열매 동기화 처리 - 출금기준으로 조회된 기준의 컴포짓키 삭제
	keyArgs = []string{walletAddr} 	
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexFruitPayment, keyArgs)
	if msgErr != nil {
		// 동기화 컴포짓 데이터 취득 실패
		return shim.Error(mgmt.GetCCErrMsg("CC_SNC_E01", msgErr.Error()))
	}
	defer keyResultsIterator.Close()

	var syncSuccessMsg []RES_FRUIT_BATCH_PAYMENT_STRT
	paymentFruitTotal := new(big.Int)
	paymentFruitAmt := new(big.Int)
	paymentDataSyncCntInt := 0

	if paymentWalletYn == "Y" {
		receiverDataSyncCntStr := mgmt.GetVarData(cst.RCV_DATA_SYNC_CNT)
		receiverDataSyncCntInt, _ := strconv.Atoi(receiverDataSyncCntStr)
		paymentDataSyncCntInt = receiverDataSyncCntInt

	} else {
		senderDataSyncCntStr := mgmt.GetVarData(cst.SND_DATA_SYNC_CNT)
		senderDataSyncCntInt, _ := strconv.Atoi(senderDataSyncCntStr)
		paymentDataSyncCntInt = senderDataSyncCntInt
	}


	chkSyncCnt := 0
	syncData := make([]string, 2)
	logger.Info("walletAddr = "+walletAddr)
	logger.Info("paymentWalletYn = "+paymentWalletYn)
	logger.Info("paymentDataSyncCntInt = "+strconv.Itoa(paymentDataSyncCntInt))
	logger.Info("firtFruitPaymentRefID = " +firtFruitPaymentRefID)
	logger.Info("lastFruitPaymentRefID = "+lastFruitPaymentRefID)
	logger.Info("--------------------------------------")

	var paymentCnt int
	var keyExistChkCnt int
	for paymentCnt = 0; keyResultsIterator.HasNext(); paymentCnt++ {
		keyRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			// 동기화 대상 없음
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}
		
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(keyRange.Key)
		if splitKeyErr != nil {
				// 동기화 대상 없음
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		logger.Info("paymentCnt = "+strconv.Itoa(paymentCnt))
		logger.Info("keyParts[4]")
		logger.Info(keyParts[4])
	

		// 동기화 시작 키,종료 키 체크
		if firtFruitPaymentRefID == lastFruitPaymentRefID { // 시작키와 종료키가 같으면 - 최초
			if paymentCnt == 0 {
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("sameKey firtFruitPaymentRefID = "+keyParts[4])
				logger.Info("--------------------------------------")
			} else if paymentCnt == paymentDataSyncCntInt {
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("sameKey lastFruitPaymentRefID = "+keyParts[4])
				logger.Info("--------------------------------------")
			}

			if firtFruitPaymentRefID == keyParts[4] { // 요청키 존재
				keyExistChkCnt = 1
			} else { // 미존재시
				if paymentCnt == paymentDataSyncCntInt { // 동기화 갯수에 도달
					keyExistChkCnt = 2
				}
			}

			
		} else {
			
			if paymentCnt == 0 && keyParts[4] == firtFruitPaymentRefID { // 최초 아님
				keyExistChkCnt = 1
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("firtFruitPaymentRefID = "+firtFruitPaymentRefID)
				logger.Info("--------------------------------------")
			} else if paymentCnt == paymentDataSyncCntInt && keyParts[4] == lastFruitPaymentRefID {
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("lastFruitPaymentRefID = "+lastFruitPaymentRefID)
				logger.Info("--------------------------------------")
			}
		}
		
 
		logger.Info("keyExistChkCnt = "+strconv.Itoa(keyExistChkCnt))
		logger.Info("chkSyncCnt = "+strconv.Itoa(chkSyncCnt))

		if chkSyncCnt > 0 && chkSyncCnt <= 2 {

			// 다음 동기화 대상 loop에는 삭제 skip 로직 추가
			err = stub.DelState(keyRange.Key)
			if err != nil {
				// 동기화 대상 삭제 실패
				return shim.Error(mgmt.GetCCErrMsg("CC_SNC_E03", walletAddr))
			}

			paymentFruitAmt.SetString(keyParts[6], 10)

			logger.Info("paymentFruitAmt")
			logger.Info(paymentFruitAmt)

			paymentFruitTotal = paymentFruitTotal.Add(paymentFruitTotal, paymentFruitAmt)
			syncSuccessMsg = append(syncSuccessMsg, RES_FRUIT_BATCH_PAYMENT_STRT{keyParts[4], ""})

			//logger.Info("paymentFruitTotal")
			//logger.Info(paymentFruitTotal)
		}
		
	

		// 설정된 갯수나 컴포짓키까지만 삭제
		if chkSyncCnt == 2 {
			chkSyncCnt = 3
		}

		// 두번째 동기화 대상 시작키
		if paymentCnt == paymentDataSyncCntInt+1 {
			logger.Info("--------------------------------------")
			logger.Info("next first refID = " +keyParts[4])
			logger.Info("--------------------------------------")
			syncData[0] = keyParts[4]
		}

		// 두번째 동기화 대상 종료키
		if paymentCnt == (paymentDataSyncCntInt*2)+1 {
			logger.Info("--------------------------------------")
			logger.Info("next last refID = "+keyParts[4])
			logger.Info("--------------------------------------")
			syncData[1] = keyParts[4]
			break
		}

	}

	logger.Info("---------------최종 체크 데이터 ----------------------")
	logger.Info("lst paymentCnt = "+strconv.Itoa(paymentCnt))    // loop 돈 횟수
	logger.Info("chkSyncCnt = "+strconv.Itoa(chkSyncCnt))
	logger.Info("keyExistChkCnt = "+strconv.Itoa(keyExistChkCnt))
	logger.Info("paymentFruitTotal = "+fmt.Sprint(paymentFruitTotal))
	logger.Info("--------------------------------------")

	if chkSyncCnt != 3 || keyExistChkCnt == 0 {
		// 동기화 실패 - 대상 refID 불일치

		if paymentWalletYn == "Y"  && chkSyncCnt != 3 { // 결제 지갑의 동기화 실패 중 갯수 부정확시 원장의 갯수로 조절

			logger.Info("------------------실패시 결제 지갑 갯수 원장 기준 반영 전-----------------------")
			logger.Info(walletAddr +" = "+strconv.Itoa(s.ReceivedPaymentDataMap[walletAddr])) 

			logger.Info("------------------실패시 결제 지갑 갯수 원장 기준 반영 후-----------------------")
			s.ReceivedPaymentDataMap[walletAddr] = paymentCnt
			logger.Info(walletAddr +" = "+strconv.Itoa(s.ReceivedPaymentDataMap[walletAddr])) 
		}
		return shim.Error(mgmt.GetCCErrMsg("CC_SNC_E02", walletAddr, firtFruitPaymentRefID, lastFruitPaymentRefID, strconv.Itoa(paymentCnt), strconv.Itoa(keyExistChkCnt)))
	}

	walletStrt := STRT_WALLET{}
	json.Unmarshal(walletAsBytes, &walletStrt)

	// 데이터 계산
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(walletStrt.FruitAmount, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", walletStrt.FruitAmount))
	}

	senderFruitAfter := new(big.Int)
	senderFruitAfter = senderFruitAfter.Add(senderFruitBefore, paymentFruitTotal)
	walletStrt.FruitAmount = fmt.Sprint(senderFruitAfter)

	// 지갑 갱신
	walletAsBytes, _ = json.Marshal(walletStrt)
	err = stub.PutState(walletAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walletAddr))
	}


	// ---------------응답 구성----------------------
	resInfo := STRT_SYNC_WALLET{walletAddr, refID, fmt.Sprint(paymentFruitTotal), syncSuccessMsg,  syncData, procDate, txId}
	resInfoJSON, _ := json.Marshal(resInfo)

	// ---------------이벤트 설정--------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성-------
	keyArgs = []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID, txId}
	err = s.createCoreCompKey(stub, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	
	// 결제 지갑 동기화 완료 갯수 반영
	if paymentWalletYn == "Y" {
		
		totalCnt := s.ReceivedPaymentDataMap[walletAddr]
		chkCnt := paymentDataSyncCntInt+1
		if totalCnt >= chkCnt {
			s.ReceivedPaymentDataMap[walletAddr] = totalCnt - chkCnt
		} else {
			s.ReceivedPaymentDataMap[walletAddr] = 0
		}
		

		logger.Info("------------------ 결제 지갑 동기화 완료 반영  전 -----------------------")
		logger.Info(totalCnt)


		logger.Info("------------------ 결제 지갑 동기화 완료 반영  후 -----------------------")
		logger.Info(s.ReceivedPaymentDataMap[walletAddr])
	} 

	return shim.Success(resInfoJSON)
}


func (s *MCCWalletChaincode) syncWalletOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 6 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "6"))
	}

	// 2. 파라메터 변수 할당
	var isOk bool	// 공통 성공 실패 체크 변수
	refID := args[0]
	procDate := args[1]
	svcPrefix := args[2] // 서비스코드

	walletAddr := args[3]
	firtFruitPaymentRefID := args[4]
	lastFruitPaymentRefID := args[5]

	txId := stub.GetTxID()

	// 4. 기본 유효성 체크
	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}


	// 5. 비지니스 로직
	// 지갑 존재 체크
	walletAsBytes, _ := stub.GetState(walletAddr)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletAddr))
	}

	// 지갑 갱신
	// 결제열매 동기화 처리 - 출금기준으로 조회된 기준의 컴포짓키 삭제
	keyArgs = []string{walletAddr} 	
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexFruitPayment, keyArgs)
	if msgErr != nil {
		// 동기화 컴포짓 데이터 취득 실패
		return shim.Error(mgmt.GetCCErrMsg("CC_SNC_E01", msgErr.Error()))
	}
	defer keyResultsIterator.Close()

	var syncSuccessMsg []RES_FRUIT_BATCH_PAYMENT_STRT
	paymentFruitTotal := new(big.Int)
	paymentFruitAmt := new(big.Int)
	paymentDataSyncCntStr := mgmt.GetVarData(cst.SND_DATA_SYNC_CNT)
	paymentDataSyncCntInt, _ := strconv.Atoi(paymentDataSyncCntStr)
	chkSyncCnt := 0
	syncData := make([]string, 2)

	logger.Info("paymentDataSyncCntInt = "+strconv.Itoa(paymentDataSyncCntInt))
	logger.Info("firtFruitPaymentRefID = " +firtFruitPaymentRefID)
	logger.Info("lastFruitPaymentRefID = "+lastFruitPaymentRefID)
	logger.Info("--------------------------------------")

	var paymentCnt int
	for paymentCnt = 0; keyResultsIterator.HasNext(); paymentCnt++ {
		keyRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			// 동기화 대상 없음
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}
		
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(keyRange.Key)
		if splitKeyErr != nil {
				// 동기화 대상 없음
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		logger.Info("paymentCnt = "+strconv.Itoa(paymentCnt))
		logger.Info("keyParts[4]")
		logger.Info(keyParts[4])
	

		// 동기화 시작 키,종료 키 체크
		if firtFruitPaymentRefID == lastFruitPaymentRefID {
			if paymentCnt == 0 {
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("sameKey firtFruitPaymentRefID = "+keyParts[4])
				logger.Info("--------------------------------------")
			} else if paymentCnt == paymentDataSyncCntInt {
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("sameKey lastFruitPaymentRefID = "+keyParts[4])
				logger.Info("--------------------------------------")
			}
		} else {
			if paymentCnt == 0 && keyParts[4] == firtFruitPaymentRefID {
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("firtFruitPaymentRefID = "+firtFruitPaymentRefID)
				logger.Info("--------------------------------------")
			} else if paymentCnt == paymentDataSyncCntInt && keyParts[4] == lastFruitPaymentRefID {
				chkSyncCnt++
				logger.Info("--------------------------------------")
				logger.Info("lastFruitPaymentRefID = "+lastFruitPaymentRefID)
				logger.Info("--------------------------------------")
			}
		}
		

		logger.Info("chkSyncCnt = "+strconv.Itoa(chkSyncCnt))

		if chkSyncCnt > 0 && chkSyncCnt <= 2 {

			// 다음 동기화 대상 loop에는 삭제 skip 로직 추가
			err = stub.DelState(keyRange.Key)
			if err != nil {
				// 동기화 대상 삭제 실패
				return shim.Error(mgmt.GetCCErrMsg("CC_SNC_E03", walletAddr))
			}


			paymentFruitAmt.SetString(keyParts[6], 10)

			logger.Info("paymentFruitAmt")
			logger.Info(paymentFruitAmt)

			paymentFruitTotal = paymentFruitTotal.Add(paymentFruitTotal, paymentFruitAmt)
			syncSuccessMsg = append(syncSuccessMsg, RES_FRUIT_BATCH_PAYMENT_STRT{keyParts[4], ""})

			//logger.Info("paymentFruitTotal")
			//logger.Info(paymentFruitTotal)
		}
		

		// 설정된 갯수나 컴포짓키까지만 삭제
		if chkSyncCnt == 2 {
			chkSyncCnt = 3
		}

		// 두번째 동기화 대상 시작키
		if paymentCnt == paymentDataSyncCntInt+1 {
			logger.Info("--------------------------------------")
			logger.Info("next first refID = " +keyParts[4])
			logger.Info("--------------------------------------")
			syncData[0] = keyParts[4]

			
		}

		// 두번째 동기화 대상 종료키
		if paymentCnt == (paymentDataSyncCntInt*2)+1 {
			logger.Info("--------------------------------------")
			logger.Info("next second refID = "+keyParts[4])
			logger.Info("--------------------------------------")
			syncData[1] = keyParts[4]
			break
		}

		//paymentCnt++
	}

	logger.Info("--------------------------------------")
	logger.Info("lst paymentCnt = "+strconv.Itoa(paymentCnt))
	logger.Info("chkSyncCnt = "+strconv.Itoa(chkSyncCnt))
	logger.Info("paymentFruitTotal = "+fmt.Sprint(paymentFruitTotal))
	logger.Info("--------------------------------------")
	if chkSyncCnt != 3 {
		// 동기화 실패 - 대상 refID 불일치
		return shim.Error(mgmt.GetCCErrMsg("CC_SNC_E02", walletAddr, firtFruitPaymentRefID, lastFruitPaymentRefID, strconv.Itoa(paymentCnt), ""))
	}

	walletStrt := STRT_WALLET{}
	json.Unmarshal(walletAsBytes, &walletStrt)

	// 데이터 계산
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(walletStrt.FruitAmount, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", walletStrt.FruitAmount))
	}

	senderFruitAfter := new(big.Int)
	senderFruitAfter = senderFruitAfter.Add(senderFruitBefore, paymentFruitTotal)
	walletStrt.FruitAmount = fmt.Sprint(senderFruitAfter)

	// 지갑 갱신
	walletAsBytes, _ = json.Marshal(walletStrt)
	err = stub.PutState(walletAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walletAddr))
	}


	// ---------------응답 구성----------------------
	resInfo := STRT_SYNC_WALLET{walletAddr, refID, fmt.Sprint(paymentFruitTotal), syncSuccessMsg,  syncData, procDate, txId}
	resInfoJSON, _ := json.Marshal(resInfo)

	// ---------------이벤트 설정--------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성-------
	keyArgs = append(keyArgs, txId)
	err = s.createCoreCompKey(stub, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	return shim.Success(resInfoJSON)
}


// 잔액 검증 용 - 지갑 간편 조회
// args[0]: 조회 방식 [possible query types: identity / address]
// args[1]: walletID or walletAddress [depends on args[0] ]
func (s *MCCWalletChaincode) getWalletSMPInfo(stub shim.ChaincodeStubInterface, args []string) (RES_SMP_WALL, STRT_WALLET, string) {

	// 1. 파라메터 변수 할당
	var walletToSearch string
	queryType := args[0]
	queryData := args[1]
	myWalletStruct := RES_SMP_WALL{} 
	walletStruct := STRT_WALLET{}
	syncData := make([]string, 2)


	// 2. 지갑 조회
	if queryType == "identity" {

		walletToSearch = s.findWalletByID(stub, "UID", queryData)

		if walletToSearch == "" {
			return myWalletStruct, walletStruct, mgmt.GetCCErrMsg("CC_WLT_E01", queryData)
		}
	} else if queryType == "address" {

		walletToSearch = queryData
	} 

	walletAsBytes, _ := stub.GetState(walletToSearch)
	if walletAsBytes == nil {
		return myWalletStruct, walletStruct, mgmt.GetCCErrMsg("CC_WLT_E01", walletToSearch)
	}

	json.Unmarshal(walletAsBytes, &walletStruct)

	// 2. 결제 내역 조회
	keyArgs := []string{walletStruct.WalletAddr} 	
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexFruitPayment, keyArgs)
	if msgErr != nil {
		return myWalletStruct, walletStruct, mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error())
	}
	defer keyResultsIterator.Close()

	paymentFruitTotal := new(big.Int)
	paymentFruitAmt := new(big.Int)
	paymentDataSyncCntStr := mgmt.GetVarData(cst.SND_DATA_SYNC_CNT)
	paymentDataSyncCntInt, _ := strconv.Atoi(paymentDataSyncCntStr)

	logger.Info("--------------------------------------")
	logger.Info(" paymentDataSyncCntStr = "+paymentDataSyncCntStr)
	logger.Info("--------------------------------------")


	//4 키 5 받는이 6 금액

	//var tmpPaymentDataCntMap map[string]int
	//tmpPaymentDataCntMap = make(map[string]int)
	//ReceivedPaymentDataMap  map[string]int  `json:"receivedPaymentDataMap"` // 받은 지갑 누적 갯수 
	//s.ReceivedPaymentDataMap = make(map[string]int)

	// 3. 결제 금액 합산
	var paymentCnt int
	for paymentCnt = 0; keyResultsIterator.HasNext(); paymentCnt++ {
		responseRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			//return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
			return myWalletStruct, walletStruct, mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error())
		}

		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			//return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", splitKeyErr.Error()))
			return myWalletStruct, walletStruct, mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error())
		}

		
		logger.Info("--------------------------------------")
		logger.Info(" paymentCnt = "+strconv.Itoa(paymentCnt))
		logger.Info("--------------------------------------")
	
		if paymentCnt == 0 {
			syncData[0] = keyParts[4]
			 
			logger.Info("--------------------------------------")
			logger.Info(" FIRST REF ID = "+keyParts[4])
			logger.Info("--------------------------------------")

		} else if paymentCnt == paymentDataSyncCntInt {
			//syncData[1] = keyParts[4]
			syncData[1] = syncData[0]

			logger.Info("--------------------------------------")
			logger.Info(" LAST REF ID = "+keyParts[4])
			logger.Info("--------------------------------------")
		}

		
		// if paymentCnt <= paymentDataSyncCntInt {

		// 	tmpAddedCnt := tmpPaymentDataCntMap[keyParts[5]]
		// 	logger.Info("---------결제 지갑 임시 정보-------------------")
		// 	logger.Info(" 결제 지갑 ==== "+keyParts[5])
		// 	logger.Info(" 기존 = "+strconv.Itoa(tmpAddedCnt))
		// 	logger.Info(" 추가 = "+strconv.Itoa(tmpAddedCnt+1))
		// 	logger.Info("--------------------------------------")
			
		// 	tmpPaymentDataCntMap[keyParts[5]] = tmpAddedCnt+1
		// }

		paymentFruitAmt.SetString(keyParts[6], 10)
		paymentFruitTotal = paymentFruitTotal.Add(paymentFruitTotal, paymentFruitAmt)

	}

	// 
	myWalletStruct.WalletAddr = walletStruct.WalletAddr
	myWalletStruct.WalletType = walletStruct.WalletType
	myWalletStruct.TimezoneOffset = walletStruct.TimezoneOffset
	myWalletStruct.Identity = walletStruct.Identity
	myWalletStruct.CreationDate = walletStruct.CreationDate
	myWalletStruct.PaymentCnt = paymentCnt

	// 4. 지갑 내용 재구성
	if paymentCnt > 0 {

		walletFruitAmt := new(big.Int)
		walletFruitAmt.SetString(walletStruct.FruitAmount, 10)

		walletTotalFruit := new(big.Int)
		walletTotalFruit = walletTotalFruit.Add(paymentFruitTotal, walletFruitAmt)

		myWalletStruct.WalletFruit = walletStruct.FruitAmount
		myWalletStruct.PaymentFruit = fmt.Sprint(paymentFruitTotal)
		myWalletStruct.TotalFruit = fmt.Sprint(walletTotalFruit)

		// 동기화 데이터 존재시
		if paymentCnt >= paymentDataSyncCntInt {

			// 동기화 키
			myWalletStruct.FruitSyncData = syncData

			
			// 결제 지갑  누적 데이터 
			// for key, val := range tmpPaymentDataCntMap {

			// 	tmpAddedCnt := s.ReceivedPaymentDataMap[key]
			// 	s.ReceivedPaymentDataMap[key] = tmpAddedCnt+val

			// 	logger.Info("---------동기화 확정 결제 지갑 누적 정보-------------------")
			// 	logger.Info(" 결제 지갑 ==== "+key)
			// 	logger.Info(" 기존 = "+strconv.Itoa(tmpAddedCnt))
			// 	logger.Info(" 추가 = "+strconv.Itoa(tmpAddedCnt+val))
			// 	logger.Info("--------------------------------------")
			// }

		}
		
	} else {
		myWalletStruct.WalletFruit = walletStruct.FruitAmount
		myWalletStruct.PaymentFruit = "0"
		myWalletStruct.TotalFruit = walletStruct.FruitAmount
	}

	// 5. 반환
	return myWalletStruct, walletStruct, ""
}


// 지갑 상세 조회
// args[0]: 조회 방식 [possible query types: identity / address]
// args[1]: walletID or walletAddress [depends on args[0] ]
func (s *MCCWalletChaincode) getWalletInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "2"))
	}

	var resMsg []FRUIT_PAYMENT_STRT 
	var walletToSearch string
	queryType := args[0]
	queryData := args[1]

	if queryType == "identity" {

		walletToSearch = s.findWalletByID(stub, "UID", queryData)

		if walletToSearch == "" {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", queryData))
		}
	} else if queryType == "address" {

		walletToSearch = queryData
	} 

	walletAsBytes, _ := stub.GetState(walletToSearch)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletToSearch))
	}

	walletStruct := STRT_WALLET{}
	json.Unmarshal(walletAsBytes, &walletStruct)

	keyArgs := []string{walletStruct.WalletAddr} 	
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexFruitPayment, keyArgs)
	if msgErr != nil {
		//return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
	}
	defer keyResultsIterator.Close()

	paymentFruitTotal := new(big.Int)
	paymentFruitAmt := new(big.Int)
	var paymentCnt int
	for paymentCnt = 0; keyResultsIterator.HasNext(); paymentCnt++ {
		responseRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			//return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			//return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", splitKeyErr.Error()))
		}

		logger.Debug(keyParts)
			 
		paymentFruitAmt.SetString(keyParts[6], 10)
		paymentFruitTotal = paymentFruitTotal.Add(paymentFruitTotal, paymentFruitAmt)
	 
		resMsg = append(resMsg, FRUIT_PAYMENT_STRT{keyParts[0], keyParts[1]+"-"+keyParts[2]+"-"+keyParts[3]+" "+keyParts[8], keyParts[4], keyParts[5], keyParts[6], keyParts[7]}) 
	}



	myWalletStruct := RES_DTL_WALL{} 
	myWalletStruct.WalletAddr = walletStruct.WalletAddr
	myWalletStruct.WalletType = walletStruct.WalletType
	myWalletStruct.TimezoneOffset = walletStruct.TimezoneOffset
	myWalletStruct.Identity = walletStruct.Identity
	myWalletStruct.CreationDate = walletStruct.CreationDate
	myWalletStruct.PaymentCnt = paymentCnt

	// 지갑갱신
	if paymentCnt > 0 {

		walletFruitAmt := new(big.Int)
		walletFruitAmt.SetString(walletStruct.FruitAmount, 10)

		walletFruitTotal := new(big.Int)
		walletFruitTotal = walletFruitTotal.Add(paymentFruitTotal, walletFruitAmt)

		myWalletStruct.WalletFruit = walletStruct.FruitAmount
		myWalletStruct.PaymentFruit = fmt.Sprint(paymentFruitTotal)
		myWalletStruct.TotalFruit = fmt.Sprint(walletFruitTotal)

		myWalletStruct.PaymentData = resMsg
		
		
	} else {
		myWalletStruct.WalletFruit = walletStruct.FruitAmount
		myWalletStruct.PaymentFruit = "0"
		myWalletStruct.TotalFruit = walletStruct.FruitAmount
	}

	myWalletAsBytes, _ := json.Marshal(myWalletStruct)


	return shim.Success(myWalletAsBytes)
}


// 씨앗 선물 N:N
func (s *MCCWalletChaincode) procSSGM(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 12 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "12"))
	}


	// 2. 파라메터 변수 할당
	var msg string // 공통 메세지 처리 변수
	//var isOk bool	// 공통 성공 실패 체크 변수
	Rs := args[0]
	Ss := args[1]
	Xs := args[2]
	Ys := args[3]
	refIDs := args[4]
	procDates := args[5]
	svcs := args[6]

	senderWalletAddrs := args[7]	// 보낸이 지갑주소
	receiverWalletAddrs := args[8]	// 받는이 지갑 주소
	paramSeedBaseStrs := args[9]	// 개인별 씨앗 선물 기준 갯수
	paramSeedPenaltyStrs := args[10]	// 씨앗양 - 페널티 포함
	delimeter := args[11]	//  구분자

	// 데이터 자르기
	arrR := strings.Split(Rs, delimeter)
	arrS := strings.Split(Ss, delimeter)
	arrX := strings.Split(Xs, delimeter)
	arrY := strings.Split(Ys, delimeter)

	arrRefID := strings.Split(refIDs, delimeter)
	arrProcDate := strings.Split(procDates, delimeter)
	arrSvc := strings.Split(svcs, delimeter)
	svc := arrSvc[0]
	procDate := arrProcDate[0]
	mainRefID := "SSGM_"+procDate // 메인 ID

	arrSenderWalletAddr := strings.Split(senderWalletAddrs, delimeter)
	arrReceiverWalletAddr := strings.Split(receiverWalletAddrs, delimeter)
	arrParamSeedBaseStr := strings.Split(paramSeedBaseStrs, delimeter)
	arrParamSeedPenaltyStr := strings.Split(paramSeedPenaltyStrs, delimeter)

	var seedingSuccessMsg []RES_SSG 
	//seedStrt := RES_SSG{}

	for i, value := range arrSenderWalletAddr { 

		msg = ""

		R := arrR[i]
		S := arrS[i]
		X := arrX[i]
		Y := arrY[i]
		refID := arrRefID[i]
		procDate := arrProcDate[i]
		//svc := arrSvc[i]
	
		senderWalletAddr := value	// 보낸이 지갑주소
		receiverWalletAddr := arrReceiverWalletAddr[i]	// 받는이 지갑 주소
		paramSeedBaseStr := arrParamSeedBaseStr[i]	// 개인별 씨앗 선물 기준 갯수
		paramSeedPenaltyStr := arrParamSeedPenaltyStr[i]	// 씨앗양 - 페널티 포함

		paramSeedBaseFloat, err := strconv.ParseFloat(paramSeedBaseStr, 64) 
		if err != nil {
			msg = mgmt.GetCCErrMsg("CC_COM_E10", paramSeedBaseStr, err.Error())
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}

		paramSeedPenaltyFloat, err := strconv.ParseFloat(paramSeedPenaltyStr, 64)
		if err != nil {
			msg = mgmt.GetCCErrMsg("CC_COM_E10", paramSeedPenaltyStr, err.Error())
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}

		// 3. 싸인 검증
		if !(tooltip.IsNumeric([]string{R, S})) {
			msg = mgmt.GetCCErrMsg("CC_SGN_E01", R, S)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
		pubKey := mw.GetPubKeyFromXandY(X, Y)
		pubKeyToCompare := mw.GetPubKeyHash(pubKey)
		if senderWalletAddr != pubKeyToCompare {
			msg = mgmt.GetCCErrMsg("CC_SGN_E02", senderWalletAddr, pubKeyToCompare, X, Y)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
		
		isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID + procDate + senderWalletAddr + receiverWalletAddr + paramSeedBaseStr + paramSeedPenaltyStr)})
		if isVerified == false {
			//return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+" Y : "+Y+" R : "+R+" S : "+S+" refID : "+refID+" procDate : "+procDate+" senderWalletAddr : "+senderWalletAddr+" receiverWalletAddr : "+receiverWalletAddr+" paramSeedPenalty :" +paramSeedPenaltyStr)))
			msg = mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+" Y : "+Y+" R : "+R+" S : "+S+" refID : "+refID+" procDate : "+procDate+" senderWalletAddr : "+senderWalletAddr+" receiverWalletAddr : "+receiverWalletAddr+" paramSeedPenalty :" +paramSeedPenaltyStr))
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
		// 4. 기본 유효성 체크
		// 기본 설정 씨앗 갯수 초과 여부
		if paramSeedPenaltyFloat > paramSeedBaseFloat {
			msg = mgmt.GetCCErrMsg("CC_SSG_E04", paramSeedPenaltyStr, paramSeedBaseStr)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
		// 동일 지갑
		if senderWalletAddr == receiverWalletAddr {
			msg = mgmt.GetCCErrMsg("CC_WLT_E02", senderWalletAddr, receiverWalletAddr)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
		// 보내는이 지갑 - 존재 여부
		senderWalletAsBytes, _ := stub.GetState(senderWalletAddr)
		if senderWalletAsBytes == nil {
			msg = mgmt.GetCCErrMsg("CC_WLT_E01", senderWalletAddr)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
		// 보내는이 지갑 - 지갑 타입
		senderMCCWalletStrt := STRT_WALLET{}
		json.Unmarshal(senderWalletAsBytes, &senderMCCWalletStrt)
		if senderMCCWalletStrt.WalletType != cst.UserWallet {
			msg = mgmt.GetCCErrMsg("CC_WLT_E03", senderWalletAddr, senderMCCWalletStrt.WalletType, cst.UserWallet)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
		// 받는이 지갑 - 존재 여부
		receiverWalletAsBytes, _ := stub.GetState(receiverWalletAddr)
		if receiverWalletAsBytes == nil {
			msg = mgmt.GetCCErrMsg("CC_WLT_E01", receiverWalletAddr)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
		// 받는이 지갑 - 지갑 타입
		receiverMCCWalletStrt := STRT_WALLET{}
		json.Unmarshal(receiverWalletAsBytes, &receiverMCCWalletStrt)
		if receiverMCCWalletStrt.WalletType != cst.UserWallet {
			msg = mgmt.GetCCErrMsg("CC_WLT_E03", receiverWalletAddr, receiverMCCWalletStrt.WalletType, cst.UserWallet)
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}
	
	
		// 5. 비지니스 로직
		// DAD 기록
		chaincodeName := cst.DADLog_CHAINCODE
		channelID := stub.GetChannelID()
		ccArgs := []string{"put", refID, senderWalletAddr, receiverWalletAddr, paramSeedBaseStr, paramSeedPenaltyStr, procDate, strconv.Itoa(senderMCCWalletStrt.TimezoneOffset)}
		
		invokeArgs := util.ArrayToChaincodeArgs(ccArgs)
		dadLogRes := stub.InvokeChaincode(chaincodeName, invokeArgs, channelID)
	
		if dadLogRes.Status != shim.OK {
			if dadLogRes.Message == "" {
				msg = fmt.Sprintf("chaincode %s is not found or returns no value", chaincodeName)
			} else {
				msg = mgmt.GetCCErrMsg("CC_COM_E09") + " : " + dadLogRes.Message
			}
			
			seedingSuccessMsg = s.setErrStructDataWidthSSG(seedingSuccessMsg,senderWalletAddr,receiverWalletAddr,paramSeedBaseStr,paramSeedPenaltyStr,procDate,refID,msg)
			continue
		}

		


		//seedingSuccessMsg = dadLogRes.Payload
		seedStrt := RES_SSG{}
		json.Unmarshal(dadLogRes.Payload, &seedStrt)
		seedingSuccessMsg = append(seedingSuccessMsg, seedStrt)


		logger.Info("seedStrt = "+fmt.Sprint(seedStrt))
	
	}


	// type RES_SSG struct {
	// 	SenderAddr string `json:"senderAddr"`							// 보낸이 지갑주소
	// 	RecipientAddr string `json:"recipientAddr"`						// 받는이 지갑주소
	// 	ParamBaseSeed        	string `json:"paramBaseSeed"`	// 개인별 기준 선물 갯수
	// 	ParamPenaltySeed       		string `json:"paramPenaltySeed"`	// 파라메터 페널티 적용 선물 양
	// 	CorePenaltySeed       		string `json:"corePenaltySeed"`	// 코어 페널티 적용 선물 양
	// 	ProcDate          	string `json:"procDate"`	// 처리날짜
	// 	TxID          string `json:"txid"`				// txID
	// }
	
	
	// type RES_SSG_LST struct {
	// 	SeedingResList []RES_SSG `json:"seedingResList"`
	// }

	// // ---------------이벤트 설정----------------------
	// //eventID := fmt.Sprint(svc+"_"+refID)
	// eventID := fmt.Sprint("2_fdagfdsagsdagdsa")
	// err = stub.SetEvent(eventID, seedingSuccessMsg)
	// if err != nil {
	// 	return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
		
	// }

	// return shim.Success(seedingSuccessMsg)
	//resInfo := STRT_SYNC_WALLET{walletAddr, refID, fmt.Sprint(paymentFruitTotal), syncSuccessMsg,  syncData, procDate, txId}
	seedingResLst := RES_SSG_LST{seedingSuccessMsg}
	resInfoJSON, _ := json.Marshal(seedingResLst)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svc +"_"+mainRefID)
	err := stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}


	return shim.Success(resInfoJSON)

}


/*
함수 사용 유무 처리
args[0] : 함수이름
*/
func (s *MCCWalletChaincode) checkFncUsing(stub shim.ChaincodeStubInterface, args []string) (bool, string) {

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


// 코어 씨앗 선물 검증
// args[0] : poc date: format [2018-12-06 06:12:34] - UTC
// args[1] - 보낸이 지갑 주소
// args[2] - 받는이 지갑 주소
// args[3] : refID
func (s *MCCWalletChaincode) getCoreSeedCompKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 4 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "4"))
	}

	// 2. 파라메터 변수 할당
	procDate := args[0]
	senderWalletAddr := args[1]	// 보낸이 지갑주소
	receiverWalletAddr := args[2]	// 받는이 지갑 주소
	refID := args[3]

	var msg string // 공통 메세지 처리 변수

	// 보내는이 지갑 - 존재 여부
	senderWalletAsBytes, _ := stub.GetState(senderWalletAddr)
	if senderWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", senderWalletAddr))
	}

	// 보내는이 지갑 - 지갑 타입
	senderMCCWalletStrt := STRT_WALLET{}
	json.Unmarshal(senderWalletAsBytes, &senderMCCWalletStrt)
	
	// 5. 비지니스 로직
	// DAD 기록 조회
	formatedDateLTZ := tooltip.FormatDateUTC2LTZ(procDate, "2006-01-02 15:04:05", "2006-01-02", senderMCCWalletStrt.TimezoneOffset)

	chaincodeName := cst.DADLog_CHAINCODE
	channelID := stub.GetChannelID()

	//ccArgs := []string{"getDADCLogInfo", formatedDateLTZ[0:4], formatedDateLTZ[5:7], formatedDateLTZ[8:10], senderWalletAddr, receiverWalletAddr, refID, txID}
	ccArgs := []string{"getDADCLogInfo", formatedDateLTZ[0:4], formatedDateLTZ[5:7], formatedDateLTZ[8:10], senderWalletAddr, receiverWalletAddr, refID}
	invokeArgs := util.ArrayToChaincodeArgs(ccArgs)
	dadLogRes := stub.InvokeChaincode(chaincodeName, invokeArgs, channelID)

	if dadLogRes.Status != shim.OK {
		if dadLogRes.Message == "" {
			msg = fmt.Sprintf("chaincode %s is not found or returns no value", chaincodeName)
		}

		return shim.Error(fmt.Sprintf(mgmt.GetCCErrMsg("CC_COM_E09") + " : " + dadLogRes.Message + msg))
	}


	return shim.Success(dadLogRes.Payload)

}


// mcc -> fruit 전환 처리
// args[0] : R - sign
// args[1] : S - sign
// arga[2] : X - sign
// args[3] : Y - sign
// args[4] : refID - 검증 및 sign
// args[5] : today date: format [2018-12-06 06:12:34] - sign
// args[6] : 서비스 코드

// args[7] : 전환양 - sign
// args[8] : 이더리움입금지갑 - sign
// args[9] : 이더리움스캔해쉬 - sign
// args[10] : 사용자 지갑
func (s *MCCWalletChaincode) procM2F(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // 1. 파라메터 갯수
	if len(args) != 11 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "11"))
	}

	// 2. 변수 선언
	var isOk bool	// 공통 성공 실패 체크 변수
	
	R := args[0]
    S := args[1]
	X := args[2]
	Y := args[3]

	refID := args[4] // refID 	
	procDate := args[5] // 처리날짜	
	svcPrefix := args[6] // 서비스코드

	exchAmntStr := args[7] // 전환양 - string형	
	ethereumAddr := args[8] // 사용자이더리움입금지갑주소 
	ethereumHash := args[9] // 이더리움스캔해쉬
	userAddr := args[10] // 사용자 지갑

	// 3. 싸인 검증 
	// R,S 값 기본 검증 - 숫자형여부
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S)) 
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)

	// X,Y 값 기본 검증 - 사용자 지갑 퍼블릭 키 유효 여부
	if userAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", userAddr, pubKeyToCompare, X, Y))
	}

	// signing 유효 체크
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID + procDate + userAddr + ethereumAddr + ethereumHash + exchAmntStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("R : "+R+" S : "+S+" X : "+X+" Y : "+Y+" refID : "+refID+" procDate : "+procDate+" userAddr : "+userAddr+" ethereumAddr : "+ethereumAddr+" amount : "+exchAmntStr)))
			
	}

	// 4. 기본 유효성 체크
	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	// 사용자 지갑 존재 여부
	userWalletAsBytes, _ := stub.GetState(userAddr)
	if userWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", userAddr))
	}

	// 5. 비지니스 로직
	// ---------------사용자 지갑 처리----------------------
	userWalletStruct := STRT_WALLET{}
	json.Unmarshal(userWalletAsBytes, &userWalletStruct)

	// 기존 전체 열매 취득
	userFruitBefore := new(big.Int)
	userFruitBefore, isOk = userFruitBefore.SetString(userWalletStruct.FruitAmount, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "userAddr Fruit ", userWalletStruct.FruitAmount))
	}

	// 전환 신청 전체량
	exchAmntInt := new(big.Int)
	exchAmntInt, isOk = exchAmntInt.SetString(exchAmntStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "exchange amount", exchAmntStr))
	}

	// 기존 전체 열매에 전환 신청량 추가
	userFruitAfter := new(big.Int)
	userFruitAfter = userFruitAfter.Add(userFruitBefore, exchAmntInt)
	userWalletStruct.FruitAmount = fmt.Sprint(userFruitAfter)

	// 사용자 지갑 갱신
	userWalletAsBytes, _ = json.Marshal(userWalletStruct)
	err = stub.PutState(userAddr, userWalletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", userAddr))
	}

	// ---------------응답 구성----------------------
	txId := stub.GetTxID()
	resWalletInfo := RES_M2F{userAddr, ethereumAddr, ethereumHash, fmt.Sprint(userWalletStruct.FruitAmount), fmt.Sprint(exchAmntInt),  procDate, txId, refID}
	resWalletJSON, _ := json.Marshal(resWalletInfo)

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	coreChkCompErr := s.createCoreCompKey(stub, keyArgs)
	if coreChkCompErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
	}

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resWalletJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	return shim.Success(resWalletJSON)

}


// fruit -> mcc 전환 처리
// args[0] : R - sign
// args[1] : S - sign
// arga[2] : X - sign
// args[3] : Y - sign
// args[4] : refID - 검증 및 sign
// args[5] : today date: format [2018-12-06 06:12:34] - sign
// args[6] : 서비스 코드

// args[7] : 이더리움입금지갑 - sign
// args[8] : 사용자 지갑
// args[9] : 신청량
// args[10] : 기본수수료량 
// args[11] : 추가수수료량
func (s *MCCWalletChaincode) procF2M(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // 1. 파라메터 갯수
	if len(args) != 12 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "12"))
	}

	// 2. 변수 선언
	//var msg string // 공통 메세지 처리 변수
	var isOk bool	// 공통 성공 실패 체크 변수

	R := args[0]
    S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4] // refID 
	procDate := args[5] // 처리날짜
	svcPrefix := args[6] // 서비스코드
		
	ethereumAddr := args[7] // 이더리움입금지갑주소
	userAddr := args[8] // 사용자 지갑
	AmntStr := args[9] // 전체량 - string형
	baseFeeStr := args[10] // 기본수수료량 - string형
	extraFeeStr := args[11] // 추가수수료량 - string형

	// 3. 싸인 검증
	// R,S 값 기본 검증 - 숫자형여부
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S)) 
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)

	// X,Y 값 기본 검증 - 사용자 지갑 퍼블릭 키 유효 여부
	if userAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", userAddr, pubKeyToCompare, X, Y))
	}

	// signing 유효 체크
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID + procDate + ethereumAddr + userAddr + AmntStr+ baseFeeStr+ extraFeeStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("R : "+R+" S : "+S+" X : "+X+" Y : "+Y+" refID : "+refID+" procDate : "+procDate+" ethereumAddr : "+ethereumAddr+" userAddr : "+userAddr+" amount : "+AmntStr+" baseFee : "+baseFeeStr+" extraFee : "+extraFeeStr)))
			
	}


	// 4. 기본 유효성 체크
	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}



	// 전환 신청 전체량
	reqAmntInt := new(big.Int)
	reqAmntInt, isOk = reqAmntInt.SetString(AmntStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "amount", AmntStr))
	}

	// 기본수수료량
	baseFeeAmntInt := new(big.Int)
	baseFeeAmntInt, isOk = baseFeeAmntInt.SetString(baseFeeStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "baseFee", baseFeeStr))
	}	

	// 추가수수료량
	extraFeeAmntInt := new(big.Int)
	extraFeeAmntInt, isOk = extraFeeAmntInt.SetString(extraFeeStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "extraFee", extraFeeStr))
	}

	// 전체수수료량
	feeAmntInt := new(big.Int)
	feeAmntInt = feeAmntInt.Add(baseFeeAmntInt, extraFeeAmntInt)

	// 신청량과 수수료량 비교
	if reqAmntInt.Cmp(feeAmntInt) < 1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", AmntStr, fmt.Sprint(feeAmntInt)))
	}


	// 5. 비지니스 로직

	// ---------------사용자 지갑 처리----------------------

	/*


	// 사용자 지갑 존재 여부
	userWalletAsBytes, _ := stub.GetState(userAddr)
	if userWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", userAddr))
	}

	userWalletStruct := STRT_WALLET{}
	json.Unmarshal(userWalletAsBytes, &userWalletStruct)


	// 기존 전체 열매 취득
	userFruitBefore := new(big.Int)
	userFruitBefore, isOk = userFruitBefore.SetString(userWalletStruct.FruitAmount, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "userAddr Fruit ", userWalletStruct.FruitAmount))
	}

	if userFruitBefore.Cmp(reqAmntInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", userWalletStruct.FruitAmount, AmntStr ))
	}

	// 전체 열매에서 전환 신청량 차감
	userFruitAfter := new(big.Int)
	userFruitAfter = userFruitAfter.Sub(userFruitBefore, reqAmntInt)
	userWalletStruct.FruitAmount = fmt.Sprint(userFruitAfter)

	// 사용자 지갑 갱신
	userWalletAsBytes, _ = json.Marshal(userWalletStruct)
	err = stub.PutState(userAddr, userWalletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", userAddr))
	}
	*/

	// 전체 열매 기준 잔액 검증
	senderWalletSMPInfo, senderWalletStruct, errStr := s.getWalletSMPInfo(stub, []string{"address", userAddr})
	if errStr != "" {
		return shim.Error(errStr)
	}
	
	// 잔액 검증
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(senderWalletSMPInfo.TotalFruit, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderWalletSMPInfo.TotalFruit))
	}

	if senderFruitBefore.Cmp(reqAmntInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderWalletSMPInfo.TotalFruit, AmntStr))
	}


	// 지갑열매 기준 갱신 처리
	senderWalletFruitBefore := new(big.Int)
	senderWalletFruitBefore, isOk = senderWalletFruitBefore.SetString(senderWalletSMPInfo.WalletFruit, 10)

	senderFruitAfter := new(big.Int)
	senderFruitAfter = senderFruitAfter.Sub(senderWalletFruitBefore, reqAmntInt)

	// 지갑 데이터 구성
	senderWalletStruct.FruitAmount = fmt.Sprint(senderFruitAfter)

	// 지갑 갱신
	senderAsBytes, _ := json.Marshal(senderWalletStruct)
	err = stub.PutState(userAddr, senderAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", userAddr))
	}

	// ---------------응답 구성----------------------
	txId := stub.GetTxID()
	resWalletInfo := RES_F2M{userAddr, ethereumAddr, senderWalletStruct.FruitAmount, AmntStr, baseFeeStr, extraFeeStr, procDate, txId, refID}
	resWalletJSON, _ := json.Marshal(resWalletInfo)

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	coreChkCompErr := s.createCoreCompKey(stub, keyArgs)
	if coreChkCompErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
	}

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resWalletJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	return shim.Success(resWalletJSON)

}



// 열매 이체 
// args[0] : R - sign
// args[1] : S - sign
// args[2] : X - sign
// args[3] : Y - sign
// args[4] : refID - 검증 및 sign
// args[5] : date [2018-11-12 00:00:00] - sign
// args[6] : 서비스 코드

// args[7] : 보내는 지갑
// args[8] : 받는 지갑
// args[9] : 이체 양 
func (s *MCCWalletChaincode) procF2F(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 10 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "10"))
	}

	// 2. 파라메터 변수 할당
	var isOk bool	// 공통 성공 실패 체크 변수

	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svcPrefix := args[6] // 서비스코드

	sender := args[7]
	receiver := args[8]
	reqAmntStr := args[9]

	reqAmntInt := new(big.Int)
	reqAmntInt, isOk = reqAmntInt.SetString(reqAmntStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "reqAmnt", reqAmntStr))
	}

	// 3. 싸인 검증 
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)

	if sender != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", sender, pubKeyToCompare, X, Y))
	}

	// hash: date + senderWal + receiverWal + amountToSend
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID+procDate + sender + receiver + reqAmntStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+
			"Y : "+Y+" R : "+R+" S : "+S+" procDate : "+procDate+" sender : "+sender+" receiver : "+receiver+" amount : "+reqAmntStr)))
	}

	// 4. 기본 유효성 체크
	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	// 동일 지갑
	if sender == receiver {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E02", sender, receiver))
	}

	// 보내는 지갑
  //	senderAsBytes, _ := stub.GetState(sender)
  //	if senderAsBytes == nil {
  //		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", sender))
  //	}

	// 받는 지갑
	receiverAsBytes, _ := stub.GetState(receiver)
	if receiverAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", receiver))
	}

	// 5. 비지니스 로직

	// ---------------보내는이----------------------
  //	senderMCCWallet := STRT_WALLET{}
	//json.Unmarshal(senderAsBytes, &senderMCCWallet)

	/*
	if senderMCCWallet.WalletType != cst.UserWallet {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", sender, senderMCCWallet.WalletType, cst.UserWallet))
	}
	

	// calculate total fruit
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(senderMCCWallet.FruitAmount, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderMCCWallet.FruitAmount))
	}

	if senderFruitBefore.Cmp(reqAmntInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderMCCWallet.FruitAmount, reqAmntStr))
	}
   */

	// 전체 열매 기준 잔액 검증
	senderWalletSMPInfo, senderWalletStruct, errStr := s.getWalletSMPInfo(stub, []string{"address", sender})
	if errStr != "" {
		return shim.Error(errStr)
	}
	
	// 잔액 검증
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(senderWalletSMPInfo.TotalFruit, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderWalletSMPInfo.TotalFruit))
	}

	if senderFruitBefore.Cmp(reqAmntInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderWalletSMPInfo.TotalFruit, reqAmntStr))
	}


	// 지갑열매 기준 갱신 처리
	senderWalletFruitBefore := new(big.Int)
	senderWalletFruitBefore, isOk = senderWalletFruitBefore.SetString(senderWalletSMPInfo.WalletFruit, 10)


	senderFruitAfter := new(big.Int)
	senderFruitAfter = senderFruitAfter.Sub(senderWalletFruitBefore, reqAmntInt)

	// 지갑 데이터 구성
	senderWalletStruct.FruitAmount = fmt.Sprint(senderFruitAfter)

	// 지갑 갱신
	senderAsBytes, _ := json.Marshal(senderWalletStruct)
	err = stub.PutState(sender, senderAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", sender))
	}

	// ---------------받는이----------------------
	receiverMCCWallet := STRT_WALLET{}
	json.Unmarshal(receiverAsBytes, &receiverMCCWallet)

	/*
	if receiverMCCWallet.WalletType != cst.UserWallet {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", receiver, receiverMCCWallet.WalletType, cst.UserWallet))
	}
	*/

	receiverFruit := new(big.Int)
	receiverFruit, isOk = receiverFruit.SetString(receiverMCCWallet.FruitAmount, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "receiverFruit", receiverMCCWallet.FruitAmount))
	}

	// 열매 추가
	receiverNewFruit := new(big.Int)
	receiverNewFruit = receiverNewFruit.Add(receiverFruit, reqAmntInt)
	receiverMCCWallet.FruitAmount = fmt.Sprint(receiverNewFruit)

	// 지갑갱신
	receiverAsBytes, _ = json.Marshal(receiverMCCWallet)
	err = stub.PutState(receiver, receiverAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", receiver))
	}


	// ---------------응답 구성----------------------
	txId := stub.GetTxID()
	walletSendInfo := RES_F2F{sender, senderWalletStruct.FruitAmount, reqAmntStr, receiver, procDate, txId, refID}
	walletqueryJSON, _ := json.Marshal(walletSendInfo)


	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, walletqueryJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	coreChkCompErr := s.createCoreCompKey(stub, keyArgs)
	if coreChkCompErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
	}

	return shim.Success(walletqueryJSON)
}


// 시스템 열매 선물
// args[0] : R - sign
// args[1] : S - sign
// args[2] : X - sign
// args[3] : Y - sign
// args[4] : refID - 검증 및 sign
// args[5] : date [2018-11-12 00:00:00] - sign
// args[6] : 서비스 코드
// args[7] : 이벤트 코드 

// ev0001 : 이벤트 코드 - 이벤트 열매 선물
// args[8] : 보내는 지갑
// args[9] : 받는 지갑
// args[10] : 이체 양 

func (s *MCCWalletChaincode) procEvF2F(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 11 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "11"))
	}

	// 2. 파라메터 변수 할당
	var isOk bool	// 공통 성공 실패 체크 변수
	var msg string // 공통 메세지 처리 변수
	
	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svcPrefix := args[6] // 서비스코드

	evtCd := args[7] // 이벤트코드

	mrktWal := args[8]
	userWal := args[9]
	giftAmntStr := args[10]


	txId := stub.GetTxID()
	var walletqueryJSON []byte

	// 3. 싸인 검증 
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)

	if userWal != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", userWal, pubKeyToCompare, X, Y))
	}

	// hash: date + senderWal + receiverWal + amountToSend
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID+procDate + evtCd+ mrktWal + userWal + giftAmntStr)})
	if isVerified == false {
		
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+
			"Y : "+Y+" R : "+R+" S : "+S+" procDate : "+procDate+" evtCd : "+evtCd+" mrktWal : "+mrktWal+" userWal : "+userWal+" giftAmnt : "+giftAmntStr)))
			
	}


	// 4. 기본 유효성 체크
	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	// ---------------이벤트 정보 취득 ----------------------
	eventStrt, errStr := s.getEventConfData(stub, []string{"getCoreEvt", evtCd})
	if errStr != "" {
		return shim.Error(fmt.Sprintf(mgmt.GetCCErrMsg("CC_COM_E09") + " : " + errStr))
	}

	// 이벤트별 비지니스 로직
	if evtCd == "ev0001" {

		eventDate := eventStrt.EvtDt
		evtKeyArgs := []string{eventDate[0:4], eventDate[5:7], eventDate[8:10], evtCd, userWal} // refID가 아닌 받는 사용자 지갑을 키로 구성
		resultCnt, err := s.getCoreCompData(stub, evtKeyArgs)
		if err != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", err.Error()))
		}

		// 이벤트 참여 횟수 검증
		evtCnt := eventStrt.EvtCnt
		if evtCnt > 0 {
			if resultCnt >= evtCnt {
				resultCntStr := strconv.Itoa(resultCnt)
				return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", evtCd+"_"+resultCntStr+"_"+userWal))
			}
		}

		// 동일 지갑
		if mrktWal == userWal {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E02", mrktWal, userWal))
		}

		// 보내는 지갑
		mrktWalAsBytes, _ := stub.GetState(mrktWal)
		if mrktWalAsBytes == nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", mrktWal))
		}

		mrkWallet := STRT_WALLET{}
		json.Unmarshal(mrktWalAsBytes, &mrkWallet)

		// 시스템 지갑 여부
		if mrkWallet.WalletType != cst.SystemWallet {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", mrktWal, mrkWallet.WalletType, cst.SystemWallet))
		}

		// 받는 지갑
		userWalAsBytes, _ := stub.GetState(userWal)
		if userWalAsBytes == nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", userWal))
		}

		userWallet := STRT_WALLET{}
		json.Unmarshal(userWalAsBytes, &userWallet)


		// 5. 비지니스 로직
		// ---------------오늘 보낸 씨앗 횟수----------------------
		chaincodeName := cst.DADLog_CHAINCODE
		formatedDateLTZ := tooltip.FormatDateUTC2LTZ(procDate, "2006-01-02 15:04:05", "2006-01-02", userWallet.TimezoneOffset)
		ccArgs := []string{"getDADCLogInfo", formatedDateLTZ[0:4], formatedDateLTZ[5:7], formatedDateLTZ[8:10], userWal}
		invokeArgs := util.ArrayToChaincodeArgs(ccArgs)
		dadLogRes := stub.InvokeChaincode(chaincodeName, invokeArgs, stub.GetChannelID())

		if dadLogRes.Status != shim.OK {
			if dadLogRes.Message == "" {
				msg = fmt.Sprintf("chaincode %s is not found or returns no value", chaincodeName)
			}
	
			return shim.Error(fmt.Sprintf(mgmt.GetCCErrMsg("CC_COM_E09") + " : " + dadLogRes.Message + msg))
		}

		seedGiftInfo := RES_CC_LST{}
		json.Unmarshal(dadLogRes.Payload, &seedGiftInfo)

		// ---------------일일 씨앗 선물 최대 횟수 ----------------------
		seedLimitStr := mgmt.GetVarData(cst.SEED_LIMIT_CNT) 
		seedLimitInt, _ := strconv.Atoi(seedLimitStr)
		if len(seedGiftInfo.CoreChkInfo) != seedLimitInt {
			dadlogCntStr := strconv.Itoa(len(seedGiftInfo.CoreChkInfo))
			return shim.Error(mgmt.GetCCErrMsg("CC_SSG_E05", userWal, seedLimitStr, dadlogCntStr))
		}

		if giftAmntStr != eventStrt.EvtAmnt {
			return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E02", evtCd, eventStrt.EvtAmnt, giftAmntStr))
		}
		
		// 선물 열매량
		giftAmntBigInt := new(big.Int)
		giftAmntBigInt, isOk = giftAmntBigInt.SetString(giftAmntStr, 10)
		if !isOk {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "giftAmntBigInt", giftAmntStr))
		}

		// 마케팅 기존 열매 량
		mrktWalAmntBefore := new(big.Int)
		mrktWalAmntBefore, isOk = mrktWalAmntBefore.SetString(mrkWallet.FruitAmount, 10)
		if !isOk {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "mrktWalAmntBefore ", mrkWallet.FruitAmount))
		}

		if mrktWalAmntBefore.Cmp(giftAmntBigInt) == -1 {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", mrkWallet.FruitAmount, giftAmntStr))
		}

		// 마케팅 신규 열매 량
		mrktWalAmntAfter := new(big.Int)
		mrktWalAmntAfter.Sub(mrktWalAmntBefore, giftAmntBigInt)
		mrkWallet.FruitAmount = fmt.Sprint(mrktWalAmntAfter)
		
		// 마케팅 지갑 갱신
		mrktWalAsBytes, _ = json.Marshal(mrkWallet)
		err = stub.PutState(mrktWal, mrktWalAsBytes)
		if err != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", mrktWal))
		}

		// 사용자 기존 열매량 
		userWalFruitAmntBefore := new(big.Int)
		userWalFruitAmntBefore, _ = userWalFruitAmntBefore.SetString(userWallet.FruitAmount, 10)
		if !isOk {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "userWalFruitAmntBefore", userWallet.FruitAmount))
		}
		// 사용자 신규 열매 량
		userWalFruitAmntAfter := new(big.Int)
		userWalFruitAmntAfter = userWalFruitAmntAfter.Add(userWalFruitAmntBefore, giftAmntBigInt)
		userWallet.FruitAmount = fmt.Sprint(userWalFruitAmntAfter)

		// 사용자 지갑 갱신
		userWalAsBytes, _ = json.Marshal(userWallet)
		err = stub.PutState(userWal, userWalAsBytes)
		if err != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", userWal))
		}

		// ---------------응답 구성----------------------
		walletSendSeedInfo := RES_EV_F2F{userWal, userWallet.FruitAmount, mrktWal, giftAmntStr, procDate, txId, refID}
		walletqueryJSON, _ = json.Marshal(walletSendSeedInfo)

		// ---------------이벤트 검증용 컴포짓 키 구성----------------------
		evtKeyArgs = append(evtKeyArgs, refID)
		coreChkCompErr := s.createCoreCompKey(stub, evtKeyArgs)
		if coreChkCompErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
		}
	}

	if walletqueryJSON == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", userWal))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	coreChkCompErr := s.createCoreCompKey(stub, keyArgs)
	if coreChkCompErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
	}

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix + "_"+refID)
	err = stub.SetEvent(eventID, walletqueryJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	return shim.Success(walletqueryJSON)
}

// 개별 사용자 지갑 정산 - 2019-03-11 코어검증 및 이벤트 설정 일부 추가
// args[0] : R
// args[1] : S
// args[2] : X
// args[3] : Y
// args[4] : refID
// args[5] : 처리날짜,format must be 2018-10-08 00:00:00
// args[6] : 서비스 코드

// args[7] : 사용자 지갑 주소
// args[8] : 사용자 받은 씨앗 - 페널티 포함
// args[9] : 당일 전체 씨앗  - 페널티 포함
func (s *MCCWalletChaincode) calcDDW(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 10 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "10"))
	}

	// 2. 파라메터 변수 할당
	//var msg string // 공통 메세지 처리 변수
	var isOk bool	// 공통 성공 실패 체크 변수

	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svcPrefix := args[6] // 서비스코드

	walletAddr := args[7]
	penaltySeedStr := args[8]
	penaltySeedSumStr := args[9]

	// 3. 싸인 검증
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)
	if walletAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", walletAddr, pubKeyToCompare, X, Y))
	}

	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID+procDate + walletAddr + penaltySeedStr + penaltySeedSumStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+
			"Y : "+Y+" R : "+R+" S : "+S+" procDate : "+procDate+" walletAddr : "+walletAddr+" penaltySeedStr : "+penaltySeedStr+" penaltySeedSumStr : "+penaltySeedSumStr)))
	}

	// 4. 기본 유효성 체크
	walletAsBytes, _ := stub.GetState(walletAddr)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletAddr))
	}

	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	// 5. 비지니스 로직
	// ---------------사용자 지급 갯수(페널티 적용분) 및 전체 DAD 갯수----------------------
	userPenalty := float64(0.0)
	userPenalty, err = strconv.ParseFloat(penaltySeedStr, 64)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "penaltySeed ", penaltySeedStr))
	}
	userPenaltyBigFloat := big.NewFloat(userPenalty)

	dadSum := new(big.Float)
	dadSum, isOk = dadSum.SetString(penaltySeedSumStr)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "penaltySeedSum ", penaltySeedSumStr))
	}

	// ---------------DAD 지급 최대량----------------------	
	// DAD 량
	dadLimitAmntStr := mgmt.GetVarData(cst.DAD_LIMIT_AMNT)
	dadLimitBigFloat := new(big.Float)
	dadLimitBigFloat, isOk = dadLimitBigFloat.SetString(dadLimitAmntStr)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "dadLimit", dadLimitAmntStr))
	}

	// ---------------열매 계산 및 사용자 지갑 갱신----------------------
	settleToken := new(big.Float)
	settleToken = settleToken.Quo(userPenaltyBigFloat, dadSum)

	settleTokenBigFloat := new(big.Float)
	settleTokenBigFloat = settleTokenBigFloat.Mul(settleToken, dadLimitBigFloat)

	wallet := STRT_WALLET{}
	json.Unmarshal(walletAsBytes, &wallet)

	if wallet.WalletType != cst.UserWallet {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", walletAddr, wallet.WalletType, cst.UserWallet))
	}

	walletFruit := wallet.FruitAmount

	walFruit := new(big.Float)
	walFruit, isOk = walFruit.SetString(walletFruit)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "walletFruit", walletFruit))
	}
	newWalFruitAmount := new(big.Float)
	newWalFruitAmount = newWalFruitAmount.Add(walFruit, settleTokenBigFloat)
	wallet.FruitAmount = fmt.Sprintf("%.0f", newWalFruitAmount)

	walletAsBytes, _ = json.Marshal(wallet)
	err = stub.PutState(walletAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walletAddr))
	}


	// ---------------응답 구성----------------------
	txId := stub.GetTxID()
	walletSettleDADInfo := RES_DAD_WAL{walletAddr, wallet.FruitAmount, penaltySeedStr, penaltySeedSumStr, fmt.Sprintf("%.0f", settleTokenBigFloat), procDate, txId, refID}
	walletqueryJSON, _ := json.Marshal(walletSettleDADInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix + "_"+refID)
	err = stub.SetEvent(eventID, walletqueryJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	coreChkCompErr := s.createCoreCompKey(stub, keyArgs)
	if coreChkCompErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
	}

	return shim.Success(walletqueryJSON)
}



// 개별 사용자 지갑 정산 - 2019-03-11 코어검증 및 이벤트 설정 일부 추가
// args[0] : R
// args[1] : S
// args[2] : X
// args[3] : Y
// args[4] : refID
// args[5] : 처리날짜,format must be 2018-10-08 00:00:00
// args[6] : 서비스 코드

// args[7] : 사용자 지갑 주소
// args[8] : 사용자 받은 씨앗 - 페널티 포함
// args[9] : 당일 전체 씨앗  - 페널티 포함
func (s *MCCWalletChaincode) calcDDWM(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 10 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "10"))
	}

	// 2. 파라메터 변수 할당
	//var msg string // 공통 메세지 처리 변수
	var isOk bool	// 공통 성공 실패 체크 변수

	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svcPrefix := args[6] // 서비스코드

	walletAddr := args[7]
	penaltySeedStr := args[8]
	penaltySeedSumStr := args[9]

	// 3. 싸인 검증
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)
	if walletAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", walletAddr, pubKeyToCompare, X, Y))
	}

	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID+procDate + walletAddr + penaltySeedStr + penaltySeedSumStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+
			"Y : "+Y+" R : "+R+" S : "+S+" procDate : "+procDate+" walletAddr : "+walletAddr+" penaltySeedStr : "+penaltySeedStr+" penaltySeedSumStr : "+penaltySeedSumStr)))
	}

	// 4. 기본 유효성 체크
	walletAsBytes, _ := stub.GetState(walletAddr)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletAddr))
	}

	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	// 5. 비지니스 로직
	// ---------------사용자 지급 갯수(페널티 적용분) 및 전체 DAD 갯수----------------------
	userPenalty := float64(0.0)
	userPenalty, err = strconv.ParseFloat(penaltySeedStr, 64)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "penaltySeed ", penaltySeedStr))
	}
	userPenaltyBigFloat := big.NewFloat(userPenalty)

	dadSum := new(big.Float)
	dadSum, isOk = dadSum.SetString(penaltySeedSumStr)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "penaltySeedSum ", penaltySeedSumStr))
	}

	// ---------------DAD 지급 최대량----------------------	
	// DAD 량
	dadLimitAmntStr := mgmt.GetVarData(cst.DAD_LIMIT_AMNT)
	dadLimitBigFloat := new(big.Float)
	dadLimitBigFloat, isOk = dadLimitBigFloat.SetString(dadLimitAmntStr)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "dadLimit", dadLimitAmntStr))
	}

	// ---------------열매 계산 및 사용자 지갑 갱신----------------------
	settleToken := new(big.Float)
	settleToken = settleToken.Quo(userPenaltyBigFloat, dadSum)

	settleTokenBigFloat := new(big.Float)
	settleTokenBigFloat = settleTokenBigFloat.Mul(settleToken, dadLimitBigFloat)

	wallet := STRT_WALLET{}
	json.Unmarshal(walletAsBytes, &wallet)

	if wallet.WalletType != cst.UserWallet {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", walletAddr, wallet.WalletType, cst.UserWallet))
	}

	walletFruit := wallet.FruitAmount

	walFruit := new(big.Float)
	walFruit, isOk = walFruit.SetString(walletFruit)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "walletFruit", walletFruit))
	}
	newWalFruitAmount := new(big.Float)
	newWalFruitAmount = newWalFruitAmount.Add(walFruit, settleTokenBigFloat)
	wallet.FruitAmount = fmt.Sprintf("%.0f", newWalFruitAmount)

	walletAsBytes, _ = json.Marshal(wallet)
	err = stub.PutState(walletAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walletAddr))
	}


	// ---------------응답 구성----------------------
	txId := stub.GetTxID()
	walletSettleDADInfo := RES_DAD_WAL{walletAddr, wallet.FruitAmount, penaltySeedStr, penaltySeedSumStr, fmt.Sprintf("%.0f", settleTokenBigFloat), procDate, txId, refID}
	walletqueryJSON, _ := json.Marshal(walletSettleDADInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix + "_"+refID)
	err = stub.SetEvent(eventID, walletqueryJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	coreChkCompErr := s.createCoreCompKey(stub, keyArgs)
	if coreChkCompErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
	}

	return shim.Success(walletqueryJSON)
}


// 에코 지갑 정산- 2019-03-12 코어검증 및 이벤트 설정 일부 추가
// args[0] : R
// args[1] : S
// args[2] : X
// args[3] : Y
// args[4] : refID
// args[5] : 처리날짜,format must be 2018-10-08 00:00:00
// args[6] : 서비스 코드

// args[7] - eco 지갑 id
// args[8] - eco 지갑 주소
// args[9] - 정산
func (s *MCCWalletChaincode) calcDDE(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 10 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "10"))
	}

	// 2. 파라메터 변수 할당
	//var msg string // 공통 메세지 처리 변수
	//var isOk bool	// 공통 성공 실패 체크 변수

	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svcPrefix := args[6] // 서비스코드

	ecoID := args[7]
	ecoAddr := args[8]
	amntStr := args[9]

	// daily 정산 총갯수
	amntBigInt := new(big.Int)
	amntBigInt, _ = amntBigInt.SetString(amntStr, 10)

	// 3. 싸인 검증
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)
	if ecoAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", ecoAddr, pubKeyToCompare, X, Y))
	}

	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID + procDate + ecoID + ecoAddr + amntStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+" Y : "+Y+" R : "+R+" S : "+S+" ecoID : "+ecoID+" ecoAddr : "+ecoAddr+" amount : "+amntStr+" procDate : "+procDate)))
	}

	// 4. 기본 유효성 체크
	// 중복 처리 방지 - core 검증
	keyArgs := []string{procDate[0:4], procDate[5:7], procDate[8:10], svcPrefix, refID}
	resultCnt, err := s.getCoreCompData(stub, keyArgs)
	if err != nil {
		return shim.Error(err.Error())
	}
	if resultCnt != 0 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COR_E01", refID))
	}

	ecoWalletAsBytes, _ := stub.GetState(ecoAddr)
	if ecoWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", ecoAddr))
	}

	ecoWallet := STRT_WALLET{}
	json.Unmarshal(ecoWalletAsBytes, &ecoWallet)

	if ecoWallet.WalletType != cst.EcoWallet {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", ecoAddr, ecoWallet.WalletType, cst.EcoWallet))

	}


	// 5. 비니지스 로직
	//---------------에코 지갑 DAD 계산 및 갱신----------------------
	ecoAmntBefore := new(big.Int)
	ecoAmntBefore, _ = ecoAmntBefore.SetString(ecoWallet.FruitAmount, 10)

	if ecoAmntBefore.Cmp(amntBigInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08",  ecoWallet.FruitAmount, amntStr))
	}

	ecoAmntAfter := new(big.Int)
	ecoAmntAfter = ecoAmntAfter.Sub(ecoAmntBefore, amntBigInt)
	ecoWallet.FruitAmount = fmt.Sprint(ecoAmntAfter)

	ecoWalletAsBytes, _ = json.Marshal(ecoWallet)
	err = stub.PutState(ecoAddr, ecoWalletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", ecoAddr))
	}

	// ---------------응답 구성----------------------
	txId := stub.GetTxID()
	settleEcoInfo := RES_DAD_ECO{ecoAddr, fmt.Sprint(ecoWallet.FruitAmount), fmt.Sprint(amntBigInt), procDate, txId, refID}
	settleEcoInfoJSON, _ := json.Marshal(settleEcoInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix + "_"+refID)
	err = stub.SetEvent(eventID, settleEcoInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	coreChkCompErr := s.createCoreCompKey(stub, keyArgs)
	if coreChkCompErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", coreChkCompErr.Error()))
	}

	return shim.Success(settleEcoInfoJSON)
}


// 씨앗 선물 - 2019.03.13 씨앗 갯수 페널티 적용으로 수정 및 변경
// args[0] : R
// args[1] : S
// args[2] : X
// args[3] : Y
// args[4] : refID
// args[5] : 처리날짜,format must be 2018-10-08 00:00:00
// args[6] : 서비스코드

// args[7] - 보낸이 지갑 주소
// args[8] - 받는이 지갑 주소
// args[9] - 개인별 씨앗 선물 기준 갯수
// args[10] - 씨앗양 - 페널티 포함
 
func (s *MCCWalletChaincode) procSSG(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 11 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "11"))
	}

	// 2. 파라메터 변수 할당
	var msg string // 공통 메세지 처리 변수
	//var isOk bool	// 공통 성공 실패 체크 변수
	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	svc := args[6]

	senderWalletAddr := args[7]	// 보낸이 지갑주소
	receiverWalletAddr := args[8]	// 받는이 지갑 주소
	paramSeedBaseStr := args[9]	// 개인별 씨앗 선물 기준 갯수
	paramSeedPenaltyStr := args[10]	// 씨앗양 - 페널티 포함
	
	paramSeedBaseFloat, err := strconv.ParseFloat(paramSeedBaseStr, 64) 
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", paramSeedBaseStr, err.Error()))
	}
	paramSeedPenaltyFloat, err := strconv.ParseFloat(paramSeedPenaltyStr, 64)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", paramSeedPenaltyStr, err.Error()))

	}
	
	// 3. 싸인 검증
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)
	if senderWalletAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", senderWalletAddr, pubKeyToCompare, X, Y))
	}
	
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID + procDate + senderWalletAddr + receiverWalletAddr + paramSeedBaseStr + paramSeedPenaltyStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+" Y : "+Y+" R : "+R+" S : "+S+" refID : "+refID+" procDate : "+procDate+" senderWalletAddr : "+senderWalletAddr+" receiverWalletAddr : "+receiverWalletAddr+" paramSeedPenalty :" +paramSeedPenaltyStr)))
	}

	// 4. 기본 유효성 체크
	// 기본 설정 씨앗 갯수 초과 여부
	if paramSeedPenaltyFloat > paramSeedBaseFloat {
		return shim.Error(mgmt.GetCCErrMsg("CC_SSG_E04", paramSeedPenaltyStr, paramSeedBaseStr))
	}

	// 동일 지갑
	if senderWalletAddr == receiverWalletAddr {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E02", senderWalletAddr, receiverWalletAddr))
	}

	// 보내는이 지갑 - 존재 여부
	senderWalletAsBytes, _ := stub.GetState(senderWalletAddr)
	if senderWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", senderWalletAddr))
	}

	// 보내는이 지갑 - 지갑 타입
	senderMCCWalletStrt := STRT_WALLET{}
	json.Unmarshal(senderWalletAsBytes, &senderMCCWalletStrt)
	if senderMCCWalletStrt.WalletType != cst.UserWallet {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", senderWalletAddr, senderMCCWalletStrt.WalletType, cst.UserWallet))
	}

	// 받는이 지갑 - 존재 여부
	receiverWalletAsBytes, _ := stub.GetState(receiverWalletAddr)
	if receiverWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", receiverWalletAddr))
	}

	// 받는이 지갑 - 지갑 타입
	receiverMCCWalletStrt := STRT_WALLET{}
	json.Unmarshal(receiverWalletAsBytes, &receiverMCCWalletStrt)
	if receiverMCCWalletStrt.WalletType != cst.UserWallet {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E03", receiverWalletAddr, receiverMCCWalletStrt.WalletType, cst.UserWallet))
	}

	// 5. 비지니스 로직
	// DAD 기록
	chaincodeName := cst.DADLog_CHAINCODE
	channelID := stub.GetChannelID()
	ccArgs := []string{"put", refID, senderWalletAddr, receiverWalletAddr, paramSeedBaseStr, paramSeedPenaltyStr, procDate, strconv.Itoa(senderMCCWalletStrt.TimezoneOffset)}
	invokeArgs := util.ArrayToChaincodeArgs(ccArgs)
	dadLogRes := stub.InvokeChaincode(chaincodeName, invokeArgs, channelID)

	if dadLogRes.Status != shim.OK {
		if dadLogRes.Message == "" {
			msg = fmt.Sprintf("chaincode %s is not found or returns no value", chaincodeName)
		}

		return shim.Error(fmt.Sprintf(mgmt.GetCCErrMsg("CC_COM_E09") + " : " + dadLogRes.Message + msg))
	}

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svc+"_"+refID)
	err = stub.SetEvent(eventID, dadLogRes.Payload)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	return shim.Success(dadLogRes.Payload)
}


// 지갑 조회
// args[0]: 조회 방식 [possible query types: identity / address]
// args[1]: walletID or walletAddress [depends on args[0] ]
func (s *MCCWalletChaincode) getWalletInfoOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "2"))
	}

	var walletToSearch string
	queryType := args[0]
	queryData := args[1]

	if queryType == "identity" {

		walletToSearch = s.findWalletByID(stub, "UID", queryData)

		if walletToSearch == "" {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", queryData))
		}
	} else if queryType == "address" {

		walletToSearch = queryData
	} 

	walletAsBytes, _ := stub.GetState(walletToSearch)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletToSearch))
	}

	return shim.Success(walletAsBytes)
}


/*
체인 코드 체크용
*/
func (s *MCCWalletChaincode) aliveStatusQuery(stub shim.ChaincodeStubInterface, args []string ) pb.Response {

	refID := args[0]
	server := args[1] 

	aliveStatusInfo := RES_ALIVE{server, string(time.Now().UTC().Format("2006-01-02 15:04:05")), refID, refID}
	aliveStatusJSON, _ := json.Marshal(aliveStatusInfo)

	return shim.Success(aliveStatusJSON)
}


/*
지갑 identity 중복 체크용
*/
func (s *MCCWalletChaincode) findWalletByID(stub shim.ChaincodeStubInterface, keytype string, userid string) string {

	if len(strings.TrimSpace(userid)) == 0 {
		return mgmt.GetCCErrMsg("CC_COM_E01", "1")
	}

	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexUserID, []string{keytype, userid})
	if msgErr != nil {
		return ""
	}
	defer keyResultsIterator.Close()

	if !keyResultsIterator.HasNext() {

		return ""
	}

	var i int
	for i = 0; keyResultsIterator.HasNext(); i++ {
		responseRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			return ""
		}

		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return ""
		}

		return keyParts[2]
	}

	return ""
}

// 지갑 생성
// args[0] : user_identity
// args[1] : wallet type [  user / system  / eco  ]
// args[2] : timezone offset
func (s *MCCWalletChaincode) createMCCWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 3 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "3"))
	}

	// 2. 파라메터 변수 할당
	identity := args[0]
	walletType := args[1]
	timezoneOffsetStr := args[2]

	timezoneOffset, err := strconv.Atoi(timezoneOffsetStr)
	if err != nil {
		timezoneOffset = 0
	}
	amntStr := "0" // 최초 열매값
	chkID := identity // id 중복 체크
	if walletType == cst.EcoWallet {

		// 에코 지갑 기본 열매량
		amntStr = "230000000000000000000000000"

		chkID =   cst.EcoWallet

		// 에코 지갑 체크
		walletExists := s.findWalletByID(stub, "UID", chkID)
		if walletExists != "" {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E04"))
		}
	}

	// identity 체크
	walletExists := s.findWalletByID(stub, "UID", identity)
	if walletExists != "" {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E05", identity))
	}

	// 지갑 주소 생성
	w, _ := mw.GenerateWallet()
	pubKey := w.Publickey // 공개키
	privKey := w.Privatekey // 비밀키
	walAddr := mw.GetPubKeyHash(pubKey) // 지갑 주소

	// 지갑 생성일
	currentDate := time.Now().UTC()
	walletCreationDate := string(currentDate.Format("2006-01-02 15:04:05"))

	// 지갑 구조체 
	var wallet = STRT_WALLET{
		WalletAddr:  walAddr,
		WalletType:     walletType,
		TimezoneOffset: timezoneOffset,
		FruitAmount:    amntStr,
		Identity:       identity,
		CreationDate:   walletCreationDate,			
		}

	walletAsBytes, _ := json.Marshal(wallet)

	// 지갑 등록
	err = stub.PutState(wallet.WalletAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walAddr))
	}

	// 유니크 체크용 컴포짓 키 구성
	compositeKey, compositeErr := stub.CreateCompositeKey(cst.IndexUserID, []string{"UID", identity, wallet.WalletAddr})
	if compositeErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositeKey+" : "+compositeErr.Error()))
	}

	// 컴포짓 키 등록
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositePutErr.Error()))
	}

	// 에코 지갑 체크용 데이터 구성
	if walletType == cst.EcoWallet {

		// 유니크 체크용 컴포짓 키 구성
		compositeKey, compositeErr := stub.CreateCompositeKey(cst.IndexUserID, []string{"UID", chkID, wallet.WalletAddr})
		if compositeErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositeErr.Error()))
		}

		// 컴포짓 키 등록
		compositePutErr := stub.PutState(compositeKey, []byte{0x00})
		if compositePutErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositeKey+" : "+compositeErr.Error()))
		}
	}

	// 응답 데이터 - 비밀키 포함
	walletInfo := RES_CRT_WALL{WalletAddr:  walAddr,
		WalletType:     walletType,
		TimezoneOffset: timezoneOffset,
		FruitAmount:    amntStr,
		Identity:       identity,
		CreationDate:   walletCreationDate,
		PrivateKey : privKey}

	walletqueryJSON, _ := json.Marshal(walletInfo)

	return shim.Success(walletqueryJSON)
}


// 지갑 설정
// args[0] : R
// args[1] : S
// args[2] : X
// args[3] : Y
// args[4] : refID
// args[5] : 처리날짜,format must be 2018-10-08 00:00:00
// args[6] : 서비스코드

// args[7] - idientity
// args[8] - 지갑 주소
// args[9] - 지갑타입
// args[10] - 타임존
// args[11] - 열매 량
// args[12] - 마케팅 열매 선물 refID - 비필수
func (s *MCCWalletChaincode) restructWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) < 12 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "12 or 13"))
	}

	// 2. 파라메터 변수 할당
	//var msg string // 공통 메세지 처리 변수
	//var isOk bool	// 공통 성공 실패 체크 변수
	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	//svc := args[6] // 있으나 미사용으로 주석 처리

	identity := args[7] // identity
	walletAddr := args[8]	// 지갑주소
	walletType := args[9] // 지갑타입
	timezoneOffsetStr := args[10] // 타임존
	fruitAmntStr := args[11]	// 열매 량
	

	timezoneOffset, err := strconv.Atoi(timezoneOffsetStr)
	if err != nil {
		timezoneOffset = 0
	}

	// 3. 싸인 검증
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)
	if walletAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", walletAddr, pubKeyToCompare, X, Y))
	}
	
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID + procDate + identity+walletAddr + walletType+ timezoneOffsetStr+fruitAmntStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+" Y : "+Y+" R : "+R+" S : "+S+" refID : "+refID+" procDate : "+procDate+" identity : "+identity+" walletAddr : "+walletAddr+" walletType : "+walletType+" timezoneOffsetStr : "+timezoneOffsetStr+" fruitAmntStr : "+fruitAmntStr)))
	}

	// 4. 기본 유효성 체크 
	/*
	// 지갑 존재 여부
	// 추후 관리자 설정 키가 없으면, 재구성 불가 처리 -- 
	walletAsBytes, _ := stub.GetState(walletAddr)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletAddr))
	}
*/
	// 5. 기본 비지니스 로직 처리

	currentDate := time.Now().UTC() 	// 생성일
	walletCreationDate := string(currentDate.Format("2006-01-02 15:04:05"))

	var walletInfo = STRT_WALLET{
		WalletAddr:     walletAddr,
		WalletType:  walletType,
		TimezoneOffset: timezoneOffset,
		FruitAmount:   fruitAmntStr,
		Identity: identity,
		CreationDate: walletCreationDate}

	walletAsBytes, _ := json.Marshal(walletInfo)
	err = stub.PutState(walletInfo.WalletAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletAddr))
	}

	// 유니크 체크용 컴포짓 키 구성
	compositeKey, compositeErr := stub.CreateCompositeKey(cst.IndexUserID, []string{"UID", identity, walletAddr})
	if compositeErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositeKey+" : "+compositeErr.Error()))
	}

	// 컴포짓 키 등록
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", compositePutErr.Error()))
	}


	// 마케팅 이벤트 열매 ev0001기록 생성
	if len(args) == 13 {

		svc8QueueKey := args[12]	// 마케팅 열매 선물 queueKey

		evtCd := "ev0001"
		// ---------------이벤트 정보 취득 ----------------------
		eventStrt, errStr := s.getEventConfData(stub, []string{"getCoreEvt", evtCd})
		if errStr != "" {
			return shim.Error(fmt.Sprintf(mgmt.GetCCErrMsg("CC_COM_E09") + " : " + errStr))
		}

		eventDate := eventStrt.EvtDt
		evtKeyArgs := []string{eventDate[0:4], eventDate[5:7], eventDate[8:10], evtCd, walletAddr, svc8QueueKey} // refID가 아닌 받는 사용자 지갑을 키로 구성
		evtKeyArgs = append(evtKeyArgs, refID)
		coreChkCompErr := s.createCoreCompKey(stub, evtKeyArgs)
		if coreChkCompErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", procDate, refID))
		}
	}
	
	return shim.Success(walletAsBytes)
	

}


// 지갑 타임존 변경
// args[0] : R
// args[1] : S
// args[2] : X
// args[3] : Y
// args[4] : refID
// args[5] : 처리날짜,format must be 2018-10-08 00:00:00
// args[6] : 서비스코드

// args[7] : wallet address
// args[8] : timezoneOffset
func (s *MCCWalletChaincode) updateTimezoneOffset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 9 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "9"))
	}

	R := args[0]
	S := args[1]
	X := args[2]
	Y := args[3]
	refID := args[4]
	procDate := args[5]
	//svc := args[6] // 사용하지 않음

	walletAddr := args[7]
	timezoneOffsetStr := args[8]
	timezoneOffset, err := strconv.Atoi(timezoneOffsetStr)
	if err != nil {
		timezoneOffset = 0
	}


	// 3. 싸인 검증
	if !(tooltip.IsNumeric([]string{R, S})) {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E01", R, S))
	}

	pubKey := mw.GetPubKeyFromXandY(X, Y)
	pubKeyToCompare := mw.GetPubKeyHash(pubKey)
	if walletAddr != pubKeyToCompare {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E02", walletAddr, pubKeyToCompare, X, Y))
	}
	
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID + procDate + walletAddr + timezoneOffsetStr)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+" Y : "+Y+" R : "+R+" S : "+S+" refID : "+refID+" procDate : "+procDate+" walletAddr : "+walletAddr+" timezoneOffset : "+timezoneOffsetStr)))
	}


	walletAsBytes, _ := stub.GetState(walletAddr)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walletAddr))
	}

	mccWallet := STRT_WALLET{}
	json.Unmarshal(walletAsBytes, &mccWallet)

	mccWallet.TimezoneOffset = timezoneOffset
	walletAsBytes, _ = json.Marshal(mccWallet)

	err = stub.PutState(walletAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walletAddr))
	}

	res := STRT_WALLET{mccWallet.WalletAddr,
		mccWallet.WalletType,
		mccWallet.TimezoneOffset,
		mccWallet.FruitAmount,
		mccWallet.Identity,
		mccWallet.CreationDate}

	resJSON, _ := json.Marshal(res)

	return shim.Success(resJSON)
}

