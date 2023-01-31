package system

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/echo-rest-server/util"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	"github.com/pkg/errors"
	CM_STRT "github.com/mycreditchain/rest-api/struct"
	"encoding/json"
	"github.com/mycreditchain/common/msg"
)

// bind용 구조체 - 적용되는 파라메터 타입에 대한 설정 포함
type System struct {
	//CCFnc  string // 고정항목 : 체인코드 함수
	//Id  string `json:"id" form:"id" query:"id"`

	CnId  			string `json:"cnId" form:"cnId" valid:"required"` 
	CcId			string `json:"ccId" form:"ccId" valid:"required"` 
	CcPath			string `json:"ccPath" form:"ccPath" valid:"required"`
	//CcGoPath		string `json:"ccGoPath" form:"ccGoPath" valid:"required"`
	CcVs			string `json:"ccVs" form:"ccVs" valid:"required"`
	CcInit			string `json:"ccInit" form:"ccInit" valid:"required"`
	CcPolicy		string `json:"ccPolicy" form:"ccPolicy" query:"ccPolicy"`
}

type CC struct {

	CcId		string `json:"ccId" form:"ccId" valid:"required"` 					// 체인코드 아이디
	CcFnc		string `json:"ccFnc" form:"ccFnc" query:"ccFnc" valid:"required`	// 체인코드 함수
	Param1		string `json:"param1" form:"param1" query:"param1"`	// 파라메터1
	Param2		string `json:"param2" form:"param2" query:"param2"`	// 파라메터2
	Param3		string `json:"param3" form:"param3" query:"param3"`	// 파라메터3
	Param4		string `json:"param4" form:"param4" query:"param4"`	// 파라메터4
	Param5		string `json:"param5" form:"param5" query:"param5"`	// 파라메터5
	Param6		string `json:"param6" form:"param6" query:"param6"`	// 파라메터6
	Param7		string `json:"param7" form:"param7" query:"param7"`	// 파라메터7
	Param8		string `json:"param8" form:"param8" query:"param8"`	// 파라메터8
	Param9		string `json:"param9" form:"param9" query:"param9"`	// 파라메터9
	Param10		string `json:"param10" form:"param10" query:"param10"`	// 파라메터10
}

type Alive struct {

	RefID				string `json:"refID" form:"refID" valid:"required`	//refID
	RestServer				string `json:"restServer" form:"restServer" valid:"required`	//refID
}


type CoreChkBase struct {

	ProcDate			string `json:"procDate" form:"procDate" valid:"required,cstDateTime"`	//처리날짜 - 포맷 - 2018-10-12 14:15:59
	Svc					string `json:"svc" form:"svc" valid:"numeric" valid:"required`	//서비스 코드
	RefID				string `json:"refID" form:"refID" valid:"required`	//refID
	TxID				string `json:"txID" form:"txID"`	//txID

	
	
}

type CoreChkSeed struct {

	ProcDate			string `json:"procDate" form:"procDate" valid:"required,cstDateTime"`	//처리날짜 - 포맷 - 2018-10-12 14:15:59
	SenderAddr			string `json:"senderAddr" form:"senderAddr" valid:"required"`	//보내는이 지갑주소
	RecipientAddr		string `json:"recipientAddr" form:"recipientAddr" valid:"required"`	//받는이 지갑주소

	RefID				string `json:"refID" form:"refID" valid:"required`	//refID
	TxID				string `json:"txID" form:"txID"`	//txID
}


type CoreChkTx struct {
	TxID				string `json:"txID" form:"txID" valid:"required"`	//txID
}


type CCMultipleData struct {


	CcFnc		string `json:"ccFnc" form:"ccFnc" query:"ccFnc" valid:"required`	// 체인코드 함수
	Param1		string `json:"param1" form:"param1" query:"param1"`	// 파라메터1
	Param2		string `json:"param2" form:"param2" query:"param2"`	// 파라메터2
	Param3		string `json:"param3" form:"param3" query:"param3"`	// 파라메터3
	Param4		string `json:"param4" form:"param4" query:"param4"`	// 파라메터4
	Param5		string `json:"param5" form:"param5" query:"param5"`	// 파라메터5
	Param6		string `json:"param6" form:"param6" query:"param6"`	// 파라메터6
	Param7		string `json:"param7" form:"param7" query:"param7"`	// 파라메터7
	Param8		string `json:"param8" form:"param8" query:"param8"`	// 파라메터8
	Param9		string `json:"param9" form:"param9" query:"param9"`	// 파라메터9
	Param10		string `json:"param10" form:"param10" query:"param10"`	// 파라메터10
}

const (

	ChainCodeId = "mcc"

)

// 호출 서비스별 URL 맵핑
func SetSvcUrl(g *echo.Group, fbc *blockchain.FabricSetup) {

	// 체인 코드 query
	g.Any("/query", func(c echo.Context ) error {
		return ccQuery(c, fbc)
	})

	// 체인 코드 invoke
	g.Any("/invoke", func(c echo.Context ) error {
		return ccInvoke(c, fbc)
	})

	// 채널 구성
	g.POST("/setupChannel", func(c echo.Context ) error {
		return setupChannel(c, fbc)
	})

	// 코어 체크
	g.POST("/alive", func(c echo.Context ) error {
		return postAlive(c, fbc)
	})

	// 코어 체크
	g.GET("/alive", func(c echo.Context ) error {
		return getAlive(c, fbc)
	})


	// queue Key 로 조회
	g.GET("/queue", func(c echo.Context ) error {
		return getTxInfoByQueueKey(c, fbc)
	})

	// 트랜잭션 ID 로 조회
	g.GET("/tx/:txID", func(c echo.Context ) error {
		return getPayloadByTx(c, fbc)
	})


	g.POST("/procCommMultipleData", func(c echo.Context ) error {
		return procCommMultipleData(c, fbc)
	})

	// 체인 코드 설치 및 인스턴스 생성
	g.POST("/setupCC", func(c echo.Context ) error {
		if fbc.RestServer != "master" {
			resultSt :=	msg.GetMsgStruct("RS_RES_E99")

			return  c.JSON(404, resultSt) 
		}
		return setupCC(c, fbc)
	})

	// 체인 코드 설치 및 업그레이드
	g.POST("/upgradeCC", func(c echo.Context ) error {
		if fbc.RestServer != "master" {
			resultSt :=	msg.GetMsgStruct("RS_RES_E99")

			return  c.JSON(404, resultSt) 
		}
		return UpgradeCC(c, fbc)
	})
}

// 큐 키로 tx 정보 조회 조회
func getTxInfoByQueueKey(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 코어 검증 조회

	resData, resErr := ChkCoreRefID(c, fbc)

	c.Logger().Debug("system--getTxInfoByQueueKey =", resData)
	c.Logger().Debug("system--getTxInfoByQueueKey =", resErr)

	if resErr != nil {  // 에러면(키가 없으면) -> 신규 재처리 진행 
		
		c.Logger().Debug("코어 검증키 없음  ===", resData)

		resultSt := msg.SetErrMsgStruct("RS_RES_E03", resErr.Error())
	
		return  c.JSON(200, resultSt)
		
	} else { // 에러가 아니면(키가 존재시) -> 전문 역동기화로 내려줌

		c.Logger().Debug("코어 검증키 존재 ===", resData)

		//return GetPayloadByTxId(c, fbc, resData)

		//resultSt := msg.SetErrMsgStruct("RS_RES_N01", resData)
	
		return  util.CcRes(c, resData, resErr)


	}
}




// 트랜잭션 아이디 조회
func getPayloadByTx(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(CoreChkTx)
	util.BindParam(param, c)

	resData, err := fbc.GetTransactionByTxId(param.TxID)

	return util.CcRes(c, resData, err) 
}


// 공통 일괄 처리 
func procCommMultipleData(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(CC)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = ChainCodeId 	// 체인코드 ID
	ccst.CcFnc = param.CcFnc    		// 체인코드 함수명
	ccst.SetCCArgs(param.Param1,param.Param2,param.Param3,param.Param4,param.Param5,param.Param6,param.Param7,param.Param8,param.Param9,param.Param10)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)

}

// 트랜잭션 페이로드 취득
func GetPayloadByTxId(c echo.Context, fbc *blockchain.FabricSetup, resData string) (err error) {

	ccLst := CM_STRT.RES_CC_LST{}
	json.Unmarshal([]byte(resData), &ccLst)

	if len(ccLst.CoreChkInfo) > 0 {

		refID := ccLst.CoreChkInfo[0].RefID
		txID := ccLst.CoreChkInfo[0].TxID
	
		resData, err = fbc.GetTransactionByTxId(txID)

		if err != nil {
			return err
		}
	
		resData = "{\"retryID\":\""+refID+"\","+resData[1:]

		return util.CcRes(c, resData, err)

	} else {
		return errors.WithMessage(err, "CORE REF KEY ERROR")
	}

}






// 코어 검증 체크
func ChkCoreRefID(c echo.Context, fbc *blockchain.FabricSetup) (result string, err error) {

	svc := c.FormValue("svc")
	
	if svc == "2" { // 씨앗 선물
		return chkCoreSeed(c, fbc)
	} else {
		return chkCoreBase(c, fbc)
	}
}



// 코어 검증 씨앗 체크
func chkCoreSeed(c echo.Context, fbc *blockchain.FabricSetup) (result string, err error) {

	// 파라메터 맵핑
	param := new(CoreChkSeed)
	util.BindParam(param, c)

	// 체인코드 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc" 	// 체인코드 ID
	ccst.CcFnc = "getCoreSeedCompKey"    		// 체인코드 함수명

	procDate := param.ProcDate
	ccst.SetCCArgs(procDate, param.SenderAddr, param.RecipientAddr, param.RefID, "")			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcQueryData(c, fbc, ccst)
}


// 코어 검증 기본형 체크
func chkCoreBase(c echo.Context, fbc *blockchain.FabricSetup) (result string, err error) {


	// 파라메터 맵핑
	param := new(CoreChkBase)
	util.BindParam(param, c)

	// 체인코드 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc" 	// 체인코드 ID
	ccst.CcFnc = "getCoreCompKey"    		// 체인코드 함수명

	procDate := param.ProcDate
	ccst.SetCCArgs(procDate[0:4], procDate[5:7], procDate[8:10], param.Svc, param.RefID, param.TxID)			// 체인코드 파라메터

	
	// 체인코드 호출 및 리턴
	return util.CcQueryData(c, fbc, ccst)
}

// health 체크
func getAlive(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	
	// 체인코드 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc" 	// 체인코드 ID
	ccst.CcFnc = "aliveStatusQuery"    		// 체인코드 함수명
	ccst.SetCCArgs("none", fbc.RestServer)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcQueryRes(c, fbc, ccst)
}


// health 체크
func postAlive(c echo.Context, fbc *blockchain.FabricSetup) (err error) {


	// 파라메터 맵핑
	param := new(Alive)
	util.BindParam(param, c)

	// 체인코드 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = "mcc" 	// 체인코드 ID
	ccst.CcFnc = "aliveStatusQuery"    		// 체인코드 함수명
	ccst.SetCCArgs(param.RefID, fbc.RestServer)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcQueryRes(c, fbc, ccst)

}

// 체인코드 설치
func setupCC(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(System)
	util.BindParam(param, c)

	// 체인코드 호출
	err2 := fbc.SetupCC(param)

	// 체인 코드 결과 리턴
	return util.CcRes(c, "", err2)
}

// 체인코드 업그레이드
func UpgradeCC(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(System)
	util.BindParam(param, c)

	// 체인코드 호출
	err2 := fbc.UpgradeCC(param)

	// 체인 코드 결과 리턴
	return util.CcRes(c, "", err2)
}


func CcQuery(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	return ccQuery(c, fbc)
}

func CcInvoke(c echo.Context, fbc *blockchain.FabricSetup) (err error) {
	return ccInvoke(c, fbc)
}

// 채널 구성
func setupChannel(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(System)
	util.BindParam(param, c)

	// 체인코드 호출
	err2 := fbc.SetupChannel(param)

	// 체인 코드 결과 리턴
	return util.CcRes(c, "", err2)
}



// 체인 코드 query
func ccQuery(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(CC)
	util.BindParam(param, c)

	// 체인코드 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = param.CcId 	// 체인코드 ID
	ccst.CcFnc = param.CcFnc    		// 체인코드 함수명
	ccst.SetCCArgs(param.Param1,param.Param2,param.Param3,param.Param4,param.Param5,param.Param6,param.Param7,param.Param8,param.Param9,param.Param10)			// 체인코드 파라메터
	
	// 체인코드 호출 및 리턴
	return util.CcQueryRes(c, fbc, ccst)
}



// 체인 코드 invoke
func ccInvoke(c echo.Context, fbc *blockchain.FabricSetup) (err error) {

	// 파라메터 맵핑
	param := new(CC)
	util.BindParam(param, c)

	// 체인 코드 호출 정보 구성
	ccst := new(blockchain.ChainCodeInfo)
	ccst.CcId = param.CcId 	// 체인코드 ID
	ccst.CcFnc = param.CcFnc    		// 체인코드 함수명
	ccst.SetCCArgs(param.Param1,param.Param2,param.Param3,param.Param4,param.Param5,param.Param6,param.Param7,param.Param8,param.Param9,param.Param10)			// 체인코드 파라메터

	// 체인코드 호출 및 리턴
	return util.CcInvokeRes(c, fbc, ccst)
}
/*
func (setup *FabricSetup) InstallAndInstantiateCC() error {

	// Create the chaincode package that will be sent to the peers
	fmt.Println("chain code path = "+setup.ChaincodePath)

	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("ccPkg created")

	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: setup.ChainCodeVersion, Package: ccPkg}
	fmt.Println(installCCReq)
	
	_, err = setup.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")
	}
	fmt.Println("Chaincode installed")

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.chainhero.io"})

	resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodeGoPath, Version: setup.ChainCodeVersion, Args: [][]byte{[]byte("initLedger")}, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")

	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	// Creation of the client which will enables access to our channel events
	setup.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	fmt.Println("Chaincode Installation & Instantiation Successful")
	return nil
}
*/