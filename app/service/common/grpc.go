package common

import (
	"context"
	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	"github.com/gin-gonic/gin"
	"time"
)

// IGrpc 定义GRPC接口
type IGrpc interface {
	SetGCtx(ctx *gin.Context)
	SetName(name string)
	GetGCtx() *gin.Context
	GetName() string
}

type BasicGrpc struct {
	GCtx *gin.Context
	Name string
}

func (bg *BasicGrpc) SetGCtx(ctx *gin.Context) {
	bg.GCtx = ctx
}

func (bg *BasicGrpc) SetName(name string) {
	bg.Name = name
}

func (bg *BasicGrpc) GetGCtx() *gin.Context {
	return bg.GCtx
}

func (bg *BasicGrpc) GetName() string {
	return bg.Name
}

func (bg *BasicGrpc) GetConn() (*bootstrap.IdleConn, context.Context) {
	conn, newCtx := bootstrap.GetGrpcConn(bg.GCtx, bg.Name)

	return conn, newCtx
}

// PutConn 接池释放，放回连接池
func (bg *BasicGrpc) PutConn(conn *bootstrap.IdleConn, name string) {
	bootstrap.PutGrpcConn(name, conn)
}

func (bg *BasicGrpc) Call(handle func(*bootstrap.IdleConn, context.Context) error) {
	// 获取连接池
	conn, newCtx := bg.GetConn()
	// 连接释放
	defer bg.PutConn(conn, bg.Name)
	//设置超时时间
	ctx, cancel := context.WithTimeout(newCtx, 5*time.Second)
	defer cancel()
	// 执行rpc调用
	err := handle(conn, ctx)
	// rpc错误上报
	bg.CheckGrpcError(err)
}

// CheckGrpcError TODO error wrap
func (bg *BasicGrpc) CheckGrpcError(err error) {

}
