package utils

import (
	"encoding/xml"
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"
)

func CutPrefix(s, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		return strings.TrimPrefix(s, prefix), true
	}
	return s, false
}

func EitherCutPrefix(s string, prefix ...string) (string, bool) {
	// 任一前缀匹配则返回剩余部分
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return strings.TrimPrefix(s, p), true
		}
	}
	return s, false
}

// TrimEqual trim space and equal
func TrimEqual(s, prefix string) (string, bool) {
	if strings.TrimSpace(s) == prefix {
		return "", true
	}
	return s, false
}

func EitherTrimEqual(s string, prefix ...string) (string, bool) {
	// 任一前缀匹配则返回剩余部分
	for _, p := range prefix {
		if strings.TrimSpace(s) == p {
			return "", true
		}
	}
	return s, false
}

// ContainsSpecificContent 使用正则表达式检查文本内容是否包含特定内容
func ContainsSpecificContent(text string, pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(text), nil
}

// ContainsSpecificContentV2 使用正则表达式检查文本内容是否包含特定内容
func ContainsSpecificContentV2(text string, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(text)
}

// ContainsSpecificContentV3 使用正则表达式检查文本内容是否包含特定内容，并将特定内容替换为掉
func ContainsSpecificContentV3(text string, regex string) (bool, string) {
	pattern := regexp.MustCompile(regex)
	matches := pattern.FindAllString(text, -1)
	if len(matches) > 0 {
		return true, regexp.MustCompile(regex).ReplaceAllString(text, "")
	}
	return false, text
}

func isEmpty(str *string) bool {
	return str == nil || *str == ""
}

func AddHoursToTime(t string, hours int) string {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		log.Printf("转换时间报错： %v  错误信息：%v", t, err.Error())
		return time.Now().Format("2006-01-02 15:04:05")
	}
	return parsedTime.Add(time.Duration(hours) * time.Hour).Format("2006-01-02 15:04:05")
}

func IsEndTimeBeforeNow(endTimeStr string, nowTimeStr string) bool {
	layout := "2006-01-02 15:04:05"
	endTime, err1 := time.Parse(layout, endTimeStr)
	nowTime, err2 := time.Parse(layout, nowTimeStr)
	if err1 != nil || err2 != nil {
		fmt.Println("Error parsing time")
		return false
	}

	return endTime.Before(nowTime)
}

// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// H is a shortcut for map[string]any
type H map[string]any

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "xml",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func chooseData(custom, wildcard any) any {
	if custom != nil {
		return custom
	}
	if wildcard != nil {
		return wildcard
	}
	panic("negotiation config is invalid")
}

func parseAccept(acceptHeader string) []string {
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if i := strings.IndexByte(part, ';'); i > 0 {
			part = part[:i]
		}
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func nameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func GetSubstrings(s string, start int, end int) string {
	if start < 0 || end > len(s) || start > end {
		return ""
	}
	return s[start:end]
}

// Unquote 去除JSON字符串中的转义字符
func Unquote(s string) string {
	length := len(s)
	if length < 2 || s[0] != '"' || s[length-1] != '"' {
		return s
	}
	s = s[1 : length-1]
	s = strings.Replace(s, `\"`, `"`, -1)
	return s
}

func extractNames(text string) []string {
	// 使用正则表达式匹配中文姓名，假设姓名为2-3个汉字
	re := regexp.MustCompile(`[\\4e00-\\u9fa5]{2,}`)
	names := re.FindAllString(text, -1)
	return names
}
func getFirstAndLastDay(yearMonth string) (string, string) {
	t, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		panic(err)
	}
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC)
	return firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02")
}
