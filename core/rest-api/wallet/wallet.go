package wallet

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/echo-rest-server/util"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	CM_STRT "github.com/mycreditchain/rest-api/struct"
)


// 파라메터 구조체 - 적용되는 파라메터 타입에 대한 설정 포함
type (
	


	CreateStrt struct {
		Identity		string `json:"identity" form:"identity" valid:"required"`	// 사용자유니크키 - 아이디, 이메일등
		WalletType		string `json:"walletType" form:"walletType" valid:"required"`	// 지갑 타잎 - User / System / Eco 
		TimezoneOffset	string `json:"timezoneOffset" form:"timezoneOffset" valid:"required,int"`	// Difference from UTC
	}
	
	SearchStrt struct {
		SearchType		string `json:"searchType" query:"searchType" valid:"required"`	// 조회타입 - “identity” or “address”
		UserPath		string `json:"userPath" query:"userPath" valid:"required"`	// 사용자정의주소
	}


	TimeZoneStrt struct {

		CM_STRT.COMMSIGN

		WalletAddr		string `json:"walletAddr" form:"walletAddr" valid:"required"`	// 지갑 주소
		TimezoneOffset	string `json:"timezoneOffset" form:"timezoneOffset" valid:"required,int"`	// Difference from UTC

	}

	UpdateStrt struct {

		TimeZoneStrt

		Identity		string `json:"identity" form:"identity" valid:"required"`	// 사용자유니크키 - 아이디, 이메일등
		WalletType		string `json:"walletType" form:"walletType" valid:"required"`	// 지갑 타잎 - User / System / Eco 
		Amount			string `json:"amount" form:"amount" valid:"required"`	// 보내는양

		QueueKey			string `json:"queueKey" form:"queueKey"`	// 마케팅 refID
		
	}


	SyncWallet struct {

		CM_STRT.COMM


		//RefID				string `json:"refID" form:"refID" valid:"required"`	//refID
		//ProcDate			string `json:"procDate" form:"procDate" valid:"cstDateTime,required"`	//처리날짜 - 포맷 - 2018-10-12 14:15:59
		//Svc					string `json:"svc" form:"svc" valid:"required"`	//서비스 코드
		//ProcType			string `json:"procType" form:"procType"`	// 처리타입

		
		WalletAddr		string `json:"walletAddr" form:"walletAddr" valid:"required"`	// 지갑 주소
		StartRefID	string `json:"startRefID" form:"startRefID" valid:"required"`	// 첫번째 참조키
		EndRefID	string `json:"endRefID" form:"endRefID" valid:"required"`	// 마지막 참조키
		PaymentWalletYn	string `json:"paymentWalletYn" form:"paymentWalletYn"`	// 결제 지갑 여부
	}
)


// 서비스 전체 적용 상수 -  관리 항목으로 기능 개선 필요
const (

	ChainCodeId = "mcc"
)

// 호출 서비스별 URL 맵핑
func SetSvcUrl(g *echo.Group, fbc *blockchain.FabricSetup) {
	
	// 지갑 생성
	g.POST("", func(c echo.Context ) error {
		return create(c, fbc)
	})

	// 개별 지갑 조회
	g.GET("/:searchType/:userPath", func(c echo.Context ) error {
		return search(c, fbc)
	})


	// 지갑 설정
	g.PUT("", func(c echo.Context ) error {
		return update(c, fbc)
	})

	// 지갑 timezone변경
	g.PUT("/timezone", func(c echo.Context ) error {
		return updateTimezoneOffset(c, fbc)
	})


	// 지갑 동기화
	g.PUT("/sync", func(c echo.Context ) error {
		return sync(c, fbc)
	})
}

// 재갑 설정
func update(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(UpdateStrt)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId					// 체인코드 ID
	ccst.CcFnc = "restructWallet"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{ param.Identity, param.WalletAddr, param.WalletType, param.TimezoneOffset, param.Amount, param.QueueKey})			// 체인코드 파라메터
	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}

// 지갑 타임 존 변경
func updateTimezoneOffset(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(TimeZoneStrt)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId					// 체인코드 ID
	ccst.CcFnc = "updateTimezoneOffset"			// 체인코드 함수명
	ccst.SetCCSignedArgs(param.COMMSIGN, []string{ param.WalletAddr, param.TimezoneOffset})			// 체인코드 파라메터
	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}

// 지갑 생성
func create(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(CreateStrt)
	util.BindParam(param, c)




	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId					// 체인코드 ID
	ccst.CcFnc = "createMCCWallet"			// 체인코드 함수명
	ccst.SetCCArgs(param.Identity, param.WalletType, param.TimezoneOffset)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}


// 개별 지갑 정보 조회 - 타패키지 호출용
func Search(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	return search(c, fbc)
}

// 개별 지갑 정보 조회
func search(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	
	// 파라메터 맵핑
	param := new(SearchStrt)
	util.BindParam(param, c)

	// 체인코드 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId 					// 체인코드 ID
	ccst.CcFnc = "getWalletInfo"   			// 체인코드 함수명
	ccst.SetCCArgs(param.SearchType, param.UserPath)			// 체인코드 파라메터
	
	// 체인코드 호출 및 리턴
	return util.CcQueryRes(c, fbc, ccst)
}

// 지갑동기화
func sync(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(SyncWallet)
	util.BindParam(param, c)

	
	// 유효성체크
	// 1. 동일지갑 및 지갑 파라메터 존재
	// 2. 금액 0 및 존재

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId
	ccst.CcFnc = "syncWallet"			// 체인코드 함수명
	//ccst.SetCCSignedArgs(param.COMM, []string{ param.WalletAddr, param.StartRefID, param.EndRefID})			// 체인코드 파라메터
	ccst.SetCCCommArgs(param.COMM, []string{ param.WalletAddr, param.StartRefID, param.EndRefID, param.PaymentWalletYn})			// 체인코드 파라메터
	//ccst.SetCCArgs(param.Svc, param.Svc, param.ProcDate, param.WalletAddr, param.StartRefID, param.EndRefID)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}