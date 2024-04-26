package authRouter

import (
	"develop-template/app/http/controller/api/authController"
	config2 "develop-template/helper/config"
	"github.com/gin-gonic/gin"
)

/*
*
router := gin.Default()

	router.GET("/Auth/:name/*action", func(c *gin.Context) {
		name := c.Param("name") // 可以获取路径中的 name 参数
		action := c.Param("action") // 可以获取 *action 之后的所有路径

		// 此时:
		// name = "abcedfg"
		// action = "sssss.dweb/"

		c.String(http.StatusOK, "name: %s, action: %s", name, action)
	})

router.Run()
*/
func RegisterAuth(router *gin.Engine) {
	domainMark := config2.GetDomainMark("domain")

	route := router.Group(domainMark + "/auth")
	{
		/*Auth
		GET
		 http://www.abc.com/auth/abcedfg/xxx?xx
		*/
		route.GET("/:token/*action", func(ctx *gin.Context) {

			authController.NewController(ctx).GetAuthProtocol()
		})

		route.GET("/test-auth", func(ctx *gin.Context) {
			authController.NewController(ctx).TestAuth()
		})
	}
}
