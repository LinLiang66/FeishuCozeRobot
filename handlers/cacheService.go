package handlers

import (
	"FeishuCozeRobot/model"
	"context"
	"encoding/json"
	"log"
)

type AppCache struct {
	UserName          string  `json:"user_name,omitempty"`
	UserID            string  `json:"user_id,omitempty"`
	AppID             string  `json:"appid,omitempty"`
	AppSecret         string  `json:"app_secret,omitempty"`
	AppRoleType       float64 `json:"app_role_type,omitempty"`
	VerificationToken string  `json:"verification_token,omitempty"`
	EncryptKey        string  `json:"encrypt_key,omitempty"`
	RobotAppID        string  `json:"robot_appid,omitempty"`
	RobotApiSecret    string  `json:"robot_api_secret,omitempty"`
	RobotApiKey       string  `json:"robot_api_key,omitempty"`
	RobotDomain       string  `json:"robot_domain,omitempty"`
	RobotSparkUrl     string  `json:"robot_spark_url,omitempty"`
	RobotTemperature  float64 `json:"robot_temperature,omitempty"`
	JobTitle          string  `json:"job_title,omitempty"`
	JobNumber         string  `json:"job_number,omitempty"`
}

// GetAppCache 获取飞书应用信息
func GetAppCache(ctx context.Context, appId string) (AppCache, bool) {
	if RedisClient.KEYEXISTS(ctx, "robot:robot_app_key:"+appId) {
		str, err := RedisClient.GetStr(ctx, "robot:robot_app_key:"+appId)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		var appCache AppCache
		err = json.Unmarshal([]byte(str), &appCache)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		return appCache, true
	}
	return AppCache{}, false
}

func GetCozeMessageContext(Ctx context.Context, userid string) []model.ChatHistory {
	var MessageContext []model.ChatHistory
	if RedisClient.KEYEXISTS(Ctx, "robot:message_context_coze:"+userid) {
		str, err := RedisClient.GetStr(Ctx, "robot:message_context_coze:"+userid)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		err = json.Unmarshal([]byte(str), &MessageContext)
		if err != nil {
			log.Printf("failed to start server: %v", err)
		}
		return MessageContext
	}
	return MessageContext
}

func SetCozeMessageContext(Ctx context.Context, userid string, MessageContext []model.ChatHistory) {
	bytes, _ := json.Marshal(MessageContext)
	RedisClient.SetStr(Ctx, "robot:message_context_coze:"+userid, string(bytes), 0)
}
