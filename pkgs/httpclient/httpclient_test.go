package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestNewHttpClient(t *testing.T) {
	type args struct {
		opts []ClientOption
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Client
		wantErr bool
	}{
		{}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHttpClient(
				SetTimeout(10*time.Microsecond),
				SetProxy("socks5://xxxxxxxxxxxxx"),
			)
			if err != nil {
				log.Fatalf("创建HttpClient失败: %v", err)

			}
			req, err := http.NewRequest("GET", "https://www.google.com", nil)
			if err != nil {
				log.Fatalf("创建请求失败: %v", err)
			}
			// 发送请求
			resp, err := got.Do(req)
			if err != nil {
				log.Fatalf("请求失败: %v", err)
			}
			defer resp.Body.Close()

			// 打印响应状态
			fmt.Printf("响应状态: %s\n", resp.Status)

			// 打印响应内容
			var body []byte
			if _, err := resp.Body.Read(body); err != nil {
				log.Fatalf("读取响应内容失败: %v", err)
			}
			fmt.Printf("响应内容: %s\n", body)
		})
	}
}
