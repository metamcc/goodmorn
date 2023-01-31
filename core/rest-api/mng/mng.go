package mng

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	"github.com/mycreditchain/echo-rest-server/util"
	"github.com/mycreditchain/common/msg"
	"github.com/mycreditchain/echo-rest-server/handler"

	//"strconv"
)

// 파라메터 구조체 - 적용되는 파라메터 타입에 대한 설정 포함
type (


	GetConfig struct {
		Typ					string `json:"typ" form:"typ" valid:"required"`	// 타입 - msg,var,fnc
	}

	SetConfig struct {
		Typ					string `json:"typ" form:"typ" valid:"required"`	// 타입 - msg,var,fnc
		Action				string `json:"action" form:"action" valid:"required"`	// 처리기준 - 1 : 추가/수정, 0 : 삭제
		Key					string `json:"key" form:"key" `	// 키
		Value				string `json:"value" form:"value" `	// 값
	}

	

	GetEvent struct {
		EvtCd				string `json:"evtCd" form:"evtCd" valid:"required"`	// 이벤트 코드
	}


	SetEvent struct {
		EvtCd				string `json:"evtCd" form:"evtCd" valid:"required"`	// 이벤트코드
		EvtDt				string `json:"evtDt" form:"evtDt" valid:"required"`	// 이벤트기준일자
		EvtAmnt				string `json:"evtAmnt" form:"evtAmnt" valid:"required"`	// 이벤트열매량
		EvtCnt				string `json:"evtCnt" form:"evtCnt" valid:"required"`	// 이벤트참가횟수 - 0 - 매번 허용 , 1 : 유일 1회 허용,  2 - 2회 허용....
		Etc					string `json:"etc" form:"etc"`	// 이벤트 기타
	}



	Retry struct {
		RtyCode				string `json:"rtyCode" form:"rtyCode" valid:"required"`	// 재처리 대상 코드
		ActType				string `json:"actType" form:"actType"`	// 처리 구분
	}

) 


// 서비스 전체 적용 상수 -  관리 항목으로 기능 개선 필요
const (

	ChainCodeId = "confcc"

)

  
// 호출 서비스별 URL 맵핑
func SetSvcUrl(g *echo.Group, fbc *blockchain.FabricSetup) {

	// 체인코드 환경 정보 취득
	g.GET("/config/:typ", func(c echo.Context ) error {
		return getConfig(c, fbc)
	})

	// 체인코드 환경 정보 구성
	g.POST("/config", func(c echo.Context ) error {
		return setConfig(c, fbc)
	})


	// 체인코드 이벤트 정보 취득
	g.GET("/event/:evtCd", func(c echo.Context ) error {
		return getEevent(c, fbc)
	})

	// 체인코드 이벤트 정보 구성
	g.POST("/event", func(c echo.Context ) error {
		return setEevent(c, fbc)
	})


	// 재처리 대상 관리
	g.POST("/config/retry", func(c echo.Context ) error {
		return setKafkaRetry(c, fbc)
	})


	// 일괄 처리 관리
	g.POST("/config/batch", func(c echo.Context ) error {
		return setKafkaBatch(c, fbc)
	})
}



// 재처리 대상 관리
func setKafkaRetry(c echo.Context, fbc *blockchain.FabricSetup) (err error) {


	// 파라메터 맵핑
	param := new(Retry)
	util.BindParam(param, c)

	//c.Logger().Debugf("%+v", msgKey+logData)
	//c.Logger().Debug(param.RtyCode)
	//c.Logger().Debug(param.ActType)

	if param.ActType == "delete" { // 설정
		delete(handler.RETRY_TARGET, param.RtyCode)
	} else { // 삭제
		handler.RETRY_TARGET[param.RtyCode] = "RETRY"
	}

	resultSt :=	msg.GetMsgStruct("RS_RES_N01")


	return  c.JSON(200, resultSt)
}


// 일괄 처리 관리
func setKafkaBatch(c echo.Context, fbc *blockchain.FabricSetup) (err error) {



	
	resultSt :=	msg.GetMsgStruct("RS_RES_N01")


	return  c.JSON(200, resultSt)
}

// 체인코드 환경 정보 취득
func getConfig(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(GetConfig)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID

	if param.Typ == "msg" {
		ccst.CcFnc = "getCoreMsg"
	} else if param.Typ == "var" {
		ccst.CcFnc = "getCoreVar"
	} else if param.Typ == "fnc" {
		ccst.CcFnc = "getFncEnable"
	} 

	// 체인코드 호출 및 리턴
	return util.CcQueryRes(c, fbc, ccst)
}

// 체인코드 환경 정보 구성
func setConfig(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(SetConfig)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID

	// 체인코드 함수명
	if param.Typ == "msg" {
		ccst.CcFnc = "setCoreMsg"
	} else if param.Typ == "var" {
		ccst.CcFnc = "setCoreVar"
	} else if param.Typ == "fnc" {
		ccst.CcFnc = "setFncEnable"
	} 

	ccst.SetCCArgs(param.Action, param.Key, param.Value)			  // 체인코드 파라메터
	
	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}




// 체인코드 이벤트 정보 취득
func getEevent(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(GetEvent)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID

	ccst.CcFnc = "getCoreEvt"
	ccst.SetCCArgs(param.EvtCd)			  // 체인코드 파라메터


	// 체인코드 호출 및 리턴
	return util.CcQueryRes(c, fbc, ccst)
}


// 체인코드 이벤트 정보 구성
func setEevent(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(SetEvent)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId	// 체인코드 ID

	ccst.CcFnc = "setCoreEvt"
	ccst.SetCCArgs(param.EvtCd, param.EvtDt, param.EvtAmnt, param.EvtCnt, param.Etc)			  // 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}

