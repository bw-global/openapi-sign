package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
)

// Signature: Generate signature
// method: Request method
// version: API version
// appkey: Unique identifier of the application
// appsecret: Application secret key
// scope: Permission scope
// scopeValue: Permission scope value
// body: Request body (supports JSON object)
// Return value: Generated signature string or error message
func Signature(method, version, appkey, appsecret, scope, scopeValue string, body interface{}) (string, error) {
	if isNull(method) || isNull(version) || isNull(appkey) || isNull(appsecret) || isNull(scope) || isNull(scopeValue) {
		return "", fmt.Errorf("method, version, appkey, appsecret, scope and scopeValue cannot be null or empty")
	}
	qmap := make(map[string]string, 6)
	qmap["method"] = method
	qmap["version"] = version
	qmap["appkey"] = appkey
	qmap["scope"] = scope
	qmap["scopeValue"] = scopeValue
	qmap["appsecret"] = appsecret
	return generateSignature(qmap, body)
}

func generateSignature(params map[string]string, body interface{}) (string, error) {
	// 获取所有参数键并排序
	var keys []string
	for k := range params {
		if !isNull(k) && !isNull(params[k]) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接参数
	var builder strings.Builder
	//builder.WriteString(secret)

	for _, k := range keys {
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(params[k])
		builder.WriteString("&")
	}
	if body != nil {
		bodyBytes, err := sortedJSONMarshal(body)
		if err != nil {
			return "", fmt.Errorf("marshal body error: %w", err)
		}
		builder.Write(bodyBytes)
	}

	fmt.Println("拼接字符串:", builder.String())
	// MD5加密
	hasher := md5.New()
	hasher.Write([]byte(builder.String()))
	hash := hasher.Sum(nil)

	// 转换为大写的十六进制
	return strings.ToUpper(hex.EncodeToString(hash)), nil
}

// 创建辅助函数处理 JSON 序列化
func sortedJSONMarshal(v interface{}) ([]byte, error) {
	// 先正常序列化再反序列化，确保数据格式统一
	b, err := sonic.Marshal(v)
	if err != nil {
		return nil, err
	}

	var parsed interface{}
	if err := sonic.Unmarshal(b, &parsed); err != nil {
		return nil, err
	}

	// 内部递归函数
	var sortValue func(interface{}) ([]byte, error)
	sortValue = func(val interface{}) ([]byte, error) {
		switch v := val.(type) {
		case nil:
			return []byte("null"), nil

		case bool:
			if v {
				return []byte("true"), nil
			}
			return []byte("false"), nil

		case float64: // JSON 数字统一为 float64
			return []byte(fmt.Sprintf("%g", v)), nil

		case string:
			return sonic.Marshal(v)

		case []interface{}:
			// 处理数组
			if len(v) == 0 {
				return []byte("[]"), nil
			}

			var result strings.Builder
			result.WriteByte('[')

			for i, item := range v {
				itemBytes, err := sortValue(item)
				if err != nil {
					return nil, err
				}
				result.Write(itemBytes)

				if i < len(v)-1 {
					result.WriteByte(',')
				}
			}

			result.WriteByte(']')
			return []byte(result.String()), nil

		case map[string]interface{}:
			// 处理对象
			if len(v) == 0 {
				return []byte("{}"), nil
			}

			// 获取所有键并排序
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			var result strings.Builder
			result.WriteByte('{')

			for i, k := range keys {
				// 序列化键
				key, err := sonic.Marshal(k)
				if err != nil {
					return nil, err
				}
				result.Write(key)
				result.WriteByte(':')

				// 递归序列化值
				valueBytes, err := sortValue(v[k])
				if err != nil {
					return nil, err
				}
				result.Write(valueBytes)

				if i < len(keys)-1 {
					result.WriteByte(',')
				}
			}

			result.WriteByte('}')
			return []byte(result.String()), nil

		default:
			// 其他类型直接序列化
			return sonic.Marshal(v)
		}
	}

	return sortValue(parsed)
}

func isNull(str string) bool {
	return str == "" || len(strings.TrimSpace(str)) == 0
}
