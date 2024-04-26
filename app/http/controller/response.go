package controller

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// ------------响应内容基础支撑----------

type Response struct {
	GCtx *gin.Context
}

func NewResponse(ctx *gin.Context) Response {
	return Response{GCtx: ctx}
}

// ResJson 响应的数据结构
type ResJson struct {
	Code    int32       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (res Response) ResJson(resJson ResJson) {
	res.GCtx.JSON(http.StatusOK, resJson)
}

func (res Response) Fail(code int32, errMsg string) {
	res.ResJson(ResJson{Code: code, Message: errMsg, Data: struct{}{}})
}

func (res Response) Fail400(errMsg string) {
	if errMsg == "" {
		res.FailAbnormal()
		return
	}
	res.ResJson(ResJson{Code: 400, Message: errMsg, Data: ""})
}

// FailRefuse 请求拒绝
func (res Response) FailRefuse() {
	res.Fail400("非法请求")
}

// FailAbnormal 服务器异常
func (res Response) FailAbnormal() {
	res.Fail400("服务器繁忙，请稍后再试")
}

func (res Response) Success(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res.ResJson(ResJson{Code: 0, Data: data})
}

func (res Response) SuccessWithMsg(data interface{}, msg string) {
	if data == nil {
		data = gin.H{}
	}

	if msg == "" {
		msg = "操作成功"
	}

	res.ResJson(ResJson{Code: 0, Data: data, Message: msg})
}

// SetHeader 设置响应头
func (res Response) SetHeader(name string, val string) {
	res.GCtx.Header(name, val)
}

// SetSrvHeader 设置golang服务器独有请求头，区分PHP/Golang
func (res Response) SetSrvHeader() {
	if res.GCtx.Writer == nil {
		return
	}
	res.SetHeader("web-srv", "go")
}

// NilToObj nil判断，如果是nil 返回空对象
func (res Response) NilToObj(value interface{}) interface{} {
	rv := reflect.ValueOf(value)
	if rv.IsNil() {
		return map[string]interface{}{}
	}
	return value
}

// NilToArray nil判断，如果是nil 返回空数组
func (res Response) NilToArray(values interface{}) interface{} {
	rv := reflect.ValueOf(values)
	if rv.IsNil() {
		return []interface{}{}
	}
	return values
}
