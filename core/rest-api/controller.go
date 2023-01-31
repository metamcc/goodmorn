package api

import (
	"github.com/labstack/echo"
	"github.com/mycreditchain/rest-api/system"
	"github.com/mycreditchain/rest-api/fruit"
	"github.com/mycreditchain/rest-api/wallet"
	"github.com/mycreditchain/rest-api/test"
	"github.com/mycreditchain/rest-api/dad"
	"github.com/mycreditchain/rest-api/mng"
	"github.com/mycreditchain/rest-api/event"
	"github.com/mycreditchain/echo-rest-server/blockchain"
	//"github.com/mycreditchain/echo-rest-server/util"
	"github.com/mycreditchain/common/const"
	//"encoding/json"
	//"fmt"
)


// URL 그룹 맵핑
func SetUrlGroup(e *echo.Echo, fbc *blockchain.FabricSetup) {

/*
	// test 관련
	group99 := e.Group("/test")
	{
		test.SetSvcUrl(group99, fbc)
	}

			// 시스템 관련
		group0 := e.Group("/system")
		{
			system.SetSvcUrl(group0, fbc)
		}


				// 관리
		group4 := e.Group("/mng")
		{
			mng.SetSvcUrl(group4, fbc)
		}
*/

		// TEST 용
		/*
		group99 := e.Group("/test", func (next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
			
			//	target := url.QueryEscape(c.Request().URL.Path)
			//	c.Redirect(http.StatusTemporaryRedirect, "/login?target="+target)
			//c.Logger().Debug("c.Request().URL ======", c.Request().URL)
			//c.Logger().Debug("c.Request().URL.Path ======", c.Request().URL.Path)

			//c.Logger().Debug("fbc.RestServer ======", fbc.RestServer)
				if fbc.RestServer != "master" {

					req := c.Request()

					if req.Method == echo.PUT || req.Method == echo.POST {

						path := c.ParamNames()

						//if len(path) > 0 {
							f := make(url.Values)
							for _, p := range path {
			
								f.Set(p,  c.Param(p))
							}

							req2,_ := http.NewRequest(req.Method, const_rest.MASTER_INTERNAL_ADDR+req.URL.Path, strings.NewReader(f.Encode()) )
							req2.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
							c.SetRequest(req2)

							resultSt :=	msg.GetMsgStruct("RS_RES_N01")


							return  c.JSON(200, resultSt)
						//}

					} else {
						return c.Redirect(http.StatusSeeOther, const_rest.MASTER_INTERNAL_ADDR+c.Request().URL.Path)
					}



				//	return c.Redirect(http.StatusSeeOther, const_rest.MASTER_INTERNAL_ADDR+c.Request().URL.Path)


				} else {
					return next(c)
				}

			}
		})
		{
			test.SetSvcUrl(group99, fbc)
		}
*/

		// test 관련
		group99 := e.Group("/test")
		{
			test.SetSvcUrl(group99, fbc)
		}

		// 시스템 관련
		group0 := e.Group("/system")
		{
			system.SetSvcUrl(group0, fbc)
		}


		// 관리
		group4 := e.Group("/mng")
		{
			mng.SetSvcUrl(group4, fbc)
		}


		// wallet 관련
		group1 := e.Group("/wallet", func (next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
			
				return chkCoreRefData(next, c, fbc)
			}
		})
		{
			wallet.SetSvcUrl(group1, fbc)
		}

		// DAD 관련
		group2 := e.Group("/dad", func (next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
			
				return chkCoreRefData(next, c, fbc)
			}
		})
		{
			dad.SetSvcUrl(group2, fbc)
		}
/*
		// 열매 관련
		group3 := e.Group("/fruit", setGroupPolicy)
		{
			fruit.SetSvcUrl(group3, fbc)
		}
*/


		group3 := e.Group("/fruit", func (next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
			
				return chkCoreRefData(next, c, fbc)
			}
		})
		{
			fruit.SetSvcUrl(group3, fbc)
		}



		// 이벤트
		group5 := e.Group("/event", func (next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
			
				return chkCoreRefData(next, c, fbc)
			}
		})
		{
			event.SetSvcUrl(group5, fbc)
		}

}


// 그룹 정책 핸들러 설정 - SAMPLE
func setGroupPolicy(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
	
	
		return next(c)

	}
}



// 블록체인 처리 여부 검증
func chkCoreRefData(next echo.HandlerFunc, c echo.Context, fbc *blockchain.FabricSetup) ( err error) {

	procType := c.FormValue("procType")

	//c.Logger().Debug("chkCoreRefData--procType =", procType)

	// 재처리로 온 경우
	if procType == const_rest.RETRY_KAFKA_TOPIC {

		// 코어 검증 조회
		resData, err := system.ChkCoreRefID(c, fbc)

		c.Logger().Debug("chkCoreRefData--ChkCoreRefID resData =   ", resData)
		c.Logger().Debug("chkCoreRefData--ChkCoreRefID err     =   ", resData)

		if err != nil {  // 에러면(키가 없으면) -> 신규 재처리 진행 
			c.Logger().Debug("신규 재처리 호출")
			// if cnt := util.CheckCoreRefTarget(err); cnt > 0 {
			// 	return next(c) 
			// } 
			//return err
			return next(c) 
			
		} else { // 에러가 아니면(키가 존재시) -> 전문 역동기화로 내려줌

			c.Logger().Debug("역동기화 호출")


			return system.GetPayloadByTxId(c, fbc, resData)
		}

	} 

	return next(c)
}



	/* 그룹 정책 적용
		g.Use(middleware.BasicAuth(func(username, password string) bool {
		if username == "joe" && password == "secret" {
			return true
		}
		return false
	}))
	*/
