package authRespEntity

type RespCheck struct {
	Code    uint32        `json:"code"`
	Data    RespCheckData `json:"data"`
	Message string        `json:"message"`
}
type RespCheckData struct {
	Id          uint32 `json:"id"`
	Description string `json:"description"`
	Result      bool   `json:"result"`
	NowTime     int    `json:"now_time"`
}

type RespCommon struct {
	Msg       string `json:"msg"`
	Status    uint32 `json:"status"`
	Success   bool   `json:"success"`
	Time      string `json:"time"`
	Total     uint32 `json:"total"`
	TotalPage uint32 `json:"totalPage"`
}
