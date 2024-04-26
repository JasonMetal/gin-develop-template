/**
* @Author: 渣渣辉
* @Description:
* @File: reqEntity
* @Module: newsPage
* @Date: 2022-11-04 09:41:03
 */
package testReqEntity

import (
	"develop-template/app/entity"
)

// List 列表
type List struct {
	entity.PaginationSearch
}

// Item 单个操作
type Item struct {
	Id uint32 `json:"news_page_id" form:"news_page_id"  binding:"required" msg:"必填"`
}

type TestData struct {
	TestId    int32  `gorm:"column:test_id"                 json:"test_id"`
	TestName  string `gorm:"column:test_name"              json:"test_name"`
	TestUnion string `gorm:"column:test_union"            json:"test_union"`
	//TestTag    int32  `gorm:"column:test_tag"             json:"test_tag"`
	TestStatus int32 `gorm:"column:test_status"              json:"test_status"`
}

// Update 更新
type Update struct {
	Id uint32 `json:"news_page_id"  binding:"required" msg:"必填"`
}
