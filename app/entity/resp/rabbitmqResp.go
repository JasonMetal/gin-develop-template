package resp

// RabbitMQMessageResp RabbitMQ消息响应
type RabbitMQMessageResp struct {
	QueueName    string `json:"queue_name"`
	Message      string `json:"message"`
	SentAt       string `json:"sent_at"`
	MessageCount int    `json:"message_count,omitempty"`
}

// RabbitMQExchangeResp 交换机消息响应
type RabbitMQExchangeResp struct {
	ExchangeName string `json:"exchange_name"`
	ExchangeType string `json:"exchange_type"`
	RoutingKey   string `json:"routing_key"`
	Message      string `json:"message"`
	SentAt       string `json:"sent_at"`
}

// RabbitMQTaskResp 任务消息响应
type RabbitMQTaskResp struct {
	QueueName string `json:"queue_name"`
	TaskName  string `json:"task_name"`
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// RabbitMQBatchResp 批量发送响应
type RabbitMQBatchResp struct {
	QueueName    string `json:"queue_name"`
	TotalCount   int    `json:"total_count"`
	SuccessCount int    `json:"success_count"`
	FailedCount  int    `json:"failed_count"`
	SentAt       string `json:"sent_at"`
}

// RabbitMQHealthResp 健康检查响应
type RabbitMQHealthResp struct {
	Status    string `json:"status"`
	Connected bool   `json:"connected"`
	Message   string `json:"message,omitempty"`
}

// ============= 查询相关响应 =============

// QueueInfoResp 队列信息响应
type QueueInfoResp struct {
	Name       string `json:"name"`        // 队列名称
	Messages   int    `json:"messages"`    // 消息数量
	Consumers  int    `json:"consumers"`   // 消费者数量
	Durable    bool   `json:"durable"`     // 是否持久化
	AutoDelete bool   `json:"auto_delete"` // 是否自动删除
	Exclusive  bool   `json:"exclusive"`   // 是否排他
}

// QueueMessageResp 队列消息响应
type QueueMessageResp struct {
	Body          string            `json:"body"`           // 消息体
	ContentType   string            `json:"content_type"`   // 内容类型
	DeliveryMode  uint8             `json:"delivery_mode"`  // 投递模式：1=非持久化, 2=持久化
	Priority      uint8             `json:"priority"`       // 优先级
	CorrelationId string            `json:"correlation_id"` // 关联ID
	ReplyTo       string            `json:"reply_to"`       // 回复队列
	Expiration    string            `json:"expiration"`     // 过期时间
	MessageId     string            `json:"message_id"`     // 消息ID
	Timestamp     string            `json:"timestamp"`      // 时间戳
	Type          string            `json:"type"`           // 消息类型
	UserId        string            `json:"user_id"`        // 用户ID
	AppId         string            `json:"app_id"`         // 应用ID
	Headers       map[string]string `json:"headers"`        // 消息头
	DeliveryTag   uint64            `json:"delivery_tag"`   // 投递标签
	Redelivered   bool              `json:"redelivered"`    // 是否重新投递
	Exchange      string            `json:"exchange"`       // 交换机
	RoutingKey    string            `json:"routing_key"`    // 路由键
}

// PeekMessagesResp 查看消息响应
type PeekMessagesResp struct {
	Queue    string             `json:"queue"`    // 队列名称
	Total    int                `json:"total"`    // 返回的消息数量
	Messages []QueueMessageResp `json:"messages"` // 消息列表
}

// ConsumeMessagesResp 消费消息响应
type ConsumeMessagesResp struct {
	Queue    string             `json:"queue"`    // 队列名称
	Total    int                `json:"total"`    // 消费的消息数量
	Messages []QueueMessageResp `json:"messages"` // 消息列表
}

// PurgeQueueResp 清空队列响应
type PurgeQueueResp struct {
	Queue        string `json:"queue"`         // 队列名称
	DeletedCount int    `json:"deleted_count"` // 删除的消息数量
	Success      bool   `json:"success"`       // 是否成功
}

// DeleteQueueResp 删除队列响应
type DeleteQueueResp struct {
	Queue        string `json:"queue"`         // 队列名称
	DeletedCount int    `json:"deleted_count"` // 队列中删除的消息数量
	Success      bool   `json:"success"`       // 是否成功
}
