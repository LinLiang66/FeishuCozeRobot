package handlers

import (
	"FeishuCozeRobot/model"
	"FeishuCozeRobot/utils"
	"context"
	"fmt"
	"strings"
)

// 责任链
func chain(data *ActionInfo, actions ...Action) bool {
	for _, v := range actions {
		if !v.Execute(data) {
			return false
		}
	}
	return true
}

type MessageHandler struct {
}

func (m MessageHandler) cardHandler(ctx context.Context,
	cardAction *model.CardAction) (interface{}, error) {
	messageHandler := NewCardHandler(m)
	return messageHandler(ctx, cardAction)
}

func judgeMsgType(event *model.MessageEvent) (string, error) {
	msgType := event.Event.Message.MessageType
	switch msgType {
	case "text", "image", "audio", "post", "file":
		return msgType, nil
	default:
		return "", fmt.Errorf("unknown message type: %v", msgType)
	}
}

func (m MessageHandler) msgReceivedHandler(ctx context.Context, event *model.MessageEvent) error {
	handlerType := judgeChatType(event)
	if handlerType == "otherChat" {
		return nil
	}
	msgType, err := judgeMsgType(event)
	if err != nil {
		fmt.Printf("error getting message type: %v\n", err)
		return nil
	}
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageID
	chatId := event.Event.Message.ChatID
	mention := event.Event.Message.Mentions
	sessionId := event.Event.Sender.SenderID.UserID
	if sessionId == nil || *sessionId == "" {
		sessionId = event.Event.Sender.SenderID.OpenID
	}
	var qParsed string
	var file []model.File
	var imageKey []string
	if msgType == "post" {
		qParsed, file, imageKey = getMessag(*content)
	} else if msgType == "file" {
		file = append(file, parseFileKey(*content))
	} else {
		qParsed = strings.Trim(parseContent(*content), " ")
		file = append(file, parseFileKey(*content))
		imageKey = append(imageKey, parseImageKey(*content))
	}
	msgInfo := MsgInfo{
		handlerType: handlerType,
		msgType:     msgType,
		MsgId:       msgId,
		UserId:      event.Event.Sender.SenderID.UserID,
		chatId:      chatId,
		QParsed:     qParsed,
		fileKey:     file,
		imageKey:    imageKey,
		sessionId:   sessionId,
		Appid:       event.Header.AppID,
		mention:     mention,
	}
	data := &ActionInfo{
		Ctx:     &ctx,
		handler: &m,
		Info:    &msgInfo,
	}
	actions := []Action{
		&ProcessedUniqueAction{}, //避免重复处理
		&EmptyAction{},           //空消息处理
		&ClearAction{},           //清除消息处理
		&ProcessMentionAction{},  //判断机器人是否应该被调用
		&RobotAction{},           //大模型兜底处理
	}
	chain(data, actions...)
	return nil
}

func (m MessageHandler) judgeIfMentionMe(mention []*model.Mention) bool {
	if len(mention) != 1 {
		return false
	}
	return utils.ContainsSpecificContentV2(*mention[0].Name, "小航|小诺|小肉")
}
