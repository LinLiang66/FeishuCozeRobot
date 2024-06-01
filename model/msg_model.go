package model

import larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

type CreateMessageResp struct {
	Code int                   `json:"code"`
	Data CreateMessageRespData `json:"data"`
	Msg  string                `json:"msg"`
}

type CreateMessageRespData struct {
	MessageID string `json:"message_id"`
}

type Button struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
	Type string `json:"type"`
}

type Note struct {
	Elements []Element `json:"elements"`
}

type FeiShuMessage struct {
	Title    string      `json:"title"`
	Elements [][]Element `json:"elements"`
}

type Element struct {
	Href      string   `json:"href,omitempty"`
	UserId    string   `json:"user_id,omitempty"`
	Style     []string `json:"style,omitempty"`
	ImageKey  string   `json:"image_key,omitempty"`
	FileKey   string   `json:"file_key,omitempty"`
	EmojiType string   `json:"emoji_type,omitempty"`
	FileName  string   `json:"file_name,omitempty"`
	Tag       string   `json:"tag,omitempty"`
	Text      string   `json:"text,omitempty"`
	TextAlign string   `json:"text_align,omitempty"`
	Buttons   []Button `json:"buttons,omitempty"`
	Note      *Note    `json:"note,omitempty"`
}

type SendTextMessage struct {
	MsgType   string                `json:"msg_type,omitempty"`
	Content   Content               `json:"content,omitempty"`
	Card      *larkcard.MessageCard `json:"card,omitempty"`
	Sign      string                `json:"sign,omitempty"`
	Timestamp string                `json:"timestamp,omitempty"`
}

type Content struct {
	Text string `json:"text,omitempty"`
}

type Event struct {
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
}

type MessageEvent struct {
	Schema    string  `json:"schema"`
	Header    Header  `json:"header"`
	Event     Event   `json:"event"`
	Sender    Sender  `json:"sender"`
	Message   Message `json:"message"`
	Type      string  `json:"type"`
	Challenge string  `json:"challenge"`
}

type Header struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
}

type Sender struct {
	SenderID   SenderID `json:"sender_id"`
	SenderType string   `json:"sender_type"`
	TenantKey  string   `json:"tenant_key"`
}

type SenderID struct {
	UnionID *string `json:"union_id"`
	UserID  *string `json:"user_id"`
	OpenID  *string `json:"open_id"`
}

type Message struct {
	MessageID   *string    `json:"message_id"`
	RootID      string     `json:"root_id"`
	ParentID    string     `json:"parent_id"`
	CreateTime  string     `json:"create_time"`
	UpdateTime  string     `json:"update_time"`
	ChatID      *string    `json:"chat_id"`
	ChatType    string     `json:"chat_type"`
	MessageType string     `json:"message_type"`
	Content     *string    `json:"content"`
	Mentions    []*Mention `json:"mentions"`
	UserAgent   string     `json:"user_agent"`
}

type Mention struct {
	Key       string  `json:"key"`
	ID        ID      `json:"id"`
	Name      *string `json:"name"`
	TenantKey string  `json:"tenant_key"`
}

type ID struct {
	UnionID string `json:"union_id"`
	UserID  string `json:"user_id"`
	OpenID  string `json:"open_id"`
}

type File struct {
	FileKey  string `json:"file_key,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

type PostMessage struct {
	Title string      `json:"title"`
	Data  [][]Element `json:"content"`
}

type PrivacyCardMessageRequestBody struct {
	ChatId  *string               `json:"chat_id"`
	UserId  *string               `json:"user_id"`
	MsgType string                `json:"msg_type"`
	Card    *larkcard.MessageCard `json:"card"`
}
