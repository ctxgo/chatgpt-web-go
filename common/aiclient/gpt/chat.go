package gpt

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"chatgpt-web-new-go/pkgs/limitdata"
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

type Chat struct {
	client            *Client
	enableChatContext bool
	history           limitdata.LimitLister[openai.ChatCompletionMessage]
	maxHistory        int
	sync.RWMutex
}

type ChatOpt func(*Chat)

func WithMaxHistory(h int) ChatOpt {
	return func(c *Chat) {
		c.maxHistory = h
	}
}

func NewChat(opts ...ChatOpt) *Chat {

	// Create a session with default values
	chat := &Chat{
		enableChatContext: true,
		maxHistory:        defaultMaxHistory,
	}

	for _, opt := range opts {
		opt(chat)
	}

	if chat.maxHistory == 0 || chat.maxHistory > defaultMaxHistory {
		chat.maxHistory = defaultMaxHistory
	}

	chat.history = limitdata.NewLimitList(
		limitdata.SetMaxLenth[openai.ChatCompletionMessage](chat.maxHistory),
		limitdata.SetStep[openai.ChatCompletionMessage](2))

	return chat
}

func (s *Chat) SetMaxHistory(maxHistory int) {
	s.Lock()
	defer s.Unlock()
	s.maxHistory = maxHistory
}

func (s *Chat) GetMaxHistory() int {
	s.RLock()
	defer s.RUnlock()
	return s.maxHistory
}

func (s *Chat) SetEnableChatContext(en bool) {
	s.Lock()
	defer s.Unlock()
	s.enableChatContext = en
}

func (s *Chat) GetEnableChatContext() bool {
	s.RLock()
	defer s.RUnlock()
	return s.enableChatContext
}

func (s *Chat) SetClient(c types.IClient) error {
	s.Lock()
	defer s.Unlock()
	if client, ok := c.(*Client); ok {
		s.client = client
		return nil
	}
	return errors.Errorf("unkonw client type: %v", c)
}

func (s *Chat) GetClient() types.IClient {
	s.RLock()
	defer s.RUnlock()
	return s.client
}

func (s *Chat) GetChatType() string {
	return "openai"
}

func (s *Chat) loadMessage(message openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	defer s.AddHistory(message.Content, openai.ChatMessageRoleUser)

	if s.enableChatContext {
		return append(s.history.Get(), message)
	}
	return []openai.ChatCompletionMessage{message}
}

func getMessageSender(msg types.IMessage) (sender MessageSender, err error) {
	switch m := msg.(type) {
	case types.TextMessage:
		sender = TextMessageSender(m)
	case types.ImageMessage:
		sender = ImageMessageSender(m)
	default:
		return nil, errors.Errorf("unknow message type, meg: %v", msg)
	}
	return
}

// 发送消息并接收数据流
func (s *Chat) SendStreamMessage(ctx context.Context, msg types.IMessage) (<-chan types.StreamData, error) {
	s.Lock()
	defer s.Unlock()
	sender, err := getMessageSender(msg)
	if err != nil {
		return nil, err
	}
	streamData := make(chan types.StreamData)
	data, err := sender.SendStream(ctx, s)
	if err != nil {
		return nil, err
	}
	go func() {
		s.Lock()
		defer s.Unlock()
		defer close(streamData)
		var content string
		for d := range data {
			streamData <- d
			if d.Err == nil {
				content += d.Data
			}
		}
		s.addHistory(msg.GetPrompt(), openai.ChatMessageRoleUser)
		s.addHistory(content, openai.ChatMessageRoleAssistant)
	}()
	return streamData, nil

}

// 发送消息并接收文本
func (s *Chat) SendMessage(ctx context.Context, msg types.IMessage) (string, error) {
	s.Lock()
	defer s.Unlock()
	sender, err := getMessageSender(msg)
	if err != nil {
		return "", err
	}
	return sender.Send(ctx, s)
}

func (s *Chat) CreateImage(ctx context.Context, msg types.IMessage) (*types.ImageData, error) {
	s.Lock()
	defer s.Unlock()
	m, ok := msg.(types.ImageCreater)
	if !ok {
		return nil, errors.Errorf("Unsupported request type")
	}
	return ImageCreater(m).Create(ctx, s)
}

func (s *Chat) AddHistory(message string, role string) {
	s.Lock()
	defer s.Unlock()
	s.addHistory(message, role)
}

func (s *Chat) addHistory(message string, role string) {
	if !s.enableChatContext || message == "" || role == "" {
		return
	}
	switch role {
	case types.MessageRoleModel:
		role = openai.ChatMessageRoleAssistant
	case types.MessageRoleUser:
		role = openai.ChatMessageRoleUser
	}
	s.history.Add(openai.ChatCompletionMessage{
		Content: message,
		Role:    role,
	})
}
