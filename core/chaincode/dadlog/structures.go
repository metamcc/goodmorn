package main

// MCCDADLogs struct
type STRT_DAD_MAIN struct {

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
}


// 코어 체크 응답 구조체 리턴용 리스트 구조체
type RES_CC_LST struct {
	CoreChkInfo []RES_CC `json:"coreChkList"`
}

// 코어 체크 응답 구조체
type RES_CC struct {
	ProcDate          	string `json:"procDate"`	// 로컬 처리날짜
	SenderAddr string `json:"senderAddr"`							// 보낸이 지갑주소
	RecipientAddr string `json:"recipientAddr"`						// 받는이 지갑주소
	RefID     string `json:"refID "`				// 참조키
	ProcDateUTC          	string `json:"procDateUTC"`	// utc 처리날짜
	CorePenaltySeed       		string `json:"corePenaltySeed"`	// 코어 페널티 적용 선물 양
	TxID          string `json:"txid"`				// txID
}


// 일반 응답 구조체
type RES_NOR struct {
	RES			string `json:"res"`	// 처리결과
	ProcDate          	string `json:"procDate"`	// 처리날짜
	TxID          string `json:"txid"`				// txID
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