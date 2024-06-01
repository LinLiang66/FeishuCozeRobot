package model

import (
	"encoding/json"
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"reflect"
	"strings"
)

type Payload struct {
	CardLink     CardLink               `json:"card_link"`
	Config       map[string]interface{} `json:"config"`
	I18NElements I18NElements           `json:"i18n_elements"`
	I18NHeader   I18NHeader             `json:"i18n_header"`
}

type CardLink struct {
	AndroidURL string `json:"android_url"`
	IosURL     string `json:"ios_url"`
	PCURL      string `json:"pc_url"`
	URL        string `json:"url"`
}

type I18NElements struct {
	ZhCN []ZhCNElement `json:"zh_cn"`
}

type ZhCNElement struct {
	Elements []ZhCNElementClass `json:"elements"`
	Fallback *ZhCNFallback      `json:"fallback,omitempty"`
	Name     *string            `json:"name,omitempty"`
	Tag      string             `json:"tag"`
}

type ZhCNElementClass struct {
	BackgroundStyle   string         `json:"background_style"`
	Columns           []PurpleColumn `json:"columns"`
	Content           *string        `json:"content,omitempty"`
	FlexMode          string         `json:"flex_mode"`
	HorizontalAlign   string         `json:"horizontal_align"`
	HorizontalSpacing string         `json:"horizontal_spacing"`
	Margin            string         `json:"margin"`
	Tag               string         `json:"tag"`
	Token             *string        `json:"token,omitempty"`
}

type PurpleColumn struct {
	BackgroundStyle string          `json:"background_style"`
	Elements        []PurpleElement `json:"elements"`
	Tag             string          `json:"tag"`
	VerticalAlign   string          `json:"vertical_align"`
	VerticalSpacing string          `json:"vertical_spacing"`
	Weight          *int64          `json:"weight,omitempty"`
	Width           string          `json:"width"`
}

type PurpleElement struct {
	ActionType         *string            `json:"action_type,omitempty"`
	BackgroundStyle    *string            `json:"background_style,omitempty"`
	Columns            []FluffyColumn     `json:"columns,omitempty"`
	ComplexInteraction *bool              `json:"complex_interaction,omitempty"`
	Content            *string            `json:"content,omitempty"`
	FlexMode           *string            `json:"flex_mode,omitempty"`
	HorizontalAlign    *string            `json:"horizontal_align,omitempty"`
	HorizontalSpacing  *string            `json:"horizontal_spacing,omitempty"`
	Icon               *Icon              `json:"icon,omitempty"`
	Margin             *string            `json:"margin,omitempty"`
	Name               string             `json:"name"`
	Options            []Option           `json:"options,omitempty"`
	Placeholder        *FluffyPlaceholder `json:"placeholder,omitempty"`
	Required           *bool              `json:"required,omitempty"`
	Size               *string            `json:"size,omitempty"`
	Tag                string             `json:"tag"`
	Text               *ElementText       `json:"text,omitempty"`
	TextAlign          *string            `json:"text_align,omitempty"`
	TextSize           *string            `json:"text_size,omitempty"`
	Type               string             `json:"type"`
	Value              *string            `json:"value,omitempty"`
	Width              *string            `json:"width,omitempty"`
}

type FluffyColumn struct {
	BackgroundStyle *string         `json:"background_style,omitempty"`
	Elements        []FluffyElement `json:"elements,omitempty"`
	Tag             *string         `json:"tag,omitempty"`
	VerticalAlign   *string         `json:"vertical_align,omitempty"`
	VerticalSpacing *string         `json:"vertical_spacing,omitempty"`
	Weight          *int64          `json:"weight,omitempty"`
	Width           *string         `json:"width,omitempty"`
}

type FluffyElement struct {
	Content      *string            `json:"content,omitempty"`
	DefaultValue *string            `json:"default_value,omitempty"`
	Fallback     *ElementFallback   `json:"fallback,omitempty"`
	Name         *string            `json:"name,omitempty"`
	Placeholder  *PurplePlaceholder `json:"placeholder,omitempty"`
	Required     *bool              `json:"required,omitempty"`
	Tag          string             `json:"tag"`
	TextAlign    *string            `json:"text_align,omitempty"`
	TextSize     *string            `json:"text_size,omitempty"`
	Width        *string            `json:"width,omitempty"`
}

type ElementFallback struct {
	Tag  string     `json:"tag"`
	Text PurpleText `json:"text"`
}

type PurpleText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type PurplePlaceholder struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Icon struct {
	Color string `json:"color"`
	Tag   string `json:"tag"`
	Token string `json:"token"`
}

type Option struct {
	Text  OptionText `json:"text"`
	Value string     `json:"value"`
}

type OptionText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type FluffyPlaceholder struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type ElementText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type ZhCNFallback struct {
	Tag  string     `json:"tag"`
	Text FluffyText `json:"text"`
}

type FluffyText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type I18NHeader struct {
	ZhCN I18NHeaderZhCN `json:"zh_cn"`
}

type I18NHeaderZhCN struct {
	Template    string        `json:"template"`
	TextTagList []TextTagList `json:"text_tag_list"`
	Title       Title         `json:"title"`
	UdIcon      UdIcon        `json:"ud_icon"`
}

type TextTagList struct {
	Color *string          `json:"color,omitempty"`
	Tag   *string          `json:"tag,omitempty"`
	Text  *TextTagListText `json:"text,omitempty"`
}

type TextTagListText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Title struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type UdIcon struct {
	Tag   string `json:"tag"`
	Token string `json:"token"`
}

func EscapeAndFormatJson(card *Payload) string {
	formattedBytes, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		return ""
	}
	return string(formattedBytes)
}

type ColumnSet struct {
	Tag_              string   `json:"tag"`
	FlexMode          string   `json:"flex_mode"`
	BackgroundStyle   string   `json:"background_style"`
	Columns           []Column `json:"columns"`
	HorizontalSpacing string   `json:"horizontal_spacing"`
	larkcard.MessageCardElement
}

func (m *ColumnSet) Tag() string {
	return "column_set"
}

func (m *ColumnSet) MarshalJSON() ([]byte, error) {
	return messageCardColumnSetJson(m)
}

func messageCardColumnSetJson(e *ColumnSet) ([]byte, error) {
	data, err := StructToMap(e)
	if err != nil {
		return nil, err
	}
	data["tag"] = e.Tag()
	return json.Marshal(data)
}

type MarkdownElement struct {
	Tag_      string `json:"tag"`
	TextAlign string `json:"text_align"`
	Content   string `json:"content"`
	Element
}

func (m *MarkdownElement) Tag() string {
	return "markdown"
}

func (m *MarkdownElement) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

func messageCardElementJson(e *MarkdownElement) ([]byte, error) {
	data, err := StructToMap(e)
	if err != nil {
		return nil, err
	}
	data["tag"] = e.Tag()
	return json.Marshal(data)
}

type Column struct {
	Tag           string                        `json:"tag,omitempty"`
	Width         string                        `json:"width,omitempty"`
	Weight        float64                       `json:"weight,omitempty"`
	VerticalAlign string                        `json:"vertical_align,omitempty"`
	Elements      []larkcard.MessageCardElement `json:"elements,omitempty"`
}

type CardAction struct {
	Header     map[string]string
	Body       []byte
	RequestURI string
	*larkevent.EventReq
	OpenID        string `json:"open_id,omitempty"`
	UserID        string `json:"user_id,omitempty"`
	OpenMessageID string `json:"open_message_id,omitempty"`
	OpenChatId    string `json:"open_chat_id,omitempty"`
	TenantKey     string `json:"tenant_key,omitempty"`
	Token         string `json:"token,omitempty"`
	Timezone      string `json:"timezone,omitempty"`
	Challenge     string `json:"challenge,omitempty"`
	Type          string `json:"type,omitempty"`
	Appid         string `json:"app_id,omitempty"`
	Action        *struct {
		Value     interface{}       `json:"value,omitempty"`
		Tag       string            `json:"tag,omitempty"`
		Option    string            `json:"option,omitempty"`
		Timezone  string            `json:"timezone,omitempty"`
		FormValue map[string]string `json:"form_value,omitempty"`
		Name      string            `json:"name,omitempty"`
	} `json:"action,omitempty"`
}

// StructToMap 定义一个函数，将任何结构体转换为 map[string]interface{}
func StructToMap(i interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, fmt.Errorf("input is a nil pointer")
	}

	// 获取结构体的真实值（如果传入的是指针）
	value := reflect.Indirect(v)

	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input should be a struct")
	}

	m := make(map[string]interface{})
	t := value.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		// 如果结构体字段有 json 标签，则使用标签的名称作为键，否则使用结构体字段名
		key := jsonTag
		if key == "" {
			key = field.Name
		}

		// 获取字段值并加入到 map 中
		fv := value.Field(i)
		m[key] = fv.Interface()
	}

	return m, nil
}

type PersonInfo struct {
	UserName string `json:"person"`
	Time     string `json:"time"`
	WeekRate string `json:"week_rate"`
}

func (m *PersonList) Tag() string {
	return "person_list"
}
func (m *PersonList) MarshalJSON() ([]byte, error) {
	data, err := StructToMap(m)
	if err != nil {
		return nil, err
	}
	data["tag"] = m.Tag()
	return json.Marshal(data)
}

func (m *PersonList) Text() string {
	var result string
	for _, person := range m.Persons {
		result += person.Id + " "
	}
	return strings.TrimSpace(result)
}

type PersonList struct {
	Persons    []Persons `json:"persons"`
	Size       string    `json:"size"`
	Lines      int       `json:"lines"`
	ShowAvatar bool      `json:"show_avatar"`
	ShowName   bool      `json:"show_name"`
}

type Persons struct {
	Id string `json:"id"`
}

func (c *ChartCard) Tag() string {
	return "chart"
}
func (c *ChartCard) MarshalJSON() ([]byte, error) {
	data, err := StructToMap(c)
	if err != nil {
		return nil, err
	}
	data["tag"] = c.Tag()
	return json.Marshal(data)
}

type ChartCard struct {
	ChartSpec struct {
		Type  string `json:"type"`
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Data struct {
			Values interface{} `json:"values"`
		} `json:"data"`
		XField string `json:"xField"`
		YField string `json:"yField"`
	} `json:"chart_spec"`
}

type ChartData struct {
	Time  string `json:"time"`
	Value int    `json:"value"`
}

func (a *TableCard) Tag() string {
	return "table"
}
func (a *TableCard) MarshalJSON() ([]byte, error) {
	data, err := StructToMap(a)
	if err != nil {
		return nil, err
	}
	data["tag"] = a.Tag()
	return json.Marshal(data)
}

type TableCard struct {
	PageSize    int           `json:"page_size"`
	RowHeight   string        `json:"row_height"`
	HeaderStyle HeaderStyle   `json:"header_style"`
	Columns     []TableColumn `json:"columns"`
	Rows        interface{}   `json:"rows"`
}
type HeaderStyle struct {
	Bold            bool   `json:"bold"`
	BackgroundStyle string `json:"background_style"`
	Lines           int    `json:"lines"`
	TextSize        string `json:"text_size"`
	TextAlign       string `json:"text_align"`
}

type TableColumn struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	DataType    string `json:"data_type"`
	Width       string `json:"width,omitempty"`
	Format      struct {
		Symbol    string `json:"symbol"`
		Precision int    `json:"precision"`
	} `json:"format,omitempty"`
}

type Customer struct {
	CustomerName  string              `json:"customer_name"`
	CustomerScale []map[string]string `json:"customer_scale"`
	CustomerArr   float64             `json:"customer_arr"`
	CustomerYear  string              `json:"customer_year"`
}
