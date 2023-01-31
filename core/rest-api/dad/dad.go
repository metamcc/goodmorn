package dad

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	"github.com/mycreditchain/echo-rest-server/util"
	//"strconv"
	CM_STRT "github.com/mycreditchain/rest-api/struct"



)

// 파라메터 구조체 - 적용되는 파라메터 타입에 대한 설정 포함
type (

/*

	MultifulSeeds struct {

		ValidationChecks		string `json:"validationChecks" form:"validationChecks" valid:"required"`	//해쉬값 체크 - , R, S, x, y
		
		RHashs				string `json:"rHashs" form:"rHashs" valid:"required"`	//R 해쉬
		SHashs				string `json:"sHashs" form:"sHashs" valid:"required"`	//S 해쉬
		XValues				string `json:"xValues" form:"xValues" valid:"required"`	//x Value
		YValues				string `json:"yValues" form:"yValues" valid:"required"`	//y Value

		QueueKeys			string `json:"queueKeys" form:"queueKeys" valid:"required"`	//refID
		SenderAddrs			string `json:"senderAddrs" form:"senderAddrs" valid:"required"`	//보내는이 지갑주소
		RecipientAddrs		string `json:"recipientAddrs" form:"recipientAddrs" valid:"required"`	//받는이 지갑주소
		RefID				string `json:"refID" form:"refID" valid:"required"`	//refID
		Amounts				string `json:"amounts" form:"amounts" valid:"required"`	//보내는양
		Date				string `json:"date" form:"date" valid:"cstDateTime,required"`	//날짜 - 포맷 -2018-10-12 14:15:59
		Separator			string `json:"separator" form:"separator" valid:"required"`	//refID
		
		

	}
*/
	STRT_SSG struct {

		CM_STRT.COMMSIGN

		SenderAddr			string `json:"senderAddr" form:"senderAddr" valid:"required"`	//보내는이 지갑주소
		RecipientAddr		string `json:"recipientAddr" form:"recipientAddr" valid:"required"`	//받는이 지갑주소
		BaseAmount				string `json:"baseAmount" form:"baseAmount" valid:"required"`	// 씨앗 선물 기준 양
		PenaltyAmount				string `json:"penaltyAmount" form:"penaltyAmount" valid:"required"`	// 페널티 포함 씨앗 선물 양
		
	}

	STRT_DAD_WLT struct {

		CM_STRT.COMMSIGN

		WalletAddr				string `json:"walletAddr" form:"walletAddr" valid:"required"`	// 에코 지갑 아이디
		PenaltySeed				string `json:"penaltySeed" form:"penaltySeed" valid:"required"`	//에코 지갑 주소
		PenaltySeedSum				string `json:"penaltySeedSum" form:"penaltySeedSum" valid:"required"`		// 정산량

	}

	STRT_DAD_ECO struct {

		CM_STRT.COMMSIGN

		EcoID				string `json:"ecoID" form:"ecoID" valid:"required"`	// 에코 지갑 아이디
		EcoAddr				string `json:"ecoAddr" form:"ecoAddr" valid:"required"`	//에코 지갑 주소
		Amount				string `json:"amount" form:"amount" valid:"required"`		// 정산량
	}






	STRT_SSGM struct {

		CM_STRT.CoreSign

						
		RefID				string `json:"refID" form:"refID" valid:"required"`	//refID
		ProcDate			string `json:"procDate" form:"procDate" valid:"required"`	//처리날짜
		Svc					string `json:"svc" form:"svc" valid:"required"`	//서비스 코드
		ProcType			string `json:"procType" form:"procType"`	// 처리타입
		Delimiter			string `json:"delimiter" form:"delimiter" valid:"required"`	// 구분 문자

		SenderAddr			string `json:"senderAddr" form:"senderAddr" valid:"required"`	//보내는이 지갑주소
		RecipientAddr		string `json:"recipientAddr" form:"recipientAddr" valid:"required"`	//받는이 지갑주소
		BaseAmount				string `json:"baseAmount" form:"baseAmount" valid:"required"`	// 씨앗 선물 기준 양
		PenaltyAmount				string `json:"penaltyAmount" form:"penaltyAmount" valid:"required"`	// 페널티 포함 씨앗 선물 양
		
	}
	

	

	STRT_DAD_WLTM struct {

		CM_STRT.CoreSign

						
		RefID				string `json:"refID" form:"refID" valid:"required"`	//refID
		ProcDate			string `json:"procDate" form:"procDate" valid:"required"`	//처리날짜
		Svc					string `json:"svc" form:"svc" valid:"required"`	//서비스 코드
		ProcType			string `json:"procType" form:"procType"`	// 처리타입
		Delimiter			string `json:"delimiter" form:"delimiter" valid:"required"`	// 구분 문자

		WalletAddr				string `json:"walletAddr" form:"walletAddr" valid:"required"`	// 에코 지갑 아이디
		PenaltySeed				string `json:"penaltySeed" form:"penaltySeed" valid:"required"`	//에코 지갑 주소
		PenaltySeedSum				string `json:"penaltySeedSum" form:"penaltySeedSum" valid:"required"`		// 정산량

	}
	
) 


// 서비스 전체 적용 상수 -  관리 항목으로 기능 개선 필요
const (

	ChainCodeId = "mcc"

)

  
// 호출 서비스별 URL 맵핑
func SetSvcUrl(g *echo.Group, fbc *blockchain.FabricSetup) {

	// 씨앗 선물
	g.PUT("/SSG", func(c echo.Context ) error {
		return ssg(c, fbc)
	})

	// 개별 DAD
	g.PUT("/DDW", func(c echo.Context ) error {
		return ddw(c, fbc)
	})

	// ECO 정산
	g.PUT("/DDE", func(c echo.Context ) error {
		return dde(c, fbc)
	})

	// 씨앗 일괄 선물
	g.PUT("/SSGM", func(c echo.Context ) error {
		return ssgm(c, fbc)
	})	

	// 멀티 개인별 DAD
	g.PUT("/DDWM", func(c echo.Context ) error {
		return ddwm(c, fbc)
	})
}


// 씨앗 전송
func ssg(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(STRT_SSG)
	util.BindParam(param, c)


	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID
	ccst.CcFnc = "SSG"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{param.SenderAddr, param.RecipientAddr,  param.BaseAmount, param.PenaltyAmount} )			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}

// 개별 사용자 지갑 정산 처리
func ddw(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(STRT_DAD_WLT)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID
	ccst.CcFnc = "DDW"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{param.WalletAddr, param.PenaltySeed, param.PenaltySeedSum})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}



// 에코 지갑 정산 처리
func dde(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	
	// 파라메터 맵핑
	param := new(STRT_DAD_ECO)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID
	ccst.CcFnc = "DDE"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{param.EcoID, param.EcoAddr, param.Amount})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}



// 멀티 씨앗 전송
func ssgm(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(STRT_SSGM)
	util.BindParam(param, c)


	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID
	ccst.CcFnc = "SSGM"			// 체인코드 함수명
	ccst.SetCCArgs(param.CoreSign.RHash,param.CoreSign.SHash,param.CoreSign.XValue,param.CoreSign.YValue,param.RefID,param.ProcDate,param.Svc,param.SenderAddr,param.RecipientAddr,param.BaseAmount,param.PenaltyAmount,param.Delimiter)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}


// 멀티 개인별 사용자 지갑 정산 처리
func ddwm(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(STRT_DAD_WLTM)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID
	ccst.CcFnc = "DDWM"			// 체인코드 함수명
//	ccst.SetCCSignedArgs(param.COMMSIGN, []string{param.WalletAddr, param.PenaltySeed, param.PenaltySeedSum})			// 체인코드 파라메터
	ccst.SetCCArgs(param.CoreSign.RHash,param.CoreSign.SHash,param.CoreSign.XValue,param.CoreSign.YValue,param.RefID,param.ProcDate,param.Svc,param.WalletAddr, param.PenaltySeed, param.PenaltySeedSum,param.Delimiter)			// 체인코드 파라메터


	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}