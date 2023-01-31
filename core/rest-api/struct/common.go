package strt


type (

	// sing 검증 공통 구조체
	CoreSign struct {
				
		RHash				string `json:"rHash" form:"rHash" valid:"required"`	//R 해쉬
		SHash				string `json:"sHash" form:"sHash" valid:"required"`	//S 해쉬
		XValue				string `json:"xValue" form:"xValue" valid:"required"`	//x Value
		YValue				string `json:"yValue" form:"yValue" valid:"required"`	//y Value

	}


	// core 검증 공통 구조체
	CoreCheck struct {
				
		RefID				string `json:"refID" form:"refID" valid:"required"`	//refID
		ProcDate			string `json:"procDate" form:"procDate" valid:"cstDateTime,required"`	//처리날짜 - 포맷 - 2018-10-12 14:15:59
		
	}

	// 공통 파라메터 구조체
	CommParam struct {
			
		
		Svc					string `json:"svc" form:"svc" valid:"required"`	//서비스 코드
		ProcType			string `json:"procType" form:"procType"`	// 처리타입
	}



	// 공통 파라메터 구조체
	COMMSIGN struct {
		
		CoreSign
		COMM
	}


	// 공통 파라메터 구조체
	COMM struct {
		
		CoreCheck
		CommParam
	}

	// 코어 체크 응답 구조체 - 리스트
	RES_CC_LST struct {
		CoreChkInfo []RES_CC `json:"coreChkList"`
	}

	// 코어 체크 응답 구조체 - 객체
	RES_CC struct {
		ProcDate          	string `json:"procDate"`	// 처리날짜
		SvcCd     string `json:"svcCd"`				// 서비스코드
		RefID     string `json:"refID"`				// 참조키
		TxID          string `json:"txid"`				// txID
		RetryID     string `json:"retryID"`				// 재처리 참조키
	}
)