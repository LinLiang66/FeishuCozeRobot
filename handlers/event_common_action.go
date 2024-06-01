package handlers

import (
	"FeishuCozeRobot/model"
	"FeishuCozeRobot/utils"
	"context"
)

type MsgInfo struct {
	handlerType HandlerType
	msgType     string
	MsgId       *string
	chatId      *string
	UserId      *string
	QParsed     string
	fileKey     []model.File
	imageKey    []string
	sessionId   *string
	Appid       string
	mention     []*model.Mention
	ctx         context.Context
	BotId       string
	BotToken    string
}

type ActionInfo struct {
	handler *MessageHandler
	Ctx     *context.Context
	Info    *MsgInfo
}

type Action interface {
	Execute(a *ActionInfo) bool
}

// ProcessedUniqueAction 消息去重处理
type ProcessedUniqueAction struct { //幂等判断消息唯一性
}

func (*ProcessedUniqueAction) Execute(a *ActionInfo) bool {
	if RedisClient.KEYEXISTS(*a.Ctx, "robot:message_event:"+*a.Info.MsgId) {
		return false
	}
	err := RedisClient.SetStrWithExpire(*a.Ctx, "robot:message_event:"+*a.Info.MsgId, "Message has been handle", 25200)
	if err != nil {
		return false
	}
	return true
}

// RobotAction   扣子大模型兜底处理
type RobotAction struct { /*大模型兜底处理*/
}

func (*RobotAction) Execute(a *ActionInfo) bool {
	go CozeSendStream(a)
	return false
}

type ProcessMentionAction struct { //是否机器人应该处理
}

func (*ProcessMentionAction) Execute(a *ActionInfo) bool {
	//私聊直接过
	if a.Info.handlerType == UserHandler {
		return true
	}
	// 群聊判断是否提到机器人
	if a.Info.handlerType == GroupHandler {
		return a.handler.judgeIfMentionMe(a.Info.mention)
	}
	return false
}

// EmptyAction 空内容处理
type EmptyAction struct { /*空消息*/
}

func (*EmptyAction) Execute(a *ActionInfo) bool {
	if len(a.Info.QParsed) == 0 && a.Info.msgType != "file" {
		sendMsg(*a.Ctx, "🤖️：您好，请问有什么可以帮到您~", a.Info.chatId, a.Info.Appid)
		return false
	}
	return true
}

// ClearAction 清空上下文
type ClearAction struct { /*清除消息*/
}

func (*ClearAction) Execute(a *ActionInfo) bool {
	if foundClear, _ := utils.ContainsSpecificContent(a.Info.QParsed,
		"clear|清除"); foundClear {
		sendClearCacheCheckCard(*a.Ctx, a.Info)
		return false
	}
	return true
}

type RolePlayAction struct { /*角色扮演*/
}
