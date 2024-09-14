package gemini

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"context"
	"io"

	"github.com/google/generative-ai-go/genai"
)

type MessageSender interface {
	Send(ctx context.Context, s *Chat) (string, error)
	SendStream(ctx context.Context, s *Chat) (<-chan types.StreamData, error)
}

func sendMessage(ctx context.Context, s *Chat, part ...genai.Part) (string, error) {
	resp, err := s.session.SendMessage(ctx, part...)
	if err != nil {
		return "", err
	}
	respStr := retrieveResponce(resp)
	return respStr, nil

}

func sendStreamMessage(ctx context.Context, s *Chat, part ...genai.Part) <-chan types.StreamData {
	responseChan := make(chan types.StreamData, 10)

	go func() {
		defer close(responseChan)
		// 发送消息并接收响应
		iter := s.session.SendMessageStream(ctx, part...)
		iterMesssage(iter, responseChan)
	}()

	return responseChan
}

// 文本消息
type TextMessageSender types.TextMessage

func (sender TextMessageSender) Send(ctx context.Context, s *Chat) (string, error) {
	return sendMessage(ctx, s, genai.Text(sender.Prompt))
}

func (sender TextMessageSender) SendStream(ctx context.Context, s *Chat) (<-chan types.StreamData, error) {
	return sendStreamMessage(ctx, s, genai.Text(sender.Prompt)), nil
}

// images消息
type ImageMessageSender types.ImageMessage

func (sender ImageMessageSender) buildPart() ([]genai.Part, error) {
	data, err := io.ReadAll(sender.Data)
	if err != nil {
		return nil, err
	}
	mimeType, err := getImageFormat(data)
	if err != nil {
		return nil, err
	}
	return []genai.Part{
		genai.ImageData(mimeType, data),
		genai.Text(sender.Prompt),
	}, nil
}

func (sender ImageMessageSender) Send(ctx context.Context, s *Chat) (string, error) {
	prompt, err := sender.buildPart()
	if err != nil {
		return "", err
	}
	return sendMessage(ctx, s, prompt...)

}

func (sender ImageMessageSender) SendStream(ctx context.Context, s *Chat) (<-chan types.StreamData, error) {

	prompt, err := sender.buildPart()
	if err != nil {
		return nil, err
	}
	return sendStreamMessage(ctx, s, prompt...), nil

}
