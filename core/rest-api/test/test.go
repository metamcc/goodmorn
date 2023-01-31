package test

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/echo-rest-server/util"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	"github.com/mycreditchain/rest-api/system"
	//"github.com/mycreditchain/echo-rest-server/service/wallet"
	"github.com/labstack/gommon/random"
	"math/rand"
	"fmt"
	"strconv"
	"time"
	"github.com/mycreditchain/common/msg"
	"errors"
//	CM_STRT "github.com/mycreditchain/rest-api/struct"
)


type (

	Size struct {
		Key		string `json:"key" form:"key" `	//키
		Value		string `json:"value" form:"value" `	//값

	}

	Wallet1 struct {
		Identity		string `json:"identity" form:"identity" valid:"required"`	// 사용자유니크키 - 아이디, 이메일등
		WalletType		string `json:"walletType" form:"walletType" valid:"required"`	// 지갑 타잎 - User / System / Eco 
		TimezoneOffset	int `json:"timezoneOffset" form:"timezoneOffset" valid:"int"`	// Difference from UTC
	}

	MakeData struct {
		WalletAddr		string `json:"walletAddr" form:"walletAddr" valid:"required"`	// 지갑 주소
		Fruit				string `json:"fruit" form:"fruit" valid:"required"`	//이체양
		Typ				string `json:"typ" form:"typ" valid:"required"`	//지갑 종류

		
	}


	Test struct {
		Identity		string `json:"identity" form:"identity" valid:"required"`	// 사용자유니크키 - 아이디, 이메일등
		Fruit				string `json:"fruit" form:"fruit" valid:"required"`	//이체양
	}

		// 열매 이체
		F2F struct {
			//CM_STRT.COMMSIGN
	
			RefID				string `json:"refID" form:"refID" valid:"required"`	//refID
			Svc					string `json:"svc" form:"svc" valid:"required"`	//서비스 코드
			ProcDate			string `json:"procDate" form:"procDate" valid:"cstDateTime,required"`	//처리날짜 - 포맷 - 2018-10-12 14:15:59

			SenderAddr			string `json:"senderAddr" form:"senderAddr" valid:"required"`	//보내는이 지갑
			RecipientAddr		string `json:"recipientAddr" form:"recipientAddr" valid:"required"`	//받는이 지갑
			Amount				string `json:"amount" form:"amount" valid:"required"`	// 보내는양
			OrgRefID			string `json:"orgRefID" form:"orgRefID" `	// 원본참조키
			
		}
)



// 호출 서비스별 URL 맵핑
func SetSvcUrl(g *echo.Group, fbc *blockchain.FabricSetup) {

	// 생성용 테스트 용
	g.PUT("/createKeyValue", func(c echo.Context ) error {
		return createKeyValue(c, fbc)
	})

	// 지갑 생성 테스트 용
	g.GET("/createWallet", func(c echo.Context ) error {
		return createWallet(c, fbc)
	})
	
	// 지갑 임시 생성 테스트 용
	g.GET("/sendSeed", func(c echo.Context ) error {
		return sendSeed(c, fbc)
	})


	// 재처리 테스트 용
	g.PUT("/retry", func(c echo.Context ) error {
		return retry(c, fbc)
	})

	// 다용도 테스트 용
	g.PUT("/test", func(c echo.Context ) error {

		return test(c, fbc)
		
	})


		// 다용도 테스트 용
		g.GET("/test2", func(c echo.Context ) error {
			return simpleAPI(c, fbc)
		})

				// 다용도 테스트 용
				g.PUT("/makeData", func(c echo.Context ) error {
					return makeData(c, fbc)
				})


	// 체인 코드 설치 및 업그레이드
	g.POST("/upgradeCC", func(c echo.Context ) error {
		return system.UpgradeCC(c, fbc)
	})

		// 열매 결제
		g.PUT("/payment", func(c echo.Context ) error {
			return payment(c, fbc)
		})
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
	ccst.CcId = "mcc"
	ccst.CcFnc = "payFruit2"			// 체인코드 함수명
	ccst.SetCCArgs( param.RefID,param.ProcDate,param.Svc,  param.SenderAddr, param.RecipientAddr, param.Amount)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}

//  테스트 함수
func test(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	act := c.FormValue("act")



	if "1" == act {

		return retry(c, fbc)
	
	} else if "2" == act {
	
		return makeData2(c, fbc)
	}

	resultSt :=	msg.GetMsgStruct("RS_RES_N01")


	return  c.JSON(200, resultSt)
}



// 용량 테스트
func makeData2(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(Test)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc"					// 체인코드 ID
	ccst.CcFnc = "makeData2"			// 체인코드 함수명
	//ccst.SetCCArgs(strconv.Itoa(param.Key),strconv.Itoa(param.Value))			// 체인코드 파라메터
	ccst.SetCCArgs(param.Identity, param.Fruit)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}

// 용량 테스트
func makeData(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(MakeData)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc"					// 체인코드 ID
	ccst.CcFnc = "makeData"			// 체인코드 함수명
	//ccst.SetCCArgs(strconv.Itoa(param.Key),strconv.Itoa(param.Value))			// 체인코드 파라메터
	ccst.SetCCArgs(param.WalletAddr, param.Fruit, param.Typ)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}

// 단순 호출 테스트
func simpleAPI(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	resultSt :=	msg.GetMsgStruct("RS_RES_N01")
	fmt.Println("11111111111111111111111")

	//res, _ := util.GetTransactionByTxId("4adc6cf04a25146581abd68274d3b6260026c9344ad0230c42d16a2d2b027ba9")

	//system.GetTransactionByTxId("4adc6cf04a25146581abd68274d3b6260026c9344ad0230c42d16a2d2b027ba9")

	//fmt.Println(res)
	fmt.Println("2222222222222222222222")
	return  c.JSON(200, resultSt)
}

// 지갑 생성 테스트
func eventCreateWallet(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	
	// 파라메터 맵핑
	param := new(Wallet1)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc"					// 체인코드 ID
	ccst.CcFnc = "createMCCWallet"			// 체인코드 함수명
	ccst.SetCCArgs(param.Identity,param.WalletType,"0")			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}






// 용량 테스트
func createKeyValue(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	fmt.Println("createKeyValue")

	// 파라메터 맵핑
	param := new(Size)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc"					// 체인코드 ID
	ccst.CcFnc = "createKeyValue"			// 체인코드 함수명
	//ccst.SetCCArgs(strconv.Itoa(param.Key),strconv.Itoa(param.Value))			// 체인코드 파라메터
	ccst.SetCCArgs(param.Key,param.Value)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}





// 재처리 테스트
func retry(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	typ := c.FormValue("typ")
	var ranCnt int
	if typ == "1" {
		ranCnt,_ = strconv.Atoi(c.FormValue("ranCnt"))
	} else {
		ranCnt = randomInt(0, 100)

		time.Sleep(70 * time.Second)
	}
	if ranCnt < 30 {
		err1 := errors.New("MVCC_READ_CONFLICT")
		msg.CallPanicWithMsg(err1, "RETRY")
		return err1
	} else if ranCnt < 60   {
	
		err1 := errors.New("CC_WLT_E25")
		msg.CallPanicWithMsg(err1, "RETRY")
		return err1
	} else if ranCnt < 95 {

		err1 := errors.New("CC_WLT_E02")
		msg.CallPanicWithMsg(err1, "RETRY")
		return err1

	} else {

		resultSt :=	msg.GetMsgStruct("RS_RES_N01")


		return  c.JSON(200, resultSt)
	}
}

// 지갑 생성 테스트
func createWallet(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	/*
	fmt.Println("createKeyValue")

	// 파라메터 맵핑
	param := new(Size)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc"					// 체인코드 ID
	ccst.CcFnc = "createKeyValue"			// 체인코드 함수명
	//ccst.SetCCArgs(strconv.Itoa(param.Key),strconv.Itoa(param.Value))			// 체인코드 파라메터
	ccst.SetCCArgs(param.Key,param.Value)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
*/
	q := c.Request().URL.Query()
//	q.Set("identity", random.String(32))

	q.Set("ccId", "mcc")
	q.Set("ccFnc", "createMCCWallet")
	q.Set("param1", random.String(32))
	q.Set("param2", "user")
	q.Set("param3", "0")

	c.Request().URL.RawQuery = q.Encode() 

	return system.CcInvoke(c, fbc)

}



// 씨앗 전송 테스트
func sendSeed(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
/*
	baseAddr1 := "testWalletAddress"
	baseAddr2 := baseAddr1+"0"
	lstAddr3 := baseAddr2

	
	lstAddr1 := baseAddr1+strconv.Itoa(randomInt(1,9))
	lstAddr2 := baseAddr1+strconv.Itoa(randomInt(1,9))
//	fmt.Println(lstAddr1 + " /" + lstAddr2 + " /" + lstAddr3)
	if lstAddr1 == lstAddr2 {
		//lstAddr2 = lstAddr3
		lstAddr3 = baseAddr2
	} else {
		lstAddr3 = lstAddr2
	}
*/
	q := c.Request().URL.Query()

	q.Set("ccId", "mycreditchain-service")
	q.Set("ccFnc", "sendGiftSeeds")

	q.Set("param1", "testWalletAddress0")
	q.Set("param2", "testWalletAddress1")
	q.Set("param3", "1")
	q.Set("param4", "ac280f50f894743bce69f3eca8f0d87f90e040d2a1f63cca79d92d2524cca926")
	q.Set("param5", "100113774347558763635196620777057547744106511837449035806014434766006542942365")
	q.Set("param6", "102687051746910574167914069359289101150310670635673484478773643237323224056416")
	q.Set("param7", "2018-09-10 00:00:00")
	c.Request().URL.RawQuery = q.Encode() 

	//return dad.SendGiftSeeds(c, fbc)

	return system.CcInvoke(c, fbc)
}

// min ~ max 사이의 랜덤 숫자 취득
func randomInt(min, max int) int {
	
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

    return min + r1.Intn(max+1-min)
}