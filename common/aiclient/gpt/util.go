package gpt

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

// encodeToBase64URL 将二进制数据编码为 Base64 URL 安全字符串
func encodeToBase64URL(data []byte) string {
	// 检测数据类型
	mimeType := http.DetectContentType(data)

	// 转换为 Base64
	base64Encoded := base64.StdEncoding.EncodeToString(data)

	// 返回为 data URI 格式
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Encoded)
}
