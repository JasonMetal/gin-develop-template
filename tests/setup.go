package tests

import (
	"context"
	"develop-template/app/model/common"
	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"testing"
)

var DB *gorm.DB

func Setup() {
	testing.Init()
	bootstrap.DevEnv = bootstrap.EnvLocal
	bootstrap.Init()
	bootstrap.InitWeb([]gin.HandlerFunc{})
	DB = common.ConnectionObject(context.TODO()).DB
}
