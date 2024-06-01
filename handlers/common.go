package handlers

import (
	"FeishuCozeRobot/model"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// func sendCard
func msgFilter(msg string) string {
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")

}

func parseContent(content string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["text"] == nil {
		return ""
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}

func ParseContent(content string) string {

	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["text"] == nil {
		return ""
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}
func processMessage(msg interface{}) (string, error) {
	msg = strings.TrimSpace(msg.(string))
	msgB, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	msgStr := string(msgB)

	if len(msgStr) >= 2 {
		msgStr = msgStr[1 : len(msgStr)-1]
	}
	return msgStr, nil
}

func processNewLine(msg string) string {
	return strings.Replace(msg, "\\n", `
`, -1)
}

func processQuote(msg string) string {
	return strings.Replace(msg, "\\\"", "\"", -1)
}

// 将字符中 \u003c 替换为 <  等等
func processUnicode(msg string) string {
	regex := regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)
	return regex.ReplaceAllStringFunc(msg, func(s string) string {
		r, _ := regexp.Compile(`\\u`)
		s = r.ReplaceAllString(s, "")
		i, _ := strconv.ParseInt(s, 16, 32)
		return string(rune(i))
	})
}

func cleanTextBlock(msg string) string {
	msg = processNewLine(msg)
	msg = processUnicode(msg)
	msg = processQuote(msg)
	return msg
}

func parseFileKey(content string) model.File {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return model.File{}
	}
	if contentMap["file_key"] == nil && contentMap["file_name"] == nil {
		return model.File{}
	}
	return model.File{FileKey: contentMap["file_key"].(string), FileName: contentMap["file_name"].(string)}
}

func parseImageKey(content string) string {
	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if contentMap["image_key"] == nil {
		return ""
	}
	imageKey := contentMap["image_key"].(string)
	return imageKey
}

func getMessag(content string) (string, []model.File, []string) {
	var PostMessage model.PostMessage
	err := json.Unmarshal([]byte(content), &PostMessage)
	if err != nil {
		fmt.Println(err)
		return "", nil, nil
	}

	var qParsed string
	var file []model.File
	var imageKey []string
	for _, element := range PostMessage.Data {
		if element[0].Tag == "text" {

			Text := element[0].Text
			qParsed += Text
		} else if element[0].Tag == "img" {
			imageKey = append(imageKey, element[0].ImageKey)
		} else if element[0].Tag == "file" {
			file = append(file, model.File{
				FileKey:  element[0].FileKey,
				FileName: element[0].FileName,
			})
		}
	}

	return qParsed, file, imageKey
}
