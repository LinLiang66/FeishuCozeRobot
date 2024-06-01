package model

// CozeMessage 请求内容
type CozeMessage struct {
	ConversationId string        `json:"conversation_id"`
	BotId          string        `json:"bot_id"`
	User           string        `json:"user"`
	Query          string        `json:"query"`
	Stream         bool          `json:"stream"`
	ChatHistory    []ChatHistory `json:"chat_history,omitempty"`
}

type ChatHistory struct {
	Role        string `json:"role"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Type        string `json:"type,omitempty"`
}

// CozeEvent 流式响应内容
type CozeEvent struct {
	Event   string `json:"event"`
	Message struct {
		Role        string      `json:"role"`
		Type        string      `json:"type"`
		Content     interface{} `json:"content"`
		ContentType string      `json:"content_type"`
	} `json:"message,omitempty"`
	IsFinish         bool   `json:"is_finish"`
	Index            int    `json:"index"`
	ConversationId   string `json:"conversation_id"`
	SeqId            int    `json:"seq_id"`
	ErrorInformation struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error_information,omitempty"`
}
