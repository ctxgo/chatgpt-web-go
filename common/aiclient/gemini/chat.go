package gemini

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"chatgpt-web-new-go/pkgs/limitdata"
	"context"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// Default values
const (
	defaultMaxHistory = 50
)

type Chat struct {
	client            *Client
	model             *genai.GenerativeModel
	session           *genai.ChatSession
	enableChatContext bool
	history           limitdata.LimitLister[*genai.Content]
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
	chat := &Chat{
		enableChatContext: true,
	}

	for _, opt := range opts {
		opt(chat)
	}

	if chat.maxHistory == 0 || chat.maxHistory > defaultMaxHistory {
		chat.maxHistory = defaultMaxHistory
	}

	chat.history = limitdata.NewLimitList(
		limitdata.SetMaxLenth[*genai.Content](chat.maxHistory),
		limitdata.SetStep[*genai.Content](2))

	return chat
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

func (s *Chat) SetMaxHistory(h int) {
	s.Lock()
	defer s.Unlock()
	s.maxHistory = h
}
func (s *Chat) GetMaxHistory() int {
	s.RLock()
	defer s.RUnlock()
	return s.maxHistory
}

func (s *Chat) SetClient(c types.IClient) error {
	s.Lock()
	defer s.Unlock()
	if _client, ok := c.(*Client); ok {
		s.client = _client
		return nil
	}
	return errors.Errorf("unkonw client type: %v", c)
}

func (s *Chat) GetClient() types.IClient {
	s.RLock()
	defer s.RUnlock()
	return s.client
}

func (s *Chat) AddHistory(message string, role string) {
	s.Lock()
	defer s.Unlock()
	s.addHistory(message, role)
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
	s.setupSession(msg)
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
		s.addHistory(msg.GetPrompt(), "user")
		s.addHistory(content, "model")
	}()
	return streamData, nil
}

// 发送消息并接收文本
func (s *Chat) SendMessage(ctx context.Context, msg types.IMessage) (string, error) {
	s.Lock()
	defer s.Unlock()
	s.setupSession(msg)
	sender, err := getMessageSender(msg)
	var raWerr *googleapi.Error
	if errors.As(err, raWerr) {
		return "", raWerr
	}
	content, err := sender.Send(ctx, s)
	if err != nil {
		return "", err
	}
	s.addHistory(content, "model")
	return content, nil
}

func (s *Chat) CreateImage(ctx context.Context, msg types.IMessage) (*types.ImageData, error) {
	return nil, errors.Errorf("unsupported")
}

func (s *Chat) setupSession(msg types.IMessage) {
	s.buildModel(msg)
	if s.enableChatContext {
		s.loadHistory()
	}
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

func (s *Chat) addHistory(message string, role string) {
	if !s.enableChatContext || message == "" || role == "" {
		return
	}
	switch role {
	case types.MessageRoleModel:
		role = "model"
	case types.MessageRoleUser:
		role = "user"
	}
	s.history.Add(&genai.Content{
		Parts: []genai.Part{genai.Text(message)},
		Role:  role,
	})
}

func (s *Chat) loadHistory() {
	if !s.enableChatContext {
		return
	}
	s.session.History = s.history.Get()
}

func (c *Chat) buildModel(msg types.IMessage) {
	opts := msg.GetOptions()
	c.model = newModelWithSafety(c.client.Client, opts.Model)
	c.session = c.model.StartChat()
	if opts.MaxTokens != nil {
		c.model.MaxOutputTokens = opts.MaxTokens
	}
}

// setupModelWithSafety 配置并返回带有安全设置的模型
// 参考 https://ai.google.dev/gemini-api/docs/safety-settings?hl=zh-cn
func newModelWithSafety(client *genai.Client, modelName string) *genai.GenerativeModel {
	model := client.GenerativeModel(modelName)
	model.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockNone},       // 骚扰内容
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockNone},       // 仇恨言论
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockNone}, // 性明显内容
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockNone}, // 危险内容
	}
	return model
}
