package testRespEntity

import (
	"develop-template/app/entity"
	"develop-template/app/entity/db/testDbEntity"
)

// List 列表
type List struct {
	entity.Pagination
	List []*testDbEntity.TestData `json:"list"`
}

type OnlineList struct {
	//List []testDbEntity.TestData `json:"list"`
	List []*testDbEntity.TestData `json:"list"`
}

// Info 信息
type Info struct {
	testDbEntity.TestData
}
