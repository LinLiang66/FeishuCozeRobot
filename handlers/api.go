package handlers

import (
	"FeishuCozeRobot/model"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// DownloadMessageFile 获取消息中的资源文件
func DownloadMessageFile(ctx context.Context, appid string, MsgId string, FileKey string, FilieType string, FilePath string) bool {
	if len(FileKey) == 0 {
		return false
	}
	var try int
	var res bool
	for { // 重试次数为3次
		res = DownloadMessageFileV2(ctx, appid, MsgId, FileKey, FilieType, FilePath)
		if res || try == 3 {
			break
		}
		try++
		time.Sleep(1 * time.Second) // 延时5秒
	}
	return res
}
func DownloadMessageFileV2(ctx context.Context, appid string, MsgId string, FileKey string, FilieType string, FilePath string) bool {
	client := GetLarkClient(ctx, appid)
	// 创建请求对象
	req := larkim.NewGetMessageResourceReqBuilder().
		MessageId(MsgId).
		FileKey(FileKey).
		Type(FilieType).
		Build()
	// 发起请求
	resp, err := client.Im.MessageResource.Get(ctx, req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return false
	}
	// 业务处理
	err = resp.WriteFile(FilePath)
	if err != nil {
		log.Printf("下载附件报错: %v", err.Error())
		return false
	}
	return true
}

// getMessage 飞书根据messageId获取历史消息
func getMessage(ctx context.Context, appId string, messageId string) *larkim.GetMessageResp {
	client := GetLarkClient(ctx, appId)
	// 创建请求对象
	req := larkim.NewGetMessageReqBuilder().
		MessageId(messageId).
		Build()

	resp, err := client.Im.Message.Get(ctx, req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return resp
	}

	return resp
}

// GetUser 飞书根据UserId 获取用户信息
func GetUser(Ctx context.Context, Appid string, UserId string) *larkcontact.GetUserResp {

	client := GetLarkClient(Ctx, Appid)
	// 创建请求对象
	req := larkcontact.NewGetUserReqBuilder().
		UserId(UserId).
		UserIdType(`user_id`).
		Build()
	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := client.Contact.User.Get(Ctx, req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return resp
	}
	return resp
}

// getTextFromJson 获取飞书指定消息内容
func getTextFromJson(jsonStr string) string {
	var data model.FeiShuMessage
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return ""
	}
	var Message string
	for i := 1; i < len(data.Elements); i++ {
		noteElement := data.Elements[i]
		for _, Element := range noteElement {
			if Element.Tag == "text" {
				Message += Element.Text
			}
		}
	}
	return Message
}

// SendCustomRobotMessage  发送飞书自定义消息
func SendCustomRobotMessage(Message model.SendTextMessage, WebhookUrl string) bool {
	// 构造请求内容
	messagebyte, err := json.Marshal(Message)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 构造请求对象
	req, err := http.NewRequest("POST", WebhookUrl, strings.NewReader(string(messagebyte)))
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var res model.CreateMessageResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Printf("请求结果转换为实体: %v\n", resp)
		return false
	}
	return res.Code == 0
}

// RemoveBase64Prefix 删除Base64前缀
func RemoveBase64Prefix(base64Str string) (string, error) {
	// 正则表达式匹配Base64前缀
	prefixPattern := `data:image\/.*;base64,`
	re := regexp.MustCompile(prefixPattern)
	// 移除前缀
	base64Str = re.ReplaceAllString(base64Str, "")
	// 检查Base64字符串长度是否为4的倍数
	padding := len(base64Str) % 4
	if padding != 0 {
		// 如果不是4的倍数，则添加'='作为填充字符
		base64Str += strings.Repeat("=", 4-padding)
	}
	// 尝试解码以验证Base64字符串是否有效
	_, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}

	return base64Str, nil
}

// ImageToBase64 根据文件路径将图片转换为Base64编码
func ImageToBase64(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	encodedString := base64.StdEncoding.EncodeToString(fileBytes)
	return encodedString, nil
}
