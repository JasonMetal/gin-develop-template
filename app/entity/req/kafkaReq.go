package req

// SendMessageReq 发送消息请求
type SendMessageReq struct {
	Topic   string `json:"topic" binding:"required" example:"test-topic"`
	Message string `json:"message" binding:"required" example:"Hello Kafka"`
}

// SendJSONReq 发送JSON消息请求
type SendJSONReq struct {
	Topic string                 `json:"topic" binding:"required" example:"test-json"`
	Data  map[string]interface{} `json:"data" binding:"required"`
}

// FetchMessagesReq 查询消息请求
type FetchMessagesReq struct {
	Topic     string `form:"topic" binding:"required" example:"test-topic"`
	Partition int32  `form:"partition" example:"0"`
	Offset    int64  `form:"offset" example:"0"`
	Limit     int    `form:"limit" example:"10"`
	Page      int    `form:"page" example:"1"`
}

// GetTopicInfoReq 获取主题信息请求
type GetTopicInfoReq struct {
	Topic string `form:"topic" binding:"required" example:"test-topic"`
}
