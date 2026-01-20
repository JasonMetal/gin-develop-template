package kafkaController

import (
	"develop-template/app/entity"
	"develop-template/app/entity/req"
	"develop-template/app/entity/resp"
	baseController "develop-template/app/http/controller"
	"develop-template/app/service/kafkaService"
	"fmt"
	"github.com/JasonMetal/submodule-support-go.git/helper/logger"
	"github.com/gin-gonic/gin"
	"math"
)

type controller struct {
	baseController.BaseController
}

func NewController(ctx *gin.Context) *controller {
	return &controller{baseController.NewBaseController(ctx)}
}

// SendMessage 发送消息
// @Summary 发送Kafka消息
// @Description 发送单条消息到指定主题
// @Tags Kafka
// @Accept json
// @Produce json
// @Param request body req.SendMessageReq true "发送消息请求"
// @Success 200 {object} controller.ResJson
// @Router /kafka/send [post]
func (c *controller) SendMessage() {
	var request req.SendMessageReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := kafkaService.NewKafkaService()
	if err := service.SendMessage(request.Topic, request.Message); err != nil {
		logger.L(c.GCtx).Infof("发送消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(gin.H{
		"topic":   request.Topic,
		"message": request.Message,
	}, "消息发送成功")
}

// SendJSON 发送JSON消息
// @Summary 发送JSON格式的Kafka消息
// @Description 发送JSON格式数据到指定主题
// @Tags Kafka
// @Accept json
// @Produce json
// @Param request body req.SendJSONReq true "发送JSON消息请求"
// @Success 200 {object} controller.ResJson
// @Router /kafka/send-json [post]
func (c *controller) SendJSON() {
	var request req.SendJSONReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := kafkaService.NewKafkaService()
	if err := service.SendJSON(request.Topic, request.Data); err != nil {
		logger.L(c.GCtx).Infof("发送JSON消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送JSON消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(gin.H{
		"topic": request.Topic,
		"data":  request.Data,
	}, "JSON消息发送成功")
}

// FetchMessages 查询消息
// @Summary 查询Kafka消息
// @Description 查询指定主题的消息，支持分页
// @Tags Kafka
// @Accept json
// @Produce json
// @Param topic query string true "主题名称"
// @Param partition query int false "分区编号" default(0)
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} controller.ResJson{data=resp.FetchMessagesResp}
// @Router /kafka/messages [get]
func (c *controller) FetchMessages() {
	var request req.FetchMessagesReq
	if err := c.GCtx.ShouldBindQuery(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	// 设置默认值
	if request.Limit <= 0 {
		request.Limit = 10
	}
	if request.Limit > 100 {
		request.Limit = 100
	}
	if request.Page <= 0 {
		request.Page = 1
	}

	service := kafkaService.NewKafkaService()

	// 如果没有指定offset，使用page计算
	if request.Offset == 0 && request.Page > 1 {
		// 获取分区的最早偏移量
		oldestOffset, _, err := service.GetPartitionOffset(request.Topic, request.Partition)
		if err != nil {
			logger.L(c.GCtx).Infof("获取分区偏移量失败: %v", err)
			c.Fail400(fmt.Sprintf("获取分区偏移量失败: %v", err))
			return
		}
		request.Offset = oldestOffset + int64((request.Page-1)*request.Limit)
	}

	// 获取消息
	messages, err := service.FetchMessages(request.Topic, request.Partition, request.Offset, request.Limit)
	if err != nil {
		logger.L(c.GCtx).Infof("查询消息失败: %v", err)
		c.Fail400(fmt.Sprintf("查询消息失败: %v", err))
		return
	}

	// 获取分区的偏移量信息计算总数
	oldestOffset, newestOffset, err := service.GetPartitionOffset(request.Topic, request.Partition)
	if err != nil {
		logger.L(c.GCtx).Infof("获取偏移量信息失败: %v", err)
		oldestOffset = 0
		newestOffset = 0
	}

	total := newestOffset - oldestOffset
	if total < 0 {
		total = 0
	}

	lastPage := int(math.Ceil(float64(total) / float64(request.Limit)))
	if lastPage < 1 {
		lastPage = 1
	}

	// 转换消息格式
	respMessages := make([]*resp.KafkaMessageResp, 0, len(messages))
	for _, msg := range messages {
		respMessages = append(respMessages, &resp.KafkaMessageResp{
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Key:       msg.Key,
			Value:     msg.Value,
			Timestamp: msg.Timestamp,
		})
	}

	response := resp.FetchMessagesResp{
		Messages:  respMessages,
		Topic:     request.Topic,
		Partition: request.Partition,
		Pagination: entity.Pagination{
			Total:    total,
			Page:     request.Page,
			LastPage: lastPage,
		},
	}

	c.Success(response)
}

// GetTopicInfo 获取主题信息
// @Summary 获取主题信息
// @Description 获取主题的分区信息和偏移量
// @Tags Kafka
// @Accept json
// @Produce json
// @Param topic query string true "主题名称"
// @Success 200 {object} controller.ResJson{data=resp.TopicInfoResp}
// @Router /kafka/topic-info [get]
func (c *controller) GetTopicInfo() {
	var request req.GetTopicInfoReq
	if err := c.GCtx.ShouldBindQuery(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := kafkaService.NewKafkaService()

	// 获取分区列表
	partitions, err := service.GetTopicPartitions(request.Topic)
	if err != nil {
		logger.L(c.GCtx).Infof("获取分区信息失败: %v", err)
		c.Fail400(fmt.Sprintf("获取分区信息失败: %v", err))
		return
	}

	// 获取每个分区的偏移量信息
	partitionInfos := make([]resp.PartitionInfo, 0, len(partitions))
	for _, partition := range partitions {
		oldestOffset, newestOffset, err := service.GetPartitionOffset(request.Topic, partition)
		if err != nil {
			logger.L(c.GCtx).Infof("获取分区 %d 偏移量失败: %v", partition, err)
			continue
		}

		messageCount := newestOffset - oldestOffset
		if messageCount < 0 {
			messageCount = 0
		}

		partitionInfos = append(partitionInfos, resp.PartitionInfo{
			Partition:    partition,
			OldestOffset: oldestOffset,
			NewestOffset: newestOffset,
			MessageCount: messageCount,
		})
	}

	response := resp.TopicInfoResp{
		Topic:      request.Topic,
		Partitions: partitionInfos,
	}

	c.Success(response)
}
