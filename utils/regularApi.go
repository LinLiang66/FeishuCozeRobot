package utils

import (
	"FeishuCozeRobot/model"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// SplitList 列表分割
func SplitList(list []string, size int) [][]string {
	result := make([][]string, 0)
	if list == nil || len(list) == 0 || size <= 0 {
		return result
	}
	totalSize := len(list)
	count := totalSize / size
	remainder := totalSize % size
	for i := 0; i < size; i++ {
		fromIndex := i*count + min(i, remainder)
		toIndex := fromIndex + count
		if i < remainder {
			toIndex++
		}
		result = append(result, list[fromIndex:toIndex])
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ContentMatch  使用正则表达式检查文本内容是否包含特定内容
func ContentMatch(text string, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(text)
}

func ChecklenCoze(Messages []model.ChatHistory) []model.ChatHistory {
	for getCozeLength(Messages) > 50 {
		Messages = Messages[1:]
	}
	return Messages
}
func getCozeLength(Messages []model.ChatHistory) int {
	return len(Messages)
}

func GetSubstring(text string, maxLength int) string {
	if utf8.RuneCountInString(text) <= maxLength {
		return text
	}
	re := regexp.MustCompile("(.{" + strconv.Itoa(maxLength) + "})") //获取特定长度的内容
	return re.FindStringSubmatch(text)[1]
}

func IsDateTimeGreaterThan(dateTime1 string, dateTime2 string) bool {
	layout := "2006-01-02 15:04:05"
	t1, err1 := time.Parse(layout, dateTime1)
	t2, err2 := time.Parse(layout, dateTime2)
	if err1 != nil || err2 != nil {
		return false
	}
	return t1.After(t2)
}

func GetCauseContent(text string) string {
	// 正则表达式匹配常见中文字符范围，包括基本多文种平面(BMP)的常用汉字
	pattern := "[^\\x00-\\x7F]+"
	// 使用正则表达式查找所有匹配项
	r := regexp.MustCompile(pattern)
	matches := r.FindAllString(text, -1)
	// 把所有匹配到的中文字符连接成一个新的字符串
	return strings.Join(matches, "")
}

func MatchingType(input string, regex string) (string, bool) {
	pattern := regexp.MustCompile(regex)
	matches := pattern.FindAllString(input, -1)
	if len(matches) > 0 {
		return matches[0], true
	}
	return "", false
}

func RemoveSpaces(s string) string {
	re := regexp.MustCompile("(?m)(?<!\\S)\\s(?!\\S)")
	return re.ReplaceAllString(s, "")
}
