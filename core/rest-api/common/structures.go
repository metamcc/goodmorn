package common


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
		RetryCnt			string `json:"retryCnt" form:"RetryCnt" valid:"required"`	// 재처리횟수
		ProcType 			string `json:"procType" form:"procType" valid:"required"`	// 큐처리타입
	}


	// 공통 파라메터 구조체
	COMM struct {
		
		CoreSign
		CoreCheck
		CommParam
	}
)