package req

// RabbitMQSendMessageReq 发送消息请求
type RabbitMQSendMessageReq struct {
	QueueName string `json:"queue_name" binding:"required" example:"test-queue"`
	Message   string `json:"message" binding:"required" example:"Hello RabbitMQ"`
}

// RabbitMQSendJSONReq 发送JSON消息请求
type RabbitMQSendJSONReq struct {
	QueueName string                 `json:"queue_name" binding:"required" example:"test-json-queue"`
	Data      map[string]interface{} `json:"data" binding:"required"`
}

// SendToExchangeReq 发送消息到交换机请求
type SendToExchangeReq struct {
	ExchangeName string `json:"exchange_name" binding:"required" example:"test-exchange"`
	ExchangeType string `json:"exchange_type" binding:"required" example:"direct"`
	RoutingKey   string `json:"routing_key" example:"test.routing.key"`
	Message      string `json:"message" binding:"required" example:"Hello Exchange"`
}

// SendFanoutReq 发送广播消息请求
type SendFanoutReq struct {
	ExchangeName string `json:"exchange_name" binding:"required" example:"fanout-exchange"`
	Message      string `json:"message" binding:"required" example:"Broadcast Message"`
}

// SendDirectReq 发送直接消息请求
type SendDirectReq struct {
	ExchangeName string `json:"exchange_name" binding:"required" example:"direct-exchange"`
	RoutingKey   string `json:"routing_key" binding:"required" example:"error"`
	Message      string `json:"message" binding:"required" example:"Error Log"`
}

// SendTopicReq 发送主题消息请求
type SendTopicReq struct {
	ExchangeName string `json:"exchange_name" binding:"required" example:"topic-exchange"`
	RoutingKey   string `json:"routing_key" binding:"required" example:"user.created"`
	Message      string `json:"message" binding:"required" example:"User Created Event"`
}

// SendTaskReq 发送任务消息请求
type SendTaskReq struct {
	QueueName string                 `json:"queue_name" binding:"required" example:"task-queue"`
	TaskName  string                 `json:"task_name" binding:"required" example:"send_email"`
	TaskData  map[string]interface{} `json:"task_data" binding:"required"`
}

// SendBatchReq 批量发送消息请求
type SendBatchReq struct {
	QueueName string   `json:"queue_name" binding:"required" example:"batch-queue"`
	Messages  []string `json:"messages" binding:"required,min=1"`
}

// ============= 查询相关请求 =============

// GetQueueInfoReq 获取队列信息请求
type GetQueueInfoReq struct {
	Queue string `form:"queue" json:"queue" binding:"required" example:"test-queue"` // 队列名称
}

// PeekMessagesReq 查看消息请求（不消费）
type PeekMessagesReq struct {
	Queue string `form:"queue" json:"queue" binding:"required" example:"test-queue"` // 队列名称
	Limit int    `form:"limit" json:"limit" example:"10"`                            // 查看数量，默认10，最多100
}

// ConsumeMessagesReq 消费消息请求（消费并删除）
type ConsumeMessagesReq struct {
	Queue string `form:"queue" json:"queue" binding:"required" example:"test-queue"` // 队列名称
	Limit int    `form:"limit" json:"limit" example:"10"`                            // 消费数量，默认10，最多100
}

// PurgeQueueReq 清空队列请求
type PurgeQueueReq struct {
	Queue string `json:"queue" binding:"required" example:"test-queue"` // 队列名称
}

// DeleteQueueReq 删除队列请求
type DeleteQueueReq struct {
	Queue    string `json:"queue" binding:"required" example:"test-queue"` // 队列名称
	IfUnused bool   `json:"if_unused" example:"false"`                     // 仅在没有消费者时删除
	IfEmpty  bool   `json:"if_empty" example:"false"`                      // 仅在队列为空时删除
}
