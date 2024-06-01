package handlers

import (
	"FeishuCozeRobot/model"
	"context"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

func NewClearCardHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *model.CardAction) (interface{}, error) {
		if cardMsg.Kind == ClearCardKind {
			cardMsg.MsgId = cardAction.OpenMessageID
			newCard, err, done := CommonProcessClearCache(ctx, cardMsg)
			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func CommonProcessClearCache(ctx context.Context, cardMsg CardMsg) (
	interface{}, error, bool) {

	if cardMsg.Value == "1" {
		RedisClient.DelByKey(ctx, "robot:message_context:"+cardMsg.SessionId)
		RedisClient.DelByKey(ctx, "robot:message_context_yuanqi:"+cardMsg.SessionId)
		newCard := NewSendCard(
			WithHeader("️🆑 机器人提醒", larkcard.TemplateGrey),
			WithMainMd("已删除此话题的上下文信息"),
			WithNote("我们可以开始一个全新的话题，继续找我聊天吧"),
		)
		return newCard, nil, true
	}
	if cardMsg.Value == "0" {
		newCard := NewSendCard(
			WithHeader("️🆑 机器人提醒", larkcard.TemplateGreen),
			WithMainMd("依旧保留此话题的上下文信息"),
			WithNote("我们可以继续探讨这个话题,期待和您聊天。如果您有其他问题或者想要讨论的话题，请告诉我哦"),
		)
		return newCard, nil, true
	}
	return nil, nil, false
}
