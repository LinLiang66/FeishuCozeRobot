package handlers

import (
	"FeishuCozeRobot/model"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"io"
	"net/http"
	"strings"
	"time"
)

type CardKind string
type CardChatType string

var (
	ModelCardKind      = CardKind("domain_version")   // 内置模型切换
	ClearCardKind      = CardKind("clear")            // 清空上下文
	PicModeChangeKind  = CardKind("pic_mode_change")  // 切换图片创作模式
	PicResolutionKind  = CardKind("pic_resolution")   // 图片分辨率调整
	PicTextMoreKind    = CardKind("pic_text_more")    // 重新根据文本生成图片
	PicVarMoreKind     = CardKind("pic_var_more")     // 变量图片
	RoleTagsChooseKind = CardKind("role_tags_choose") // 内置角色所属标签选择
	RoleChooseKind     = CardKind("role_choose")      // 内置角色选择
	PraiseKind         = CardKind("praise")           // 内置角色选择
	TemperatureKind    = CardKind("temperature")      // 内置模式
)

var (
	GroupChatType = CardChatType("group")
	UserChatType  = CardChatType("personal")
)

type CardMsg struct {
	Kind      CardKind
	ChatType  CardChatType
	Value     interface{}
	SessionId string
	MsgId     string
	APPID     string
	UserId    string
	Type      string
}

type MenuOption struct {
	value string
	label string
}

// 发送仅特定人可见的消息卡片
func sendPrivacyCardMessage(
	MsgInfo *MsgInfo,
	cardContent *larkcard.MessageCard,
	appid string,
) error {
	client := GetLarkClient(MsgInfo.ctx, appid)

	// 发起请求
	apiResp, err := client.Do(context.Background(),
		&larkcore.ApiReq{
			HttpMethod: http.MethodPost,
			ApiPath:    "https://open.feishu.cn/open-apis/ephemeral/v1/send",
			Body: &model.PrivacyCardMessageRequestBody{Card: cardContent,
				MsgType: larkim.MsgTypeInteractive,
				UserId:  MsgInfo.UserId,
				ChatId:  MsgInfo.chatId},
			SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
		},
	)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	res := &model.CreateMessageResp{}
	err = json.Unmarshal(apiResp.RawBody, &res)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func replyCard(ctx context.Context,
	msgId *string,
	cardContent *larkcard.MessageCard,
	appid string,
) error {
	client := GetLarkClient(ctx, appid)
	cardContentString, err := cardContent.String()
	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Uuid(uuid.New().String()).
			Content(cardContentString).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

func replyCardV2(ctx context.Context,
	msgId *string,
	cardContent string,
	appid string,
) error {
	client := GetLarkClient(ctx, appid)
	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Uuid(uuid.New().String()).
			Content(cardContent).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

func replyCardWithBackId(ctx context.Context,
	msgId *string,
	cardContent *larkcard.MessageCard,
	appid string,
) (*string, error) {
	client := GetLarkClient(ctx, appid)
	cardContentString, err := cardContent.String()
	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Uuid(uuid.New().String()).
			Content(cardContentString).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil, errors.New(resp.Msg)
	}

	return resp.Data.MessageId, nil
}

func NewSendCardV2(header string, MainMd string, Note string) *larkcard.MessageCard {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(true).
		Build()
	var aElementPool []larkcard.MessageCardElement
	aElementPool = append(aElementPool, WithMainMd(MainMd))
	aElementPool = append(aElementPool, WithNote(Note))
	// 卡片消息体
	cardContent := larkcard.NewMessageCard().
		Config(config).
		Header(WithHeader(header, larkcard.TemplateBlue)).
		Elements(
			aElementPool,
		)
	return cardContent
}

func NewSendCard(header *larkcard.MessageCardHeader, elements ...larkcard.MessageCardElement) *larkcard.MessageCard {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(true).
		Build()
	var aElementPool []larkcard.MessageCardElement
	for _, element := range elements {
		aElementPool = append(aElementPool, element)
	}
	// 卡片消息体
	cardContent := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(
			aElementPool,
		)
	return cardContent
}

func NewSendCardV3(header *larkcard.MessageCardHeader, elements []larkcard.MessageCardElement) *larkcard.MessageCard {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(true).
		Build()
	// 卡片消息体
	cardContent := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(
			elements,
		)
	return cardContent
}
func newSendCardWithOutHeader(
	elements ...larkcard.MessageCardElement) *larkcard.MessageCard {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(true).
		Build()
	var aElementPool []larkcard.MessageCardElement
	for _, element := range elements {
		aElementPool = append(aElementPool, element)
	}
	// 卡片消息体
	cardContent := larkcard.NewMessageCard().
		Config(config).
		Elements(
			aElementPool,
		)

	return cardContent
}

func newSimpleSendCard(
	elements ...larkcard.MessageCardElement) *larkcard.MessageCard {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(false).
		Build()
	var aElementPool []larkcard.MessageCardElement
	for _, element := range elements {
		aElementPool = append(aElementPool, element)
	}
	// 卡片消息体
	cardContent := larkcard.NewMessageCard().
		Config(config).
		Elements(
			aElementPool,
		)
	return cardContent
}

// withSplitLine 用于生成分割线
func withSplitLine() larkcard.MessageCardElement {
	splitLine := larkcard.NewMessageCardHr().
		Build()
	return splitLine
}

// WithHeader 用于生成消息头
func WithHeader(title string, color string) *larkcard.
	MessageCardHeader {
	if title == "" {
		title = "🤖️机器人提醒"
	}
	header := larkcard.NewMessageCardHeader().
		Template(color).
		Title(larkcard.NewMessageCardPlainText().
			Content(title).
			Build()).
		Build()
	return header
}

// WithNote 用于生成纯文本脚注
func WithNote(note string) larkcard.MessageCardElement {
	noteElement := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{larkcard.NewMessageCardPlainText().
			Content(note).
			Build()}).
		Build()
	return noteElement
}

// WithMainMd 用于生成markdown消息体
func WithMainMd(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = processNewLine(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(msg).
				Build()).
			IsShort(true).
			Build()}).
		Build()
	return mainElement
}

// WithMainPersonList  用于生成UserList消息体
func WithMainPersonList(userids []string) larkcard.MessageCardElement {
	userList := make([]model.Persons, len(userids))
	for i, userid := range userids {
		userList[i] = model.Persons{Id: userid}
	}
	return &model.PersonList{
		Persons:    userList,
		Size:       "small",
		Lines:      3,
		ShowAvatar: true,
		ShowName:   true,
	}
}

// withMainText 用于生成纯文本消息体
func withMainText(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = cleanTextBlock(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardPlainText().
				Content(msg).
				Build()).
			IsShort(false).
			Build()}).
		Build()
	return mainElement
}

// withMdImageCard 用于生成带有图片的消息体
func withMdImageCard(msg string, imageKey string) []larkcard.MessageCardElement {
	if imageKey == "" {
		return []larkcard.MessageCardElement{withMainText(msg)}
	}
	return []larkcard.MessageCardElement{WithMainMd(msg), withImageDiv(imageKey)}
}

func withImageDiv(imageKey string) larkcard.MessageCardElement {
	imageElement := larkcard.NewMessageCardImage().
		ImgKey(imageKey).
		Alt(larkcard.NewMessageCardPlainText().Content("").
			Build()).
		Preview(true).
		Mode(larkcard.MessageCardImageModelCropCenter).
		CompactWidth(true).
		Build()
	return imageElement
}

func newBtn(content string, value map[string]interface{},
	typename larkcard.MessageCardButtonType) *larkcard.
	MessageCardEmbedButton {
	btn := larkcard.NewMessageCardEmbedButton().
		Type(typename).
		Value(value).
		Text(larkcard.NewMessageCardPlainText().
			Content(content).
			Build())
	return btn
}

func newMenu(
	placeHolder string,
	value map[string]interface{},
	options ...MenuOption,
) *larkcard.
	MessageCardEmbedSelectMenuStatic {
	var aOptionPool []*larkcard.MessageCardEmbedSelectOption
	for _, option := range options {
		aOption := larkcard.NewMessageCardEmbedSelectOption().
			Value(option.value).
			Text(larkcard.NewMessageCardPlainText().
				Content(option.label).
				Build())
		aOptionPool = append(aOptionPool, aOption)

	}
	btn := larkcard.NewMessageCardEmbedSelectMenuStatic().
		MessageCardEmbedSelectMenuStatic(larkcard.NewMessageCardEmbedSelectMenuBase().
			Options(aOptionPool).
			Placeholder(larkcard.NewMessageCardPlainText().
				Content(placeHolder).
				Build()).
			Value(value).
			Build()).
		Build()
	return btn
}

// 清除卡片按钮
func withClearDoubleCheckBtn(sessionID *string) larkcard.MessageCardElement {
	confirmBtn := newBtn("确认清除", map[string]interface{}{
		"value":     "1",
		"kind":      ClearCardKind,
		"chatType":  UserChatType,
		"sessionID": *sessionID,
	}, "danger_filled",
	)
	cancelBtn := newBtn("我再想想", map[string]interface{}{
		"value":     "0",
		"kind":      ClearCardKind,
		"sessionID": *sessionID,
		"chatType":  UserChatType,
	},
		larkcard.MessageCardButtonTypeDefault)

	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{confirmBtn, cancelBtn}).
		Layout(larkcard.MessageCardActionLayoutBisected.Ptr()).
		Build()

	return actions
}

// 赞一下
func withLikeButton(sessionID *string, appid string) larkcard.MessageCardElement {
	confirmBtn := newBtn("赞一下", map[string]interface{}{
		"value":     "1",
		"kind":      PraiseKind,
		"chatType":  UserChatType,
		"appid":     appid,
		"sessionID": *sessionID,
	}, larkcard.MessageCardButtonTypePrimary,
	)
	cancelBtn := newBtn("踩一下", map[string]interface{}{
		"value":     "0",
		"kind":      PraiseKind,
		"appid":     appid,
		"sessionID": *sessionID,
		"chatType":  UserChatType,
	},
		larkcard.MessageCardButtonTypeDanger)

	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{confirmBtn, cancelBtn}).
		Layout(larkcard.MessageCardActionLayoutBisected.Ptr()).
		Build()

	return actions
}

func replyMsg(ctx context.Context, msg string, msgId *string, appid string) error {
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := GetLarkClient(ctx, appid)
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

// GetImageByte 下载图片文件
func GetImageByte(url string) ([]byte, bool) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()
	ContentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ContentType, "image/") {
		return nil, false
	}
	body, err := io.ReadAll(resp.Body)
	return body, err == nil
}

func uploadImageForByte(ctx context.Context, imageBytes []byte, appid string) (*string, error) {
	client := GetLarkClient(ctx, appid)
	resp, err := client.Im.Image.Create(context.Background(),
		larkim.NewCreateImageReqBuilder().
			Body(larkim.NewCreateImageReqBodyBuilder().
				ImageType(larkim.ImageTypeMessage).
				Image(bytes.NewReader(imageBytes)).
				Build()).
			Build())
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil, errors.New(resp.Msg)
	}
	return resp.Data.ImageKey, nil
}

func uploadImageForBase64(ctx context.Context, base64Str string, appid string) (*string, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client := GetLarkClient(ctx, appid)
	resp, err := client.Im.Image.Create(context.Background(),
		larkim.NewCreateImageReqBuilder().
			Body(larkim.NewCreateImageReqBodyBuilder().
				ImageType(larkim.ImageTypeMessage).
				Image(bytes.NewReader(imageBytes)).
				Build()).
			Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil, errors.New(resp.Msg)
	}
	return resp.Data.ImageKey, nil
}

func sendMsg(ctx context.Context, msg string, chatId *string, appid string) error {
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := GetLarkClient(ctx, appid)
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(*chatId).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

func PatchCard(ctx context.Context, msgId *string,
	cardContent *larkcard.MessageCard, appid string) error {
	client := GetLarkClient(ctx, appid)
	cardContentString, err := cardContent.String()
	resp, err := client.Im.Message.Patch(ctx, larkim.NewPatchMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewPatchMessageReqBodyBuilder().
			Content(cardContentString).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

func sendClearCacheCheckCard(ctx context.Context,
	msgInfo *MsgInfo) {
	newCard := NewSendCard(
		WithHeader("🆑 机器人提醒", larkcard.TemplateBlue),
		WithMainMd("您确定要清除对话上下文吗？"),
		WithNote("请注意，这将开始一个全新的对话，您将无法利用之前话题的历史信息"),
		withClearDoubleCheckBtn(msgInfo.sessionId))
	if msgInfo.handlerType == "group" {
		err := sendPrivacyCardMessage(msgInfo, newCard, msgInfo.Appid)
		if err != nil {
			return
		}

	} else if msgInfo.handlerType == "personal" {
		err := replyCard(ctx, msgInfo.MsgId, newCard, msgInfo.Appid)
		if err != nil {
			return
		}
	}

}

func UpdateTextCardEnd(ctx context.Context, msgInfo *MsgInfo, Message string) error {
	TimeLimit := time.Now().Format("2006-01-02 15:04:05")
	newCard := NewSendCard(
		WithHeader("🎉 处理结果", larkcard.TemplateViolet),
		larkcard.NewMessageCardMarkdown().
			Content("**🕐 响应时间：**\n"+TimeLimit).
			Build(),
		larkcard.NewMessageCardMarkdown().
			Content("**"+Message+"**").
			Build(),
		withLikeButton(msgInfo.sessionId, msgInfo.Appid),
		WithNote("🤖温馨提示✨✨：输入<帮助> 或 help 即可获取帮助菜单"))
	err := PatchCard(ctx, msgInfo.MsgId, newCard, msgInfo.Appid)
	if err != nil {
		return err
	}
	return nil
}
func UpdateTextMdCardEnd(ctx context.Context, msgInfo *MsgInfo, Message string) error {
	TimeLimit := time.Now().Format("2006-01-02 15:04:05")
	newCard := NewSendCard(
		WithHeader("🎉 处理结果", larkcard.TemplateViolet),
		larkcard.NewMessageCardMarkdown().
			Content("**🕐 响应时间：**\n"+TimeLimit).
			Build(),
		larkcard.NewMessageCardMarkdown().
			Content(Message).
			Build(),
		withLikeButton(msgInfo.sessionId, msgInfo.Appid),
		WithNote("🤖温馨提示✨✨：输入<帮助> 或 help 即可获取帮助菜单"))
	err := PatchCard(ctx, msgInfo.MsgId, newCard, msgInfo.Appid)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTextCard(ctx context.Context, msgInfo *MsgInfo, Message string) error {
	newCard := newSendCardWithOutHeader(
		withMainText(Message),
		WithNote("🤖温馨提示:正在处理中，请稍等..."))
	err := PatchCard(ctx, msgInfo.MsgId, newCard, msgInfo.Appid)
	if err != nil {
		return err
	}
	return nil
}

func SendOnProcessCard(ctx context.Context, msgInfo *MsgInfo) (*string, error) {
	newCard := newSendCardWithOutHeader(
		WithNote("🤖温馨提示:正在思考中，请稍等..."))
	id, err := replyCardWithBackId(ctx, msgInfo.MsgId, newCard, msgInfo.Appid)
	if err != nil {
		return nil, err
	}
	return id, nil
}
