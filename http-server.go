package main

import (
	"develop-template/constant"
	router "develop-template/routes"
	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	"github.com/gin-gonic/gin"
)

// ProjectName

// @title        业务应用开发API
// @description  提供业务应用开发的业务功能APIs
// @schemes      http https
func main() {
	bootstrap.SetProjectName(constant.ProjectName)
	// 初始化Web
	bootstrap.Init()
	middleFun := []gin.HandlerFunc{
		//	middleware.CheckUserAuth(),
	}
	r := bootstrap.InitWeb(middleFun)
	router.RegisterRouter(r)
	bootstrap.RunWeb(r, constant.HttpServiceHostPort)
}
