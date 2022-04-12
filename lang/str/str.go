package str

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode"
)

// IsNotBlank 判断字符串是否不是空白
//
// @return bool
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

// IsBlank 判断字符串是否是空白
//
// @return bool
func IsBlank(str string) bool {

	var strLen = len(str)
	if strLen != 0 {
		for _, r := range str {
			if !unicode.IsSpace(r) {
				return false
			}
		}
		return true
	} else {
		return true
	}
}

// IsEmpty 判断字符串是否是空
//
// @return bool
func IsEmpty(str string) bool {
	return len(str) == 0
}

// IsNotEmpty  判断字符串是否不是空
//
// @return bool
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// 去除字符串两边空格
//
// @return bool
func Trim(str string) string {
	var trim = strings.TrimSpace(str)
	return trim
}

// 判断 字符串 是否是以 substr 开头
func StartsWith(s string, substr string) bool {
	if strings.Index(s, substr) == 0 {
		return true
	} else {
		return false
	}
}

// 判断 字符串 是否是以 substr 结尾
func EndsWith(s string, substr string) bool {
	if strings.Index(s, substr) == len(s)-len(substr) {
		return true
	} else {
		return false
	}
}

// Strval 获取变量的字符串值
//
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
//
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
