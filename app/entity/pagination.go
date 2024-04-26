package entity

type Pagination struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	LastPage int   `json:"last_page"`
}

// PaginationSearch 搜索参数
type PaginationSearch struct {
	Page   int `form:"page"`
	Limit  int `form:"limit"`
	Offset int `form:"Offset"`
}
