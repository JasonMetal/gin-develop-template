package logic

import (
	"math"
)

// getLastPage 计算最后一页分页数
func GetLastPage(total int64, limit int) int {
	lastPage := math.Ceil(float64(total) / float64(limit))
	if lastPage <= 0 || limit == 0 {
		lastPage = 1
	}

	return int(lastPage)
}

// InitCodion 初始化分页
func InitCodion(initPage, initLimit int) (page, limit, offset int) {
	page = initPage
	if initPage == 0 {
		page = 1
	}
	limit = initLimit
	if limit == 0 {
		limit = 10
	}
	offset = (page - 1) * limit

	return
}
