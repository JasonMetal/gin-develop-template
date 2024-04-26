// @Author 911ja 2023/1/6/006 10:08:00
package configService

import (
	"context"
	siteCache "develop-template/app/cache/site"
	"github.com/gin-gonic/gin"
)

type service struct {
	Ctx    context.Context
	GCtx   *gin.Context
	userId uint32
}

func NewService(ctx context.Context) *service {
	return &service{Ctx: ctx}
}

func (l *service) GetSite() (getUrl string) {
	aCache := siteCache.NewSiteCache(l.Ctx)
	getUrl = aCache.GetSiteCache("Site")
	return getUrl
}
