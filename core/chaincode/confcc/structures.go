package main

// MCCConfigChaincode struct
type MCCConfigChaincode struct {
}

type STRT_EVT struct {
	EvtCd    string  `json:"evtCd"`  // 이벤트코드
	EvtDt  string  `json:"evtDt"`	// 이벤트기준일자
	EvtAmnt  string  `json:"evtAmnt"`	// 이벤트열매량
	EvtCnt int  `json:"evtCnt"` 	// 이벤트참가횟수 - 0 - 매번 허용 , 1 : 유일 1회 허용,  2 - 2회 허용....
	Etc    string  `json:"etc"`	// 기타
}


// 어드민 관리 항목 구조체
type STRT_ADM struct {

	MgmtDataMap  map[string]string  `json:"coreDataMap"`	// 데이터 맵
	Cnt int `json:"cnt"`	// 총갯수
}

// 어드민 비밀번호 구조체
type STRT_PWD struct {

	Pwd string  `json:"pwd"` // 비밀번호
}