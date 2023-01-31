package fruit

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	"github.com/mycreditchain/echo-rest-server/util"
	CM_STRT "github.com/mycreditchain/rest-api/struct"
	//"github.com/mycreditchain/rest-api/wallet"
	//"errors"
	//"github.com/mycreditchain/common/msg"
)

// 파라메터 구조체 - 적용되는 파라메터 타입에 대한 설정 포함
type (

	// 열매 이체
	F2F struct {
		CM_STRT.COMMSIGN

		SenderAddr			string `json:"senderAddr" form:"senderAddr" valid:"required"`	//보내는이 지갑
		RecipientAddr		string `json:"recipientAddr" form:"recipientAddr" valid:"required"`	//받는이 지갑
		Amount				string `json:"amount" form:"amount" valid:"required"`	// 보내는양
		OrgRefID			string `json:"orgRefID" form:"orgRefID" `	// 원본참조키
		
	}

	// 열매 -> MCC
	F2M struct {

		CM_STRT.COMMSIGN
		
		UserAddr		string `json:"userAddr" form:"userAddr" valid:"required"`	// 사용자  지갑
		EthereumAddr			string `json:"ethereumAddr" form:"ethereumAddr" valid:"required"`	// 사용자 이더리움 입금 지갑
		Amount				string `json:"amount" form:"amount" valid:"required"`	// 전체량
		BaseFee				string `json:"baseFee" form:"baseFee" valid:"required"`	// 기본수수료량
		ExtraFee				string `json:"extraFee" form:"extraFee" valid:"required"`	// 추가수수료량
	}

	// MCC -> 열매
	M2F struct {

		CM_STRT.COMMSIGN
		Amount				string `json:"amount" form:"amount" valid:"required"`	// 전환양
		EthereumAddr			string `json:"ethereumAddr" form:"ethereumAddr" valid:"required"`	// 사용자 이더리움 입금 지갑
		EthereumHash			string `json:"ethereumHash" form:"ethereumHash" valid:"required"`	// 이더리움 스캔 해쉬
		UserAddr		string `json:"userAddr" form:"userAddr" valid:"required"`	// 사용자  지갑
	}



	F2NF struct {
		CM_STRT.COMMSIGN

		SenderAddr			string `json:"senderAddr" form:"senderAddr" valid:"required"`	//보내는이 지갑
		TotalAmount				string `json:"totalAmount" form:"totalAmount" valid:"required"`	// 보내는 총양

		ReceverAddrs		string `json:"receiverAddrs" form:"receiverAddrs" valid:"required"`	//받는이 지갑

	}
) 


// 서비스 전체 적용 상수 -  관리 항목으로 기능 개선 필요
const (

	ChainCodeId = "mcc"

)


// 호출 서비스별 URL 맵핑
func SetSvcUrl(g *echo.Group, fbc *blockchain.FabricSetup) {

	// 열매 이체
	g.PUT("/F2F", func(c echo.Context ) error {
		return f2f(c, fbc)
	})

	// FRUIT -> MCC  전환 처리
	g.PUT("/F2M", func(c echo.Context ) error {
		return f2m(c, fbc)
	})	

	// MCC -> FRUIT  전환 처리
	g.PUT("/M2F", func(c echo.Context ) error {
		return m2f(c, fbc)
	})	

	// 열매 결제
	g.PUT("/payment", func(c echo.Context ) error {
		return payment(c, fbc)
	})


	// 열매 결제
	// g.PUT("/payFruit", func(c echo.Context ) error {
	// 	return payFruit(c, fbc)
	// })

	// // 열매 결제
	// g.PUT("/payFruit3", func(c echo.Context ) error {
	// 	return payFruit3(c, fbc)
	// })
}


// 열매 이체
func f2f(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(F2F)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "F2F"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{ param.SenderAddr, param.RecipientAddr, param.Amount})			// 체인코드 파라메터


	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}


// FRUIT -> MCC  전환 처리
func f2m(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(F2M)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "F2M"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{ param.EthereumAddr, param.UserAddr,  param.Amount,  param.BaseFee,  param.ExtraFee})			// 체인코드 파라메터


	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}


// MCC -> FRUIT  전환 처리
func m2f(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(M2F)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "M2F"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{param.Amount, param.EthereumAddr, param.EthereumHash, param.UserAddr})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}


// 열매 이체
func payment(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(F2F)
	util.BindParam(param, c)

	//param.COMMSIGN.

	//param.COMMSIGN.CommParam.Svc = "10"
	//commStrt.CommParam.Svc

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "payFruit"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{ param.SenderAddr, param.RecipientAddr, param.Amount, param.OrgRefID})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}

/*
// 대량 열매 이체
func payFruit3(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(F2NF)
	util.BindParam(param, c)

	
	// 유효성체크
	// 1. 동일지갑 및 지갑 파라메터 존재
	// 2. 금액 0 및 존재


	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "payFruit3"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{ param.SenderAddr, param.TotalAmount, param.ReceverAddrs})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}



// 열매 이체
func payFruit(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(F2F)
	util.BindParam(param, c)

	

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "payFruit"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{ param.SenderAddr, param.RecipientAddr, param.Amount})			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}
*/
