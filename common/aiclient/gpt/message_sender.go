package gpt

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

type MessageSender interface {
	Send(ctx context.Context, s *Chat) (string, error)
	SendStream(ctx context.Context, s *Chat) (<-chan types.StreamData, error)
}

func sendMessage(ctx context.Context, s *Chat, r openai.ChatCompletionRequest) (string, error) {
	resp, err := s.client.CreateChatCompletion(ctx, r)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil

}

func sendStreamMessage(ctx context.Context, s *Chat, r openai.ChatCompletionRequest) (<-chan types.StreamData, error) {
	r.Stream = true
	responseChan := make(chan types.StreamData, 10)
	stream, err := s.client.CreateChatCompletionStream(ctx, r)
	if err != nil {
		return nil, errors.Errorf("ChatCompletionStream error: %v\n", err)
	}
	go func() {
		defer close(responseChan)
		var text string
		defer s.AddHistory(text, openai.ChatMessageRoleAssistant)
		for {
			resp, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				responseChan <- types.StreamData{Err: err}
				return
			}
			if len(resp.Choices) != 0 {
				content := resp.Choices[0].Delta.Content
				responseChan <- types.StreamData{Data: content}
				text = text + content
			}
		}
	}()
	return responseChan, nil
}

// 文本消息
type TextMessageSender types.TextMessage

func (sender TextMessageSender) buildRequest(s *Chat) openai.ChatCompletionRequest {
	messages := s.loadMessage(openai.ChatCompletionMessage{
		Content: sender.Prompt,
		Role:    sender.Options.Model,
	})

	req := openai.ChatCompletionRequest{
		Messages: messages,
	}
	if sender.Options.MaxTokens != nil {
		req.MaxTokens = int(*sender.Options.MaxTokens)
	}
	return req
}

func (sender TextMessageSender) Send(ctx context.Context, s *Chat) (string, error) {
	return sendMessage(ctx, s, sender.buildRequest(s))
}

func (sender TextMessageSender) SendStream(ctx context.Context, s *Chat) (<-chan types.StreamData, error) {
	return sendStreamMessage(ctx, s, sender.buildRequest(s))
}

// images消息
type ImageMessageSender types.ImageMessage

func (sender ImageMessageSender) buildRequest(s *Chat) (openai.ChatCompletionRequest, error) {
	var req openai.ChatCompletionRequest
	data, err := io.ReadAll(sender.Data)
	if err != nil {
		return req, err
	}
	imageBase64Url := encodeToBase64URL(data)
	messages := s.loadMessage(openai.ChatCompletionMessage{
		Content: sender.Prompt,
		Role:    sender.Options.Model,
		MultiContent: []openai.ChatMessagePart{
			{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL:    imageBase64Url,
					Detail: openai.ImageURLDetailAuto,
				},
			},
		},
	})
	req = openai.ChatCompletionRequest{
		Messages: messages,
	}
	if sender.Options.MaxTokens != nil {
		req.MaxTokens = int(*sender.Options.MaxTokens)
	}
	return req, nil
}

func (sender ImageMessageSender) Send(ctx context.Context, s *Chat) (string, error) {
	req, err := sender.buildRequest(s)
	if err != nil {
		return "", err
	}
	return sendMessage(ctx, s, req)
}

func (sender ImageMessageSender) SendStream(ctx context.Context, s *Chat) (<-chan types.StreamData, error) {
	req, err := sender.buildRequest(s)
	if err != nil {
		return nil, err
	}
	return sendStreamMessage(ctx, s, req)
}

type ImageCreater types.ImageCreater

func (sender ImageCreater) Create(ctx context.Context, s *Chat) (*types.ImageData, error) {
	reqUrl := openai.ImageRequest{
		Prompt:         sender.Prompt,
		Size:           string(sender.Size),
		ResponseFormat: string(sender.ResponseFormat),
		N:              1,
	}

	resp, err := s.client.CreateImage(ctx, reqUrl)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, errors.Errorf("empty data")
	}
	var data types.ImageData
	switch sender.ResponseFormat {
	case types.CreateImageResponseFormatURL:
		data.URL = resp.Data[0].URL
	case types.CreateImageResponseFormatB64JSON:
		data.B64JSON = resp.Data[0].B64JSON
	}
	return &data, nil
}
