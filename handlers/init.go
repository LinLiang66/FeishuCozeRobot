package handlers

import (
	"FeishuCozeRobot/config"
	"FeishuCozeRobot/model"
	"context"
)

type MessageHandlerInterface interface {
	msgReceivedHandler(ctx context.Context, event *model.MessageEvent) error
	cardHandler(ctx context.Context, cardAction *model.CardAction) (interface{}, error)
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "personal"
)

// 扣子API相关配置
var cozeconfig *config.Config

// handlers 所有消息类型类型的处理器
var handlers MessageHandlerInterface

func InitHandlers(config *config.Config) {
	handlers = NewMessageHandler()
	cozeconfig = config
}

func Handler(ctx context.Context, event *model.MessageEvent) error {
	return handlers.msgReceivedHandler(ctx, event)
}

func CardHandler(ctx context.Context, cardAction *model.CardAction) (interface{}, error) {
	return handlers.cardHandler(ctx, cardAction)
}

func NewMessageHandler() MessageHandlerInterface {
	return &MessageHandler{}
}

func judgeChatType(event *model.MessageEvent) HandlerType {
	chatType := event.Event.Message.ChatType
	if chatType == "group" {
		return GroupHandler
	}
	if chatType == "p2p" {
		return UserHandler
	}
	return "otherChat"
}
