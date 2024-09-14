package gemini

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"context"
	"fmt"
	"os"
	"testing"
)

func TestNewGeminiClient(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				apiKey: "xxxxxxxxxxxxxxxxxxxxxxxxx",
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(context.Background(),
				tt.args.apiKey,
				WithProxy("http://xxxxxxxxx"),
			)
			if err != nil {
				panic(err)
			}
			s := NewChat()
			s.SetClient(got)
			f, err := os.Open("/test.png")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			resp, err := s.SendMessage(context.Background(), types.TextMessage{
				Options: types.Options{
					Model: "gemini-1.5-pro",
				},
				Prompt: `后续的每次回复请在回复之前说 "hi", 如 "hi,[你的回复]"`,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)
			resp, err = s.SendMessage(context.Background(), types.TextMessage{
				Options: types.Options{
					Model: "gemini-1.5-pro",
				},
				Prompt: `请随机说出一个成语`,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)

			resp, err = s.SendMessage(context.Background(), types.ImageMessage{
				TextMessage: types.TextMessage{
					Options: types.Options{
						Model: "gemini-1.5-pro",
					},
					Prompt: "描述图片中的内容",
				},
				Data: f,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)
		})
	}
}
