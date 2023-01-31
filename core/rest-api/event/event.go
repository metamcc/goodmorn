package event

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	"github.com/mycreditchain/echo-rest-server/util"
	CM_STRT "github.com/mycreditchain/rest-api/struct"
	//"strconv"
)

// 파라메터 구조체 - 적용되는 파라메터 타입에 대한 설정 포함
type (


	// 이벤트 열매 이체
	STRT_EV_F2F struct {

		CM_STRT.COMMSIGN

		MarketingAddr		string `json:"marketingAddr" form:"marketingAddr" valid:"required"`	// 마케팅 지갑주소
		UserAddr			string `json:"userAddr" form:"userAddr" valid:"required"`	//  받는이 지갑주소
		EventCd				string `json:"eventCd" form:"eventCd" valid:"required"`	//  이벤트코드
		Amount				string `json:"amount" form:"amount" valid:"required"`	// 이벤트 열매량
	}

) 


// 서비스 전체 적용 상수 -  관리 항목으로 기능 개선 필요
const (

	ChainCodeId = "mcc"

)


// 호출 서비스별 URL 맵핑
func SetSvcUrl(g *echo.Group, fbc *blockchain.FabricSetup) {

	// 시스템 열매 선물
	g.PUT("/F2F", func(c echo.Context ) error {
		return ev_f2f(c, fbc)
	})

	// 시스템 열매 선물
	g.PUT("/S2S", func(c echo.Context ) error {
		return ev_f2f(c, fbc)
	})
}


// 이벤트 열매 이체
func ev_f2f(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(STRT_EV_F2F)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "EV_F2F"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{param.EventCd, param.MarketingAddr, param.UserAddr, param.Amount})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}


// 시스템 씨앗 선물
func ev_s2s(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(STRT_EV_F2F)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "EV_S2S"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{param.MarketingAddr, param.UserAddr, param.Amount})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}