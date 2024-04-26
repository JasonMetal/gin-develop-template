package rsaGenRouter

import (
	"develop-template/app/http/controller/api/rsaGenController"
	config2 "develop-template/helper/config"
	"github.com/gin-gonic/gin"
)

func RegisterRsaGen(router *gin.Engine) {
	domainMark := config2.GetDomainMark("domain")

	route := router.Group(domainMark + "/rsa")
	{
		route.GET("/gen", func(ctx *gin.Context) {
			rsaGenController.NewController(ctx).ZipRSAKey()
		})

		route.GET("/download", func(ctx *gin.Context) {
			rsaGenController.NewController(ctx).DownloadCerts()
		})
	}

}
