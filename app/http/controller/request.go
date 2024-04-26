package controller

import (
	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ------------请求信息基础支撑----------

type Request struct {
	GCtx *gin.Context
	Val  string
}

func NewRequest(ctx *gin.Context) Request {
	return Request{GCtx: ctx}
}

// GetHeader 获取响应头
func (req Request) GetHeader(key string) Request {
	req.Val = req.GCtx.GetHeader(key)
	return req
}

func (req Request) GetQuery(key string) Request {
	req.Val = req.GCtx.Query(key)
	return req
}

func (req Request) GetQueryDefault(key string, defaultValue string) Request {
	req.Val = req.GCtx.DefaultQuery(key, defaultValue)
	return req
}

func (req Request) PostForm(key string) Request {
	req.Val = req.GCtx.PostForm(key)
	return req
}

func (req Request) Value() string {
	return req.Val
}

func (req Request) Bool() bool {
	b, _ := strconv.ParseBool(req.Val)
	return b
}

func (req Request) ShouldBindQuery(obj any) error {
	return req.GCtx.ShouldBindQuery(obj)
}

// ShouldBindJSON 使用注意，只能获取一次，无法重复获取
// ShouldBindWith for better performance if you need to call only once.
func (req Request) ShouldBindJSON(obj any) {
	_ = req.GCtx.ShouldBindJSON(obj)
}

func (req Request) IsMimeJson() bool {
	return req.GCtx.ContentType() == binding.MIMEJSON
}

func (req Request) PostToModel(obj any) Request {
	if req.IsMimeJson() {
		_ = req.GCtx.ShouldBindBodyWith(obj, binding.JSON)
	} else {
		_ = req.GCtx.ShouldBind(obj)
	}
	return req
}

// GetAllParamsFromUrl 获取url中所有参数, 不支持数组
func (req Request) GetAllParamsFromUrl() map[string]string {
	reqParams := make(map[string]string)
	params := req.GCtx.Request.URL.Query()
	for key, value := range params {
		if len(value) == 1 {
			reqParams[key] = value[0]
		}
	}

	return reqParams
}

// GetValidMsg 参数校验
func (req Request) GetValidMsg(err error, obj interface{}) string {
	//判断环境，正式环境不返回具体的参数错误
	if bootstrap.DevEnv == bootstrap.EnvProduct {
		return "无效参数"
	}
	getObj := reflect.TypeOf(obj)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			if f, exist := getObj.Elem().FieldByName(e.Field()); exist {
				filedName := f.Tag.Get("json")
				if filedName == "" {
					filedName = f.Tag.Get("form")
				}
				return filedName + f.Tag.Get("msg")
			}
		}
	}
	//传参数据类型错误，长度超出的，和未定义的错误都返回无效参数
	if err != nil {
		return "无效参数"
	}

	return err.Error()
}
