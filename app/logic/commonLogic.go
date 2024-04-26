package logic

import (
	"develop-template/app/constant"
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

// InitCondition 初始化分页
func InitCondition(initPage, initLimit int) (page, limit, offset int) {
	if initPage == 0 {
		initPage = constant.InitPage
	}

	if initLimit == 0 {
		initLimit = constant.InitLimit
	}
	return initPage, initLimit, (initPage - 1) * initLimit
}
