package helloController

import (
	baseController "develop-template/app/http/controller"
	"develop-template/app/http/middleware"
	"fmt"
	"github.com/JasonMetal/submodule-support-go.git/helper/logger"
	"github.com/gin-gonic/gin"
)

type controller struct {
	baseController.BaseController
}

func NewController(ctx *gin.Context) *controller {
	return &controller{baseController.NewBaseController(ctx)}
}
func (c *controller) Hello() {
	fmt.Println(c.GCtx.Value(middleware.XRequestIDKey))
	logger.L(c.GCtx).Infof("==== %s\n", "hello world111111111111")

	c.Success("hello world")
}
