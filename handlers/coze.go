package handlers

import (
	"FeishuCozeRobot/model"
	"FeishuCozeRobot/utils"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func CozeSendStream(a *ActionInfo) {
	userId := *a.Info.UserId
	message := a.Info.QParsed
	MessageId, err := SendOnProcessCard(*a.Ctx, a.Info)
	a.Info.MsgId = MessageId
	//TODO 多租户
	//BotID、BotToken多租户的情况下，
	//可以根据用户id获取对应的BotID、BotToken，
	//自行扩展，目前使用config/app.json中默认的BotID、BotToken
	data := gencozeParams(*a.Ctx, userId, message, cozeconfig.BotID)
	resp, err := SendCozeStreamMessage(data, cozeconfig.BotToken)
	if err != nil {
		println("发送请求失败： %v", err)
		return
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var answer = ""
	for scanner.Scan() {
		jsonStr := scanner.Text()
		if strings.HasPrefix(jsonStr, "data:") {
			jsonStr = jsonStr[5:]
		}
		var res model.CozeEvent
		err2 := json.Unmarshal([]byte(jsonStr), &res)
		if err2 == nil {
			switch res.Event {
			case "message":
				if res.Message.Role == "assistant" && res.Message.Type == "answer" {
					content := fmt.Sprintf("%v", res.Message.Content)
					answer += content
					//同时判断消息内容长度
					//如果大于5，并且不是图片，则更新卡片消息
					//减少触发更新卡片消息频率，导致限流的几率
					if len(content) > 5 && !utils.ContentMatch(content, "![图片名称]|![图片]") {
						UpdateTextCard(*a.Ctx, a.Info, answer)
					}
				}
			case "error":
				UpdateTextMdCardEnd(*a.Ctx, a.Info, "哎呀呀，我出错了，请将下方错误信息，请反馈给管理员：\n"+res.ErrorInformation.Msg)
				return
			case "done":
				//发送最终卡片消息
				//储存最终上下文参数
				UpdateTextMdCardEnd(*a.Ctx, a.Info, replaceImageLinks(answer, a.Info.Appid))
				setCozeParams(*a.Ctx, userId, answer)
				return
			}
		}
	}
}

// 获取上下文并生成参数
func gencozeParams(ctx context.Context, userId string, question string, BotID string) model.CozeMessage {
	messages := GetCozeMessageContext(ctx, userId)
	messages = append(messages, model.ChatHistory{Role: "user", ContentType: "text", Content: question})
	newMessage := utils.ChecklenCoze(messages)
	SetCozeMessageContext(ctx, userId, newMessage)

	return model.CozeMessage{
		BotId:       BotID,
		User:        userId,
		Stream:      true,
		Query:       question,
		ChatHistory: newMessage,
	}
}

// 保存上下文缓存
func setCozeParams(ctx context.Context, userId string, question string) {
	messages := GetCozeMessageContext(ctx, userId)
	messages = append(messages, model.ChatHistory{Role: "assistant", Type: "answer", ContentType: "text", Content: question})
	SetCozeMessageContext(ctx, userId, messages)
}

// SendCozeStreamMessage 发送扣子消息
func SendCozeStreamMessage(data interface{}, BotToken string) (*http.Response, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		log.Printf("参数加密报错： %v", err)
		return nil, err
	}
	// 构造请求对象
	req, err := http.NewRequest("POST", "https://api.coze.cn/open_api/v2/chat", strings.NewReader(string(marshal)))
	if err != nil {
		log.Printf("构造请求对象： %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+BotToken)

	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("发起请求报错： %v", err)
		return nil, err
	}

	return resp, nil
}

// 将MD中的图片链接替换为飞书的图片链接
func replaceImageLinks(text string, appid string) string {
	pattern := `!\[图片名称\]\((.*?)\)|!\[图片\]\((.*?)\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			oldLink := match[0]
			newLink := getImageKey(oldLink, appid)
			text = strings.Replace(text, oldLink, newLink, 1)
		}
		return text
	}
	return text

}

// 将MD中的图片链接迭代为飞书图片Key
func getImageKey(text string, appid string) string {
	pattern := `!\[图片名称\]\((.*?)\)|!\[图片\]\((.*?)\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		imageByte, cx := GetImageByte(matches[1])
		if cx {
			imgKey, err := uploadImageForByte(context.TODO(), imageByte, appid)
			if err == nil {
				return "![](" + *imgKey + ")"
			}
		}
	}
	return ""
}
