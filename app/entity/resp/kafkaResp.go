package resp

import "develop-template/app/entity"

// KafkaMessageResp Kafka消息响应
type KafkaMessageResp struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

// FetchMessagesResp 查询消息响应
type FetchMessagesResp struct {
	Messages   []*KafkaMessageResp `json:"messages"`
	Pagination entity.Pagination   `json:"pagination"`
	Topic      string              `json:"topic"`
	Partition  int32               `json:"partition"`
}

// TopicInfoResp 主题信息响应
type TopicInfoResp struct {
	Topic      string          `json:"topic"`
	Partitions []PartitionInfo `json:"partitions"`
}

// PartitionInfo 分区信息
type PartitionInfo struct {
	Partition    int32 `json:"partition"`
	OldestOffset int64 `json:"oldest_offset"`
	NewestOffset int64 `json:"newest_offset"`
	MessageCount int64 `json:"message_count"`
}
