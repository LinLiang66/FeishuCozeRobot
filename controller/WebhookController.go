package controller

import (
	"FeishuCozeRobot/handlers"
	"FeishuCozeRobot/model"
	"FeishuCozeRobot/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"time"
)

type WebhookController struct {
}

// EventHandlerFunc 接收并处理消息事件
func (controller *WebhookController) EventHandlerFunc(c *gin.Context) {
	appid := c.Param("appid")
	c.Header("Server", "Go-Gin-Server")
	if len(appid) == 0 {
		c.JSON(200, gin.H{
			"message":   "appid Cannot be empty!!",
			"code":      -1,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	}
	appCache, exist := handlers.GetAppCache(c, appid)
	if !exist {
		c.JSON(200, gin.H{
			"message":   "appid is invalid",
			"code":      -1,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
	} else {
		plainEventJsonStr, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(400, "Error reading request body")
			return
		}
		fuzzy := &larkevent.EventFuzzy{}
		err = json.Unmarshal(plainEventJsonStr, &fuzzy)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		if larkevent.ReqType(fuzzy.Type) == larkevent.ReqTypeChallenge && fuzzy.Token == appCache.VerificationToken {
			c.JSON(200, gin.H{
				"message":   "success",
				"challenge": fuzzy.Challenge,
				"code":      200,
				"success":   true,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		EventJsonStr, err := utils.EventDecrypt(fuzzy.Encrypt, appCache.EncryptKey)
		event := &model.MessageEvent{}
		err = json.Unmarshal([]byte(EventJsonStr), &event)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if larkevent.ReqType(event.Type) == larkevent.ReqTypeChallenge {
			c.JSON(200, gin.H{
				"message":   "success",
				"challenge": event.Challenge,
				"code":      200,
				"success":   true,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		c.JSON(200, gin.H{
			"message":   "success",
			"code":      200,
			"success":   true,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		if event.Header.EventType == "im.message.receive_v1" {
			go handlers.ExecuteHandler(c, event)
		}
	}
	return
}

// CardActionHandlerFunc 接收并处理卡片回调事件
func (controller *WebhookController) CardActionHandlerFunc(c *gin.Context) {
	appid := c.Param("appid")
	c.Header("Server", "Go-Gin-Server")
	if len(appid) == 0 {
		c.JSON(200, gin.H{
			"message":   "appid Cannot be empty!!",
			"code":      -1,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
		return
	}
	appCache, exist := handlers.GetAppCache(c, appid)
	if !exist {
		c.JSON(200, gin.H{
			"message":   "appid is invalid",
			"code":      -1,
			"success":   false,
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		})
	} else {
		plainEventJsonStr, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(400, "Error reading request body")
			return
		}
		fuzzy := &larkevent.EventFuzzy{}
		err = json.Unmarshal(plainEventJsonStr, &fuzzy)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		if larkevent.ReqType(fuzzy.Type) == larkevent.ReqTypeChallenge && fuzzy.Token == appCache.VerificationToken {
			c.JSON(200, gin.H{
				"message":   "success",
				"challenge": fuzzy.Challenge,
				"code":      200,
				"success":   true,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		cardAction := &model.CardAction{}
		if len(fuzzy.Encrypt) > 50 {
			EventJsonStr, err := utils.EventDecrypt(fuzzy.Encrypt, appCache.EncryptKey)
			err = json.Unmarshal([]byte(EventJsonStr), &cardAction)
			if err != nil {
				fmt.Println("Error:", err)
				c.JSON(200, gin.H{
					"message":   err.Error(),
					"code":      -1,
					"success":   false,
					"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
				})
				return
			}
		} else {
			err = json.Unmarshal(plainEventJsonStr, &cardAction)
			if err != nil {
				fmt.Println("Error:", err)
				c.JSON(200, gin.H{
					"message":   err.Error(),
					"code":      -1,
					"success":   false,
					"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
				})
				return
			}
		}

		if larkevent.ReqType(cardAction.Type) == larkevent.ReqTypeChallenge && cardAction.Token == appCache.VerificationToken {
			println("验证成功")
			c.JSON(200, gin.H{
				"message":   "success",
				"challenge": cardAction.Challenge,
				"code":      200,
				"success":   true,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		if handlers.RedisClient.KEYEXISTS(c, "robot:card_event:"+cardAction.UserID+":"+cardAction.OpenMessageID) {
			c.JSON(200, gin.H{
				"message":   "success",
				"code":      200,
				"success":   true,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		err = handlers.RedisClient.SetStrWithExpire(c, "robot:card_event:"+cardAction.UserID+":"+cardAction.OpenMessageID, "Event has been handle", 25200)
		if err != nil {
			return
		}
		handler, err := handlers.CardHandler(c, cardAction)
		if err != nil {
			c.JSON(400, gin.H{
				"message":   err.Error(),
				"code":      400,
				"success":   false,
				"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
		c.JSON(200, handler)
	}

}
