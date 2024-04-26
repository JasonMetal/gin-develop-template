package controller

import "github.com/gin-gonic/gin"

type BaseController struct {
	GCtx *gin.Context
	Response
	Request
}

func NewBaseController(ctx *gin.Context) BaseController {
	c := BaseController{}
	c.GCtx = ctx
	c.Request.GCtx = c.GCtx
	c.Response.GCtx = c.GCtx
	return c
}
