package routers

import (
	"FeishuCozeRobot/controller"
	"github.com/gin-gonic/gin"
)

// RegisterRouter 路由设置
func RegisterRouter(router *gin.Engine) {
	routerUser(router)
}

// 用户路由
func routerUser(engine *gin.Engine) {
	con := &controller.WebhookController{}

	// 添加新的路由
	webhookGroup := engine.Group("/webhook")
	{
		webhookGroup.POST("/event/:appid", con.EventHandlerFunc)     //飞书机器人消息事件处理
		webhookGroup.POST("/card/:appid", con.CardActionHandlerFunc) //飞书机器人卡片事件回调处理
	}

}
