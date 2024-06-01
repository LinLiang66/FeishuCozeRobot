package handlers

import (
	"FeishuCozeRobot/model"
	"github.com/gin-gonic/gin"
	_ "strconv"
)

func ExecuteHandler(c *gin.Context, event *model.MessageEvent) {
	err := Handler(c, event)
	if err != nil {
		return
	}
}
