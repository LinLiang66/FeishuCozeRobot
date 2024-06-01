package handlers

import (
	"FeishuCozeRobot/model"
	"context"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"time"
)

func NewPraiseCardHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *model.CardAction) (interface{}, error) {
		if cardMsg.Kind == PraiseKind {
			cardMsg.MsgId = cardAction.OpenMessageID
			newCard, err, done := CommonProcessPraiseCache(ctx, cardMsg)

			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func CommonProcessPraiseCache(ctx context.Context, cardMsg CardMsg) (
	interface{}, error, bool) {
	// 业务处理
	messagedata := getMessage(ctx, cardMsg.APPID, cardMsg.MsgId)
	if messagedata.Success() {
		Message := getTextFromJson(*messagedata.Data.Items[0].Body.Content)
		TimeLimit := time.Now().Format("2006-01-02 15:04:05")
		newCard := NewSendCard(
			WithHeader("🎉 处理结果", larkcard.TemplateViolet),
			larkcard.NewMessageCardMarkdown().
				Content("**🕐 完成时间：**\n"+TimeLimit).
				Build(),
			larkcard.NewMessageCardMarkdown().
				Content("**"+Message+"**").
				Build(),
			larkcard.NewMessageCardDiv().
				Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
					Text(larkcard.NewMessageCardLarkMd().
						Content("**📝祝您生活愉快**").
						Build()).
					IsShort(true).
					Build()}).
				Build(),
			WithNote("🤖温馨提示✨✨：输入<帮助> 或 help 即可获取帮助菜单"))
		return newCard, nil, true
	}
	return nil, nil, false

}
