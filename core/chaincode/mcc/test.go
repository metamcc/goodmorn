package main


import (
	"fmt"
	"encoding/json"
	"math/big"

	cst "github.com/mycreditchain/chaincode/mcc/constants"
	mgmt "github.com/mycreditchain/chaincode/mcc/utils/messages"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)


func (s *MCCWalletChaincode) chkCoreCompKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexCoreCheck, args)
	
	if msgErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
	}
	defer keyResultsIterator.Close()

	// Check the variable existed
	if !keyResultsIterator.HasNext() {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02"))
	}


	for keyResultsIterator.HasNext() {
		keyRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		stub.DelState(keyRange.Key)
	}


	var resMsg = MCCResStatus{
		ResMsg: "Ok"}

	resMsgAsBytes, _ := json.Marshal(resMsg)

	return shim.Success(resMsgAsBytes)


}




/*
	테스트용 
*/
func (s *MCCWalletChaincode) delCoreCompKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexCoreCheck, args)
	
	if msgErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
	}
	defer keyResultsIterator.Close()

	// Check the variable existed
	if !keyResultsIterator.HasNext() {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02"))
	}


	for keyResultsIterator.HasNext() {
		keyRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		stub.DelState(keyRange.Key)
	}


	var resMsg = MCCResStatus{
		ResMsg: "Ok"}

	resMsgAsBytes, _ := json.Marshal(resMsg)

	return shim.Success(resMsgAsBytes)


}	


	// 테스트용
func (s *MCCWalletChaincode) makeData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	walAddr := args[0]
	fruit := args[1]
	typ := args[2]

	walletAsBytes, _ := stub.GetState(walAddr)
	if walletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walAddr))
	}

	if typ == "user" {
		wallet := STRT_WALLET{}
		json.Unmarshal(walletAsBytes, &wallet)

		newFruit := new(big.Int)
		newFruit, _ = newFruit.SetString(fruit, 10)
		wallet.FruitAmount = fmt.Sprint(newFruit)

		walletAsBytes, _ = json.Marshal(wallet)
	} else if typ == "system" {
		wallet := STRT_WALLET{}
		json.Unmarshal(walletAsBytes, &wallet)

		newFruit := new(big.Int)
		newFruit, _ = newFruit.SetString(fruit, 10)

		wallet.FruitAmount = fmt.Sprint(newFruit)

		walletAsBytes, _ = json.Marshal(wallet)
	}
	err := stub.PutState(walAddr, walletAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walAddr))
	}

	return shim.Success(walletAsBytes)
}
	






// 테스트용
func (s *MCCWalletChaincode) makeData2(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	walId := args[0]
	walAddr := s.findWalletByID(stub, "UID", walId)
	fruit := args[1]


	if walAddr == "" {
		keys := []string{walId, "user", "32400"}
		return s.createMCCWallet(stub, keys)
	} else {
		walletAsBytes, _ := stub.GetState(walAddr)
		if walletAsBytes == nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", walAddr))
		}		
			wallet := STRT_WALLET{}
			json.Unmarshal(walletAsBytes, &wallet)
	
			newFruit := new(big.Int)
			newFruit, _ = newFruit.SetString(fruit, 10)
			wallet.FruitAmount = fmt.Sprint(newFruit)
	
			walletAsBytes, _ = json.Marshal(wallet)
	
		err := stub.PutState(walAddr, walletAsBytes)
		if err != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", walAddr))
		}
	
		return shim.Success(walletAsBytes)
	}
}



// 테스트용
func (s *MCCWalletChaincode) makeData3(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	cnt,_ := strconv.Atoi(args[0])
	preFixId := args[1]
	for i := 0;i<cnt; i++ {
	//	iStr := strconv.Itoa(i)
		walId := preFixId+strconv.Itoa(i)
		keys := []string{walId, "user", "32400"}
		s.createMCCWallet(stub, keys)
	}
	

	return shim.Success(nil)
}







/*

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
func (s *MCCWalletChaincode) payFruit2(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) < 6 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "6"))
	}

	// 2. 파라메터 변수 할당
	var isOk bool	// 공통 성공 실패 체크 변수
	var orgRefID string // 결제취소 - 원본 큐키 기록용

	refID := args[0]
	procDate := args[1]
	svcPrefix := args[2] // 서비스코드

	sender := args[3]
	receiver := args[4]
	reqAmntStr := args[5]

	
	txId := stub.GetTxID()
	reqAmntInt := new(big.Int)
	reqAmntInt, isOk = reqAmntInt.SetString(reqAmntStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "reqAmnt", reqAmntStr))
	}

	// 3. 싸인 검증 


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
	senderWalletSMPInfo, _, errStr := s.getWalletSMPInfo2(stub, []string{"address", sender})
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


	// ---------------응답 구성----------------------
	resInfo := STRT_FRUIT_PAYMENT{sender, receiver, reqAmntStr, senderWalletSMPInfo.FruitSyncData, procDate, txId,  refID, orgRefID}
	resInfoJSON, _ := json.Marshal(resInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	err = s.createCoreCompKey(stub, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	return shim.Success(resInfoJSON)
}

	열매결제 1안 - 컴포짓킨ㄴ +값만 쌓이고, 결제나 이체시 컴포짓키를 제거 즉시 제거(요청양만큼만)
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
	var orgRefID string
	if len(args) > 10 {
		orgRefID = args[10] // 결제취소 - 원본 큐키 기록용
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
	//senderWalletSMPInfo, orgWalletStrt, errStr := s.getWalletSMPInfo(stub, []string{"address", sender})
	//if errStr != "" {
	//	return shim.Error(errStr)
	//}

	// 잔액 검증
//	senderFruitBefore := new(big.Int)
//	senderFruitBefore, isOk = senderFruitBefore.SetString(senderWalletSMPInfo.TotalFruit, 10)
//	if !isOk {
//		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderWalletSMPInfo.TotalFruit))
//	}

//	if senderFruitBefore.Cmp(reqAmntInt) == -1 {
//		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderWalletSMPInfo.TotalFruit, reqAmntStr))
//	}




	


	// 잔액 검증


	// 기존 지갑 열매 취득


	// 출금 처리 
	// 1. 결제 데이터 우선 차감
	// 2. 결제 데이터 차감 후 모자른 금액은 지갑 데이터 차감


	// 지갑 조회
	senderWalletAsBytes, _ := stub.GetState(sender)
	if senderWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01",  sender))
	}
	senderWalletStruct := STRT_WALLET{}
	json.Unmarshal(senderWalletAsBytes, &senderWalletStruct)

	// 지갑 열매량
	senderWalletFruitBefore := new(big.Int)
	senderWalletFruitBefore, isOk = senderWalletFruitBefore.SetString(senderWalletStruct.FruitAmount, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderWalletStruct.FruitAmount))
	}

	
	// 결제 데이터 점검
	dataActFlow := 0 // 결제 데이터 전결 유무
	paymentFruitTotalAmnt := new(big.Int)
	paymentFruitAmnt := new(big.Int)

	keyArgs = []string{sender} 	
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexFruitPayment, keyArgs)
	if msgErr != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error()))
	}
	defer keyResultsIterator.Close()


	// 기존 결제 데이터 합산 및 제거 
	for keyResultsIterator.HasNext() {
		keyRange, nextErr := keyResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}
		logger.Debug(keyRange)


		err = stub.DelState(keyRange.Key)
		if err != nil {
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E13", sender))
		}

		if nextErr != nil {
			//return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		_, keyParts, splitKeyErr := stub.SplitCompositeKey(keyRange.Key)
		if splitKeyErr != nil {
			//return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", splitKeyErr.Error()))
			return shim.Error(mgmt.GetCCErrMsg("CC_COP_E02", nextErr.Error()))
		}

		logger.Info(keyParts)
		 
		paymentFruitAmnt.SetString(keyParts[6], 10)
		paymentFruitTotalAmnt.Add(paymentFruitTotalAmnt, paymentFruitAmnt)
		
		//calcFruitTotalAmnt.Add(paymentFruitTotalAmnt, senderWalletFruitBefore)
		if paymentFruitTotalAmnt.Cmp(reqAmntInt) != -1 {
			//return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderWalletSMPInfo.TotalFruit, reqAmntStr))
			dataActFlow = 1 // 결제 데이터로 모두 차감
			break
		}



	}

	// 1. 결제 데이터에서 이체량 빼고 남은 열매는 지갑으로 동기화
	// 2. 결제 데이터에서 모자른 양은 지갑 열매애서 차감

	totalBalance := new(big.Int)
	senderFruitAfter := new(big.Int)
//	


	// 결제 데이터 전결시 잔액을 지갑 열매에 추가
	if dataActFlow == 1 {

		totalBalance.Sub(paymentFruitTotalAmnt, reqAmntInt)
		senderFruitAfter.Add(senderWalletFruitBefore, totalBalance)

	} else { // 부족시 지갑 열매에서 추가 차감

		// 잔액 검증
		totalBalance.Add(paymentFruitTotalAmnt, senderWalletFruitBefore)
		if totalBalance.Cmp(reqAmntInt) != -1 { // 잔액 부족
			return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", fmt.Sprint(totalBalance), reqAmntStr))
		}

		// 전체 열매에서 요청량 차감
		senderFruitAfter.Sub(totalBalance, reqAmntInt)

	}

	// 지갑 데이터 구성
	senderWalletStruct.FruitAmount = fmt.Sprint(senderFruitAfter)

	// 지갑 갱신
	senderAsBytes, _ := json.Marshal(senderWalletStruct)
	err = stub.PutState(sender, senderAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", sender))
	}


	// ---------------받는이----------------------

	// 기록 구성
	keyArgs = []string{receiver, procDate[0:4], procDate[5:7], procDate[8:10], refID, sender, reqAmntStr, txId }
	err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}
	
	// ---------------응답 구성----------------------
	
	resInfo := STRT_FRUIT_PAYMENT2{sender, receiver, reqAmntStr,  []string{},  procDate, txId, refID, orgRefID}
	resInfoJSON, _ := json.Marshal(resInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	err = s.createCoreCompKey(stub, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	return shim.Success(resInfoJSON)
}



func (s *MCCWalletChaincode) getWalletSMPInfo(stub shim.ChaincodeStubInterface, args []string) (RES_DTL_WALL, STRT_WALLET, string) {

	//var resMsg []FRUIT_PAYMENT_STRT 
	var walletToSearch string
	queryType := args[0]
	queryData := args[1]
	myWalletStruct := RES_DTL_WALL{} 
	walletStruct := STRT_WALLET{}
	

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

	keyArgs := []string{walletStruct.WalletAddr} 	
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexFruitPayment, keyArgs)
	if msgErr != nil {

		return myWalletStruct, walletStruct, mgmt.GetCCErrMsg("CC_COP_E02", msgErr.Error())
	}
	defer keyResultsIterator.Close()

	paymentFruitTotal := new(big.Int)
	paymentFruitAmt := new(big.Int)
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


		logger.Debug(keyParts)
	
		 
		 paymentFruitAmt.SetString(keyParts[6], 10)

		 paymentFruitTotal = paymentFruitTotal.Add(paymentFruitTotal, paymentFruitAmt)

		 
		// resMsg = append(resMsg, FRUIT_PAYMENT_STRT{keyParts[0], keyParts[1]+keyParts[2]+keyParts[3], keyParts[4], keyParts[5], keyParts[6], keyParts[7]}) 
	}



	
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

		walletTotalFruit := new(big.Int)
		walletTotalFruit = walletTotalFruit.Add(paymentFruitTotal, walletFruitAmt)

		myWalletStruct.WalletFruit = walletStruct.FruitAmount
		myWalletStruct.PaymentFruit = fmt.Sprint(paymentFruitTotal)
		myWalletStruct.TotalFruit = fmt.Sprint(walletTotalFruit)

		
		//myWalletStruct.PaymentData = resMsg
		
		
	} else {
		myWalletStruct.WalletFruit = walletStruct.FruitAmount
		myWalletStruct.PaymentFruit = "0"
		myWalletStruct.TotalFruit = walletStruct.FruitAmount
	}


	return myWalletStruct, walletStruct, ""
}








// 시스템 열매 정산 1:N
func (s *MCCWalletChaincode) payFruit3(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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
	totalReqAmntStr := args[8]
	receivers := args[9]

	txId := stub.GetTxID()

	totalReqAmntInt := new(big.Int)
	totalReqAmntInt, isOk = totalReqAmntInt.SetString(totalReqAmntStr, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "reqAmnt", totalReqAmntStr))
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
	isVerified, _ := s.verifyArguments(stub, []string{pubKey, R, S, fmt.Sprintf(refID+procDate + sender + totalReqAmntStr + receivers)})
	if isVerified == false {
		return shim.Error(mgmt.GetCCErrMsg("CC_SGN_E03", fmt.Sprint("pubKey : "+pubKey+" X : "+X+
			"Y : "+Y+" R : "+R+" S : "+S+" procDate : "+procDate+" sender : "+sender+" receiver : "+receivers+" amount : "+totalReqAmntStr)))
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



	// 5. 비지니스 로직
	// ---------------보내는이----------------------
	// 보내는 이 지갑 조회 - 간편조회
	senderWalletSMPInfo, orgWalletStrt, errStr := s.getWalletSMPInfo2(stub, []string{"address", sender})
	if errStr != "" {
		return shim.Error(errStr)
	}
	
	// 잔액 검증
	senderFruitBefore := new(big.Int)
	senderFruitBefore, isOk = senderFruitBefore.SetString(senderWalletSMPInfo.TotalFruit, 10)
	if !isOk {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E10", "sender Fruit ", senderWalletSMPInfo.TotalFruit))
	}

	if senderFruitBefore.Cmp(totalReqAmntInt) == -1 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E08", senderWalletSMPInfo.TotalFruit, totalReqAmntStr))
	}

	// 기록 구성 - 지갑직접갱신, 결제데이터 기록
	
	// keyArgs = []string{sender, procDate[0:4], procDate[5:7], procDate[8:10], refID, receiver, "-"+totalReqAmntStr, txId }
	// err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
	// if err != nil {
	// 	return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	// }
	

	// ---------------받는이----------------------
	/*
		// 동일 지갑
		if sender == receiver {
			return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E02", sender, receiver))
		}

	// 보내는 이 지갑 조회 - 간편조회
	receiverWalletSMPInfo, _, errStr := s.getWalletSMPInfo2(stub, []string{"address", receiver})
	if errStr != "" {
		return shim.Error(errStr)
	}

	// 받는이 기록 구성
	var receiverSuccessMsg []RES_FRUIT_BATCH_PAYMENT_STRT 
	var receiverFailMsg []RES_FRUIT_BATCH_PAYMENT_STRT 
	var receiverInfo []string

	paymentFruitTotal := new(big.Int)
	paymentFruitAmt := new(big.Int)

	receiversArr := strings.Split(receivers, ",")
	for i, value := range receiversArr { 
		fmt.Println(i, value)
		
		receiverInfo = strings.Split(value, "|") // refID, walletAddr, fruitAmnt

		// 동일지갑
		// 지갑존재여부
		// 금액 존재여부 및 0원

		keyArgs = []string{receiverInfo[1], procDate[0:4], procDate[5:7], procDate[8:10], receiverInfo[0], sender, receiverInfo[2], txId, procDate[11:19] }
		err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
		if err != nil { // 구성 실패
			receiverFailMsg = append(receiverFailMsg, RES_FRUIT_BATCH_PAYMENT_STRT{receiverInfo[0], mgmt.GetCCErrMsg("CC_COP_E01", err.Error())})
		} else { // 구성 성공

			
			paymentFruitTotal = paymentFruitTotal.Add(paymentFruitTotal, paymentFruitAmt)

			receiverSuccessMsg = append(receiverSuccessMsg, RES_FRUIT_BATCH_PAYMENT_STRT{receiverInfo[0], ""})
		}

		// resMsg = append(resMsg, FRUIT_PAYMENT_STRT{keyParts[0], keyParts[1]+keyParts[2]+keyParts[3], keyParts[4], keyParts[5], keyParts[6], keyParts[7]}) 
	}

	// keyArgs = []string{receiver, procDate[0:4], procDate[5:7], procDate[8:10], refID, sender, reqAmntStr, txId }
	// logger.Info(keyArgs)
	// fmt.Println(keyArgs)
	// err = s.createCompKey(stub, cst.IndexFruitPayment, keyArgs)
	// if err != nil {
	// 	return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	// }





	// ---------------보내는이----------------------
	// 보내는이 데이터 최종 갱신 
	senderFruitAfter := new(big.Int)
	senderFruitAfter = senderFruitAfter.Sub(senderFruitBefore, paymentFruitTotal)
	orgWalletStrt.FruitAmount = fmt.Sprint(senderFruitAfter)

	// 지갑 갱신
	senderAsBytes, _ := json.Marshal(orgWalletStrt)
	err = stub.PutState(sender, senderAsBytes)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E03", sender))
	}

	//receivers = ""
	// ---------------응답 구성----------------------
	//resInfo := STRT_FRUIT_PAYMENT2{sender, receiver, totalReqAmntStr,  refID, senderWalletSMPInfo.FruitSyncData, procDate, txId}
	resInfo := STRT_FRUIT_PAYMENT3{sender, receivers, totalReqAmntStr,  refID, receiverSuccessMsg,  receiverFailMsg, procDate, txId}
	resInfoJSON, _ := json.Marshal(resInfo)

	// ---------------이벤트 설정----------------------
	eventID := fmt.Sprint(svcPrefix +"_"+refID)
	err = stub.SetEvent(eventID, resInfoJSON)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_EVT_E01", eventID, err.Error()))
	}

	// ---------------코어 검증용 컴포짓 키 구성----------------------
	keyArgs = append(keyArgs, txId)
	err = s.createCoreCompKey(stub, keyArgs)
	if err != nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_COP_E01", err.Error()))
	}

	return shim.Success(resInfoJSON)
}








// 열매정보키
func (s *MCCWalletChaincode) getFruitCompKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 1. 파라메터 갯수 체크
	if len(args) != 4 {
		return shim.Error(mgmt.GetCCErrMsg("CC_COM_E01", "4"))
	}

	// 2. 파라메터 변수 할당
	myWalletAddr := args[0]	// 보낸이 지갑주소
	procDate := args[1]
	refID := args[2]
	receiverWalletAddr := args[3]	// 받는이 지갑 주소
	
	var resMsg []FRUIT_PAYMENT_STRT 

	// 내 지갑 - 존재 여부
	myWalletAsBytes, _ := stub.GetState(myWalletAddr)
	if myWalletAsBytes == nil {
		return shim.Error(mgmt.GetCCErrMsg("CC_WLT_E01", myWalletAddr))
	}

	// 내 지갑 
	myMCCWalletStrt := STRT_WALLET{}
	json.Unmarshal(myWalletAsBytes, &myMCCWalletStrt)
	
	// 5. 비지니스 로직
	keyArgs := []string{myWalletAddr, procDate[0:4], procDate[5:7], procDate[8:10], refID, receiverWalletAddr}
	keyResultsIterator, msgErr := stub.GetStateByPartialCompositeKey(cst.IndexFruitPayment, keyArgs)
	
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

		resMsg = append(resMsg, FRUIT_PAYMENT_STRT{keyParts[0], keyParts[1]+keyParts[2]+keyParts[3], keyParts[4], keyParts[5], keyParts[6], keyParts[7]})
	}

	fruitPaymentInfo := RES_FRUIT_PAYMENT_LST{FruitPaymentList: resMsg}
	fruitPaymentInfoAsBytes, _ := json.Marshal(fruitPaymentInfo)

	return shim.Success(fruitPaymentInfoAsBytes)


}






*/