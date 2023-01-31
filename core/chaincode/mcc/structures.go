package main

// F2F 응답 구조체
type RES_F2F struct {
	SenderAddr string `json:"senderAddr"`							// 보낸이 지갑주소
	FruitAmntAfter        	string `json:"fruitAmntAfter"`	// 전체열매잔량 - 이체 후
	ReqAmnt		string `json:"reqAmnt"`		// 요청량 
	RecipientAddr string `json:"recipientAddr"`							// 받는 지갑주소
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`				// txID
	RefID     string `json:"refID"`				// 참조키
}

// F2M 응답 구조체
type RES_F2M struct {
	WalletAddr 		string `json:"walletAddr"`			// 사용자 지갑 주소
	EthereumAddr     	string `json:"ethereumAddr"` 	// 사용자 이더리움 입금 주소
	FruitAmntAfter        	string `json:"fruitAmntAfter"`	// 전체열매잔량 - 전환 후
	Amount		string `json:"amount"`		// 전체전환량 
	BaseFee			string `json:"baseFee"`		// 기본수수료량
	ExtraFee			string `json:"extraFee"`		// 추가수수료량
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          	string `json:"txid"`			// txID
	RefID     string `json:"refID"`				// 참조키
}


// M2F 응답 구조체
type RES_M2F struct {
	WalletAddr 		string `json:"walletAddr"`			// 지갑 주소
	EthereumAddr     	string `json:"ethereumAddr"` 	// 사용자입금이더리움 주소
	EthereumHash	string `json:"ethereumHash"` 	// 처리이더리움 해쉬
	FruitAmntAfter        	string `json:"fruitAmntAfter"`	// 전체열매잔량 - 전환 후
	
	ExchAmnt		string `json:"exchAmnt"`		// 전환량 
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          	string `json:"txid"`			// txID
	RefID     string `json:"refID"`				// 참조키
}



// 사용자 개별 정산 응답 구조체
type RES_DAD_WAL struct {
	WalletAddr 			string `json:"walletAddr"`		// 사용자 지갑 주소
	FruitAmntAfter        	string `json:"fruitAmntAfter"`	// 전체열매잔량 - 이체 후
	UserDailySeed        	string `json:"userDailyDad"`	// 일일 DAD 받은 씨앗 갯수
	TotalDailySeed       		string `json:"totalDailyDad"`	// 당일 전체 DAD 씨앗 갯수


	FruitDAD	string `json:"fruitDAD"` // 개인별 일일 배분 열매 량
	ProcDate          string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`
	RefID     string `json:"refID"`				// 참조키
}


// 에코 지갑 정산 응답 구조체
type RES_DAD_ECO struct {
	WalletAddr 		string `json:"walletAddr"` // 지갑 주소
	FruitAmntAfter   		string `json:"amount"` // 정산 후 열매량

	FruitTotalDAD	string `json:"fruitTotalDAD"` // 일일 배분 열매 량
	ProcDate          string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`
	RefID     string `json:"refID"`				// 참조키
}


// 씨앗 선물 응답 구조체
type RES_SSG struct {
	SenderAddr string `json:"senderAddr"`							// 보낸이 지갑주소
	RecipientAddr string `json:"recipientAddr"`						// 받는이 지갑주소
	ParamBaseSeed        	string `json:"paramBaseSeed"`	// 개인별 기준 선물 갯수
	ParamPenaltySeed       		string `json:"paramPenaltySeed"`	// 파라메터 페널티 적용 선물 양
	CorePenaltySeed       		string `json:"corePenaltySeed"`	// 코어 페널티 적용 선물 양
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`				// txID
	RefID     string `json:"refID"`				// 참조키

	ErrData string `json:"errData"`				// 에러내용
}


type RES_SSG_LST struct {
	SeedingResList []RES_SSG `json:"seedingResList"`
	
}


// 코어 체크 응답 구조체 리턴용 리스트 구조체
type RES_CC_LST struct {
	CoreChkInfo []RES_CC `json:"coreChkList"`
}

// 코어 체크 응답 구조체
type RES_CC struct {
	ProcDate          	string `json:"procDate"`	// 처리날짜
	SvcCd     string `json:"svcCd"`				// 서비스코드
	RefID     string `json:"refID"`				// 참조키
	TxID          string `json:"txid"`				// txID
}

// 시스템 열매 선물 응답 구조체
type RES_EV_F2F struct {
	WalletAddr 		string `json:"walletAddr"` // 받는 지갑 주소
	FruitAmntAfter        	string `json:"fruitAmntAfter"`	// 전체열매잔량 - 이체 후
	MarketingAddr string `json:"marketingAddr"`							// 보낸 마케팅 지갑주소
	GiftAmount string `json:"giftAmount"`								// 선물량
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`				// txID
	RefID     string `json:"refID"`				// 참조키
}


// MCC 지갑 구조체 - 조회 응답 구조체 포함
type STRT_WALLET struct {
	WalletAddr  	string `json:"walletAddr"`
	WalletType     string `json:"walletType"`
	TimezoneOffset int    `json:"timezoneOffset"` //distance of seconds base UTC
	FruitAmount    string `json:"fruitAmount"`    //
	Identity       string `json:"identity"`
	CreationDate   string `json:"creationDate"`
}


// MCC 지갑 생성 응답 구조체
type RES_CRT_WALL struct {
	WalletAddr  	string `json:"wallet"`
	WalletType     string `json:"walletType"`
	TimezoneOffset int    `json:"timezoneOffset"` //distance of seconds base UTC
	FruitAmount    string `json:"amount"`    //
	Identity       string `json:"identity"`
	CreationDate   string `json:"creationDate"`
	PrivateKey  	string `json:"privatekey"`
	RefID     string `json:"refID"`				// 참조키
}



// MCC 지갑 조회 상세 응답 구조체
type RES_DTL_WALL struct {
	WalletAddr  	string `json:"wallet"`
	WalletType     string `json:"walletType"`
	TimezoneOffset int    `json:"timezoneOffset"` //distance of seconds base UTC
	TotalFruit    string `json:"totalFruit"`    //
	Identity       string `json:"identity"`
	CreationDate   string `json:"creationDate"`

	WalletFruit    string `json:"walletFruit"`    //
	PaymentFruit    string `json:"paymentFruit"`    //
	PaymentCnt	int `json:"paymentCnt"`    //
	PaymentData []FRUIT_PAYMENT_STRT `json:"paymentData"`	//
	RefID     string `json:"refID"`				// 참조키
}

// MCC 지갑 조회 간편 응답 구조체
type RES_SMP_WALL struct {
	WalletAddr  	string `json:"wallet"`
	WalletType     string `json:"walletType"`
	TimezoneOffset int    `json:"timezoneOffset"` //distance of seconds base UTC
	TotalFruit    string `json:"totalFruit"`    //
	Identity       string `json:"identity"`
	CreationDate   string `json:"creationDate"`

	WalletFruit    string `json:"walletFruit"`    //
	PaymentFruit    string `json:"paymentFruit"`    //
	PaymentCnt	int `json:"iaymentCnt"`    //
	FruitSyncData []string `json:"fruitSyncData"`
	RefID     string `json:"refID"`				// 참조키
}



// 지갑조회 - 열매 결제 내역 구조체
type FRUIT_PAYMENT_STRT struct {
	MyWalletAddr  	string `json:"myWalletAddr"`
	ProcDate          	string `json:"procDate"`	// 처리날짜

	RefWalletAddr  	string `json:"refWalletAddr"`
	PaymentFruit    string `json:"paymentFruit"`    //
	TxID          string `json:"txid"`				// txID
	RefID     string `json:"refID"`				// 참조키
}

// 지갑 조회 - 열매 결제 내역 리스트
type RES_FRUIT_PAYMENT_LST struct {
	FruitPaymentList []FRUIT_PAYMENT_STRT `json:"fruitPaymentList"`
}

// 이벤트 구성 정보 구조체
type STRT_EVT struct {
	EvtCd    string  `json:"evtCd"`  // 이벤트코드
	EvtDt  string  `json:"evtDt"`	// 이벤트기준일자
	EvtAmnt  string  `json:"evtAmnt"`	// 이벤트열매량
	EvtCnt int  `json:"evtCnt"` 	// 이벤트참가횟수 - 0 - 매번 허용 , 1 : 유일 1회 허용,  2 - 2회 허용....
	Etc    string  `json:"etc"`	// 기타
}


// 이벤트 구성 정보 구조체
type RES_ALIVE struct {
	RestServer          string `json:"restServer"`				// PeerAddr
	ProcDate          	string `json:"procDate"`	// 처리날짜
	RefID          		string `json:"refID"`				// RefID
	RetryID          	string `json:"retryID"`				// RetryfID
}

// MCCFunctionStop structure
type MCCFunctionStop struct {
	IsStopped bool `json:"isstopped"`
}

// MCCResStatus structure
type MCCResStatus struct {
	ResMsg string `json:"resMsg"`
}

// MCCTokeExchPeriod structure
type MCCTokeExchPeriod struct {
	Mins string `json:"minuts"`
}

// MCCConfig structure
type MCCConfig struct {
	Configuration string `json:"configuration"`
	Messages      string `json:"message"`
}

// AllAliveConfInfo structure : getAliveStatHist
type AllAliveConfInfo struct {
	ConfInfo []MCCConfig
}


type STRT_ADM2 struct {

	MsgDataMap  map[string]string  `json:"msgDataMap"`
	VarDataMap  map[string]string  `json:"varDataMap"`
	FncDataMap  map[string]string  `json:"fncDataMap"`

	MsgCnt int `json:"msgCnt"`	// 총갯수
	VarCnt int `json:"varCnt"`	// 총갯수
	FncCnt int `json:"fncCnt"`	// 총갯수
}


type STRT_ADM struct {

	MgmtDataMap  map[string]string  `json:"coreDataMap"`	// 데이터 맵
	Cnt int `json:"cnt"`	// 총갯수
}


// 열매 결제 구초제
type STRT_FRUIT_PAYMENT struct {

	SenderAddr string `json:"senderAddr"`							// 보낸이 지갑주소
	ReceiverAddr string `json:"receiverAddr"`							// 받는 지갑주소
	ReqAmnt		string `json:"reqAmnt"`		// 요청량 

	SenderSyncData []string `json:"senderSyncData"` // 동기화 대상 키 - 보낸이
	ReceiverSyncData string `json:"receiverSyncData"` // 동기화 대상 키 - 받는 이


	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`				// txID
	RefID     string `json:"refID"`				// 참조키
	OrgRefID     string `json:"orgRefID"`				// 취소용 원본 참조키

}


// 지갑 동기화 대상 결과 구조체
type RES_FRUIT_BATCH_PAYMENT_STRT struct {
	RefID          string `json:"refID"`				// RefID
	ErrStr		string `json:"errStr"`				// txID
}


// 지갑 동기화 구조체
type STRT_SYNC_WALLET struct {

	WalletAddr string `json:"walletAddr"`							// 동기화 지갑주소
	RefID     string `json:"refID"`				// 참조키
	SyncAmnt		string `json:"syncAmnt"`		// 동기화량 

	SyncSuccessData []RES_FRUIT_BATCH_PAYMENT_STRT `json:"syncSuccessData"`
	NextFruitSyncData []string `json:"nextFruitSyncData"`			// 다음 싱크 참조키
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`				// txID
}


/*
type STRT_FRUIT_PAYMENT3 struct {


	SenderAddr string `json:"senderAddr"`							// 보낸이 지갑주소
	ReceiverAddr string `json:"receiverAddr"`							// 받는 지갑주소
	ReqAmnt		string `json:"reqAmnt"`		// 요청량 
	ID     string `json:"iD"`				// 참조키

	FruitPaymentSuccessData []RES_FRUIT_BATCH_PAYMENT_STRT `json:"fruitPaymentSuccessData"`
	FruitPaymentFailData  []RES_FRUIT_BATCH_PAYMENT_STRT `json:"fruitPaymentFailData"`
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`				// txID
}

*/