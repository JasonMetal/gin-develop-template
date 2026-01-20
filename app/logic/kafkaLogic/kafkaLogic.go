package kafkaLogic

import (
	"context"
	"develop-template/app/service/kafkaService"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	configLib "github.com/olebedev/config"
)

type logic struct {
	Ctx        context.Context
	GCtx       *gin.Context
	userId     uint32
	CurlConfig *configLib.Config
}

func NewLogic(ctx context.Context) *logic {
	//func NewLogic(ctx *gin.Context) *logic {
	return &logic{Ctx: ctx}
}

func (l *logic) TestSendMsg() {
	service := kafkaService.NewKafkaService()
	if service == nil {
		log.Printf("创建Kafka服务失败")
		return
	}

	testData := map[string]interface{}{
		"user_id": 99999,
		"action":  "login",
		"time":    time.Now().Unix(),
	}

	// 方式1: 使用SendJSON直接发送map数据（推荐）
	// SendJSON内部会自动进行JSON序列化
	err := service.SendJSON("test-json", testData)
	if err != nil {
		log.Printf("SendJSON发送失败: %v", err)
	} else {
		log.Printf("SendJSON发送成功")
	}

	// 方式2: 如果需要手动序列化，应该使用SendMessage
	jsonData, err := json.Marshal(testData)
	if err != nil {
		log.Printf("JSON序列化失败: %v", err)
		return
	}
	err = service.SendMessage("test-json-manual", string(jsonData))
	if err != nil {
		log.Printf("SendMessage发送失败: %v", err)
	} else {
		log.Printf("SendMessage发送成功")
	}

	// 方式3: 发送普通文本消息
	err = service.SendMessage("test", "111111111111111111111111")
	if err != nil {
		log.Printf("发送文本消息失败: %v", err)
	} else {
		log.Printf("发送文本消息成功")
	}
}
