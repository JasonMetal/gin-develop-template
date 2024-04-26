// Package authController
// @Description:
package authController

import (
	"develop-template/app/entity/req/authReqEntity"
	"develop-template/app/entity/resp/authRespEntity"
	"develop-template/app/error/common"
	baseController "develop-template/app/http/controller"
	"develop-template/app/logic/authLogic"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	jsonIter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	//stringHelper "github.com/JasonMetal/submodule-support-go.git/helper/strings"
)

type controller struct {
	baseController.BaseController
}

// NewController
//
//	@Description:
//	@since: time
//	@param ctx
//	@return *controller
func NewController(ctx *gin.Context) *controller {
	return &controller{baseController.NewBaseController(ctx)}
}

// GetAuthProtocol
//
//	@Description:
//	@receiver c
//	@since: 2023/4/26_17:06
func (c *controller) GetAuthProtocol() {
	logic := authLogic.NewLogic(c.GCtx)
	resp := logic.GetCheckTokenData()
	c.Success(resp)
}

// TestAuth
//
//	@Description: {"code":0, "data":map[string]interface {}{"description":"恭喜Token验证成功！", "id":1, "now_time":1.682332330584e+12, "result":true}, "message":"GetCheckTokenData"}
//	@receiver c
//	@since: 2023/4/26_17:07
func (c *controller) TestAuth() {
	reqParams := authReqEntity.AuthInput{}
	if err := c.ShouldBindQuery(&reqParams); err != nil {
		c.Fail(common.ReqParamErr, c.GetValidMsg(err, &reqParams))
		return
	}

	//fmt.Println(stringHelper.Md5("123456"))

	var resPon authRespEntity.RespCheck
	token := reqParams.Token
	resp, err := http.Get("http://127.0.0.1:8080/auth/" + token + "/xxx?xx")
	if err != nil {
		fmt.Println("token auth url err :", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if len(body) > 0 {
		jsonIter.Unmarshal(body, &resPon)
		json.Unmarshal(body, &resPon)
		fmt.Printf("resPon :%#v\n", resPon)
		if resPon.Data.Result == false {
			err = fmt.Errorf("authorization failed")
			c.GCtx.JSON(0, "authorization failed")
			return
		}
	}
	fmt.Println("token is right!!!")
	c.GCtx.JSON(0, "token is right!!!")
}
