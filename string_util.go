package siteadapt

import (
	"github.com/golang-module/carbon/v2"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// StrToByteSize 体积字符串转字节数
func StrToByteSize(text string) int64 {
	if text == "" {
		return 0
	}
	// 去除非法字符，并转换为大写
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ToUpper(text)
	// 匹配大小单位（如 KB、MB、GB 等）
	re := regexp.MustCompile(`[KMGTPI]*B?`)
	sizeStr := re.ReplaceAllString(text, "")
	// 尝试转换为浮点数
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		log.Printf("字符串转 float64 异常: %v", err)
		return 0
	}
	// 根据单位乘以相应的倍数
	switch {
	case strings.Contains(text, "PB") || strings.Contains(text, "PIB") || strings.Contains(text, "P"):
		size *= 1024 * 1024 * 1024 * 1024 * 1024
	case strings.Contains(text, "TB") || strings.Contains(text, "TIB") || strings.Contains(text, "T"):
		size *= 1024 * 1024 * 1024 * 1024
	case strings.Contains(text, "GB") || strings.Contains(text, "GIB") || strings.Contains(text, "G"):
		size *= 1024 * 1024 * 1024
	case strings.Contains(text, "MB") || strings.Contains(text, "MIB") || strings.Contains(text, "M"):
		size *= 1024 * 1024
	case strings.Contains(text, "KB") || strings.Contains(text, "KIB") || strings.Contains(text, "K"):
		size *= 1024
	}
	// 四舍五入并返回整数值
	return int64(size + 0.5)
}

// ParseInt 将字符串转换为整数
func ParseInt(str string) int {
	if len(str) == 0 {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("字符串转 int 异常: %v", err)
		return 0
	}
	return value
}

func ParseInt64(str string) int64 {
	if len(str) == 0 {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Printf("字符串转 int64 异常: %v", err)
		return 0
	}
	return i
}

func ParseFloat64(str string) float64 {
	if len(str) == 0 {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Printf("字符串转 float64 异常: %v", err)
		return 0
	}
	return f
}

func ParseBool(str string) bool {
	if len(str) == 0 {
		return false
	}
	b, err := strconv.ParseBool(str)
	if err != nil {
		log.Printf("字符串转 bool 异常: %v", err)
		return false
	}
	return b
}

func IsValidHttpUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// GetTimeStamp 解析日期字符串并返回时间戳
func GetTimeStamp(date string) int64 {
	if len(date) == 0 {
		return 0
	}
	return carbon.Parse(date).Timestamp()
}
