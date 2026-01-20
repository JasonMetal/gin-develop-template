package rabbitmqController

import (
	"develop-template/app/entity/req"
	"develop-template/app/entity/resp"
	baseController "develop-template/app/http/controller"
	"develop-template/app/service/rabbitmqService"
	"fmt"
	"time"

	"github.com/JasonMetal/submodule-support-go.git/helper/logger"
	"github.com/gin-gonic/gin"
)

type controller struct {
	baseController.BaseController
}

func NewController(ctx *gin.Context) *controller {
	return &controller{baseController.NewBaseController(ctx)}
}

// SendMessage 发送消息
// @Summary 发送RabbitMQ消息
// @Description 发送单条消息到指定队列
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendMessageReq true "发送消息请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send [post]
func (c *controller) SendMessage() {
	var request req.RabbitMQSendMessageReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.SendMessage(request.QueueName, request.Message); err != nil {
		logger.L(c.GCtx).Infof("发送消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQMessageResp{
		QueueName: request.QueueName,
		Message:   request.Message,
		SentAt:    time.Now().Format(time.RFC3339),
	}, "消息发送成功")
}

// SendJSON 发送JSON消息
// @Summary 发送JSON格式的RabbitMQ消息
// @Description 发送JSON格式数据到指定队列
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendJSONReq true "发送JSON消息请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send-json [post]
func (c *controller) SendJSON() {
	var request req.RabbitMQSendJSONReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.SendJSON(request.QueueName, request.Data); err != nil {
		logger.L(c.GCtx).Infof("发送JSON消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送JSON消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQMessageResp{
		QueueName: request.QueueName,
		Message:   fmt.Sprintf("%v", request.Data),
		SentAt:    time.Now().Format(time.RFC3339),
	}, "JSON消息发送成功")
}

// SendToExchange 发送消息到交换机
// @Summary 发送消息到交换机
// @Description 发送消息到指定交换机
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendToExchangeReq true "发送到交换机请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send-exchange [post]
func (c *controller) SendToExchange() {
	var request req.SendToExchangeReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.SendToExchange(request.ExchangeName, request.ExchangeType, request.RoutingKey, request.Message); err != nil {
		logger.L(c.GCtx).Infof("发送消息到交换机失败: %v", err)
		c.Fail400(fmt.Sprintf("发送消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQExchangeResp{
		ExchangeName: request.ExchangeName,
		ExchangeType: request.ExchangeType,
		RoutingKey:   request.RoutingKey,
		Message:      request.Message,
		SentAt:       time.Now().Format(time.RFC3339),
	}, "消息发送成功")
}

// SendFanout 发送广播消息
// @Summary 发送广播消息(Fanout模式)
// @Description 发送广播消息到Fanout交换机
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendFanoutReq true "发送广播消息请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send-fanout [post]
func (c *controller) SendFanout() {
	var request req.SendFanoutReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.PublishFanout(request.ExchangeName, request.Message); err != nil {
		logger.L(c.GCtx).Infof("发送广播消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQExchangeResp{
		ExchangeName: request.ExchangeName,
		ExchangeType: "fanout",
		Message:      request.Message,
		SentAt:       time.Now().Format(time.RFC3339),
	}, "广播消息发送成功")
}

// SendDirect 发送直接消息
// @Summary 发送直接消息(Direct模式)
// @Description 发送直接消息到Direct交换机
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendDirectReq true "发送直接消息请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send-direct [post]
func (c *controller) SendDirect() {
	var request req.SendDirectReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.PublishDirect(request.ExchangeName, request.RoutingKey, request.Message); err != nil {
		logger.L(c.GCtx).Infof("发送直接消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQExchangeResp{
		ExchangeName: request.ExchangeName,
		ExchangeType: "direct",
		RoutingKey:   request.RoutingKey,
		Message:      request.Message,
		SentAt:       time.Now().Format(time.RFC3339),
	}, "直接消息发送成功")
}

// SendTopic 发送主题消息
// @Summary 发送主题消息(Topic模式)
// @Description 发送主题消息到Topic交换机
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendTopicReq true "发送主题消息请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send-topic [post]
func (c *controller) SendTopic() {
	var request req.SendTopicReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.PublishTopic(request.ExchangeName, request.RoutingKey, request.Message); err != nil {
		logger.L(c.GCtx).Infof("发送主题消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送消息失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQExchangeResp{
		ExchangeName: request.ExchangeName,
		ExchangeType: "topic",
		RoutingKey:   request.RoutingKey,
		Message:      request.Message,
		SentAt:       time.Now().Format(time.RFC3339),
	}, "主题消息发送成功")
}

// SendTask 发送任务消息
// @Summary 发送任务消息
// @Description 发送任务消息到队列(Worker模式)
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendTaskReq true "发送任务消息请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send-task [post]
func (c *controller) SendTask() {
	var request req.SendTaskReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.SendTask(request.QueueName, request.TaskName, request.TaskData); err != nil {
		logger.L(c.GCtx).Infof("发送任务消息失败: %v", err)
		c.Fail400(fmt.Sprintf("发送任务失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQTaskResp{
		QueueName: request.QueueName,
		TaskName:  request.TaskName,
		Status:    "sent",
		CreatedAt: time.Now().Format(time.RFC3339),
	}, "任务发送成功")
}

// SendBatch 批量发送消息
// @Summary 批量发送消息
// @Description 批量发送消息到队列
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Param request body req.SendBatchReq true "批量发送消息请求"
// @Success 200 {object} controller.ResJson
// @Router /rabbitmq/send-batch [post]
func (c *controller) SendBatch() {
	var request req.SendBatchReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	if err := service.SendBatch(request.QueueName, request.Messages); err != nil {
		logger.L(c.GCtx).Infof("批量发送消息失败: %v", err)
		c.Fail400(fmt.Sprintf("批量发送失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.RabbitMQBatchResp{
		QueueName:    request.QueueName,
		TotalCount:   len(request.Messages),
		SuccessCount: len(request.Messages),
		FailedCount:  0,
		SentAt:       time.Now().Format(time.RFC3339),
	}, "批量发送成功")
}

// HealthCheck 健康检查
// @Summary RabbitMQ健康检查
// @Description 检查RabbitMQ连接状态
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Success 200 {object} controller.ResJson{data=resp.RabbitMQHealthResp}
// @Router /rabbitmq/health [get]
func (c *controller) HealthCheck() {
	service := rabbitmqService.NewRabbitMQService()
	if err := service.ValidateConnection(); err != nil {
		c.Success(resp.RabbitMQHealthResp{
			Status:    "error",
			Connected: false,
			Message:   err.Error(),
		})
		return
	}

	c.Success(resp.RabbitMQHealthResp{
		Status:    "ok",
		Connected: true,
		Message:   "RabbitMQ连接正常",
	})
}

// ============= 查询相关接口 =============

// GetQueueInfo 获取队列信息
// @Summary 获取队列信息
// @Description 获取指定队列的详细信息（消息数、消费者数等）
// @Tags RabbitMQ-Query
// @Accept json
// @Produce json
// @Param queue query string true "队列名称"
// @Success 200 {object} controller.ResJson{data=resp.QueueInfoResp}
// @Router /rabbitmq/queue/info [get]
func (c *controller) GetQueueInfo() {
	var request req.GetQueueInfoReq
	if err := c.GCtx.ShouldBindQuery(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	info, err := service.DeclareAndGetQueueInfo(request.Queue)
	if err != nil {
		logger.L(c.GCtx).Infof("获取队列信息失败: %v", err)
		c.Fail400(fmt.Sprintf("获取队列信息失败: %v", err))
		return
	}

	c.Success(resp.QueueInfoResp{
		Name:       info.Name,
		Messages:   info.Messages,
		Consumers:  info.Consumers,
		Durable:    info.Durable,
		AutoDelete: info.AutoDelete,
		Exclusive:  info.Exclusive,
	})
}

// PeekMessages 查看队列消息（不消费）
// @Summary 查看队列消息
// @Description 查看队列中的消息，但不从队列中删除（消息会重新入队）
// @Tags RabbitMQ-Query
// @Accept json
// @Produce json
// @Param queue query string true "队列名称"
// @Param limit query int false "查看数量(默认10,最多100)"
// @Success 200 {object} controller.ResJson{data=resp.PeekMessagesResp}
// @Router /rabbitmq/queue/peek [get]
func (c *controller) PeekMessages() {
	var request req.PeekMessagesReq
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

	service := rabbitmqService.NewRabbitMQService()
	messages, err := service.PeekMessages(request.Queue, request.Limit)
	if err != nil {
		logger.L(c.GCtx).Infof("查看消息失败: %v", err)
		c.Fail400(fmt.Sprintf("查看消息失败: %v", err))
		return
	}

	// 转换消息格式
	var respMessages []resp.QueueMessageResp
	for _, msg := range messages {
		respMessages = append(respMessages, resp.QueueMessageResp{
			Body:          msg.Body,
			ContentType:   msg.ContentType,
			DeliveryMode:  msg.DeliveryMode,
			Priority:      msg.Priority,
			CorrelationId: msg.CorrelationId,
			ReplyTo:       msg.ReplyTo,
			Expiration:    msg.Expiration,
			MessageId:     msg.MessageId,
			Timestamp:     msg.Timestamp.Format(time.RFC3339),
			Type:          msg.Type,
			UserId:        msg.UserId,
			AppId:         msg.AppId,
			Headers:       msg.Headers,
			DeliveryTag:   msg.DeliveryTag,
			Redelivered:   msg.Redelivered,
			Exchange:      msg.Exchange,
			RoutingKey:    msg.RoutingKey,
		})
	}

	c.Success(resp.PeekMessagesResp{
		Queue:    request.Queue,
		Total:    len(respMessages),
		Messages: respMessages,
	})
}

// ConsumeMessages 消费队列消息
// @Summary 消费队列消息
// @Description 从队列中消费消息（消息会被删除）
// @Tags RabbitMQ-Query
// @Accept json
// @Produce json
// @Param queue query string true "队列名称"
// @Param limit query int false "消费数量(默认10,最多100)"
// @Success 200 {object} controller.ResJson{data=resp.ConsumeMessagesResp}
// @Router /rabbitmq/queue/consume [get]
func (c *controller) ConsumeMessages() {
	var request req.ConsumeMessagesReq
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

	service := rabbitmqService.NewRabbitMQService()
	messages, err := service.ConsumeMessages(request.Queue, request.Limit)
	if err != nil {
		logger.L(c.GCtx).Infof("消费消息失败: %v", err)
		c.Fail400(fmt.Sprintf("消费消息失败: %v", err))
		return
	}

	// 转换消息格式
	var respMessages []resp.QueueMessageResp
	for _, msg := range messages {
		respMessages = append(respMessages, resp.QueueMessageResp{
			Body:          msg.Body,
			ContentType:   msg.ContentType,
			DeliveryMode:  msg.DeliveryMode,
			Priority:      msg.Priority,
			CorrelationId: msg.CorrelationId,
			ReplyTo:       msg.ReplyTo,
			Expiration:    msg.Expiration,
			MessageId:     msg.MessageId,
			Timestamp:     msg.Timestamp.Format(time.RFC3339),
			Type:          msg.Type,
			UserId:        msg.UserId,
			AppId:         msg.AppId,
			Headers:       msg.Headers,
			DeliveryTag:   msg.DeliveryTag,
			Redelivered:   msg.Redelivered,
			Exchange:      msg.Exchange,
			RoutingKey:    msg.RoutingKey,
		})
	}

	c.Success(resp.ConsumeMessagesResp{
		Queue:    request.Queue,
		Total:    len(respMessages),
		Messages: respMessages,
	})
}

// PurgeQueue 清空队列
// @Summary 清空队列
// @Description 删除队列中的所有消息
// @Tags RabbitMQ-Query
// @Accept json
// @Produce json
// @Param request body req.PurgeQueueReq true "清空队列请求"
// @Success 200 {object} controller.ResJson{data=resp.PurgeQueueResp}
// @Router /rabbitmq/queue/purge [post]
func (c *controller) PurgeQueue() {
	var request req.PurgeQueueReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	count, err := service.PurgeQueue(request.Queue)
	if err != nil {
		logger.L(c.GCtx).Infof("清空队列失败: %v", err)
		c.Fail400(fmt.Sprintf("清空队列失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.PurgeQueueResp{
		Queue:        request.Queue,
		DeletedCount: count,
		Success:      true,
	}, fmt.Sprintf("成功清空队列，删除了 %d 条消息", count))
}

// DeleteQueue 删除队列
// @Summary 删除队列
// @Description 删除指定的队列
// @Tags RabbitMQ-Query
// @Accept json
// @Produce json
// @Param request body req.DeleteQueueReq true "删除队列请求"
// @Success 200 {object} controller.ResJson{data=resp.DeleteQueueResp}
// @Router /rabbitmq/queue/delete [post]
func (c *controller) DeleteQueue() {
	var request req.DeleteQueueReq
	if err := c.GCtx.ShouldBindJSON(&request); err != nil {
		logger.L(c.GCtx).Infof("参数绑定失败: %v", err)
		c.Fail400("参数错误: " + err.Error())
		return
	}

	service := rabbitmqService.NewRabbitMQService()
	count, err := service.DeleteQueue(request.Queue, request.IfUnused, request.IfEmpty)
	if err != nil {
		logger.L(c.GCtx).Infof("删除队列失败: %v", err)
		c.Fail400(fmt.Sprintf("删除队列失败: %v", err))
		return
	}

	c.SuccessWithMsg(resp.DeleteQueueResp{
		Queue:        request.Queue,
		DeletedCount: count,
		Success:      true,
	}, fmt.Sprintf("成功删除队列，删除了 %d 条消息", count))
}
