package handlers

import (
	"FeishuCozeRobot/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type CardHandlerMeta func(cardMsg CardMsg, m MessageHandler) CardHandlerFunc

type CardHandlerFunc func(ctx context.Context, cardAction *model.CardAction) (
	interface{}, error)

var ErrNextHandler = fmt.Errorf("next handler")

func NewCardHandler(m MessageHandler) CardHandlerFunc {
	handlers := []CardHandlerMeta{
		NewClearCardHandler,
		NewPraiseCardHandler,
	}
	return func(ctx context.Context, cardAction *model.CardAction) (interface{}, error) {
		var cardMsg CardMsg
		if cardAction.Action == nil {
			return nil, errors.New("card action is nil")
		}
		actionValue := cardAction.Action.Value
		actionValueJson, err := json.Marshal(actionValue)
		if err != nil {
			fmt.Println("转换失败", err.Error())
			return nil, err
		}
		err = json.Unmarshal(actionValueJson, &cardMsg)
		if err != nil {
			return nil, err
		}
		for _, handler := range handlers {
			h := handler(cardMsg, m)
			i, err := h(ctx, cardAction)
			if errors.Is(err, ErrNextHandler) {
				continue
			}
			return i, err
		}
		return nil, nil
	}
}
