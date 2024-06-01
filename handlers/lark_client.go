package handlers

import (
	"context"
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

var AppClient map[string]*lark.Client

func init() {
	AppClient = make(map[string]*lark.Client)
}

func GetLarkClient(ctx context.Context, appid string) *lark.Client {
	client := AppClient[appid]
	if client != nil {
		return client
	}
	appCache, _ := GetAppCache(ctx, appid)
	NewClient := lark.NewClient(appCache.AppID, appCache.AppSecret)
	AppClient[appid] = NewClient
	return NewClient
}
