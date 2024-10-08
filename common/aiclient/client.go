package aiclient

import (
	"chatgpt-web-new-go/common/aiclient/gemini"
	"chatgpt-web-new-go/common/aiclient/gpt"
	"chatgpt-web-new-go/common/aiclient/types"
	"chatgpt-web-new-go/common/config"
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// 会话有效期15分钟
const defaultChatExpiration = 15 * time.Minute

// 每隔16分钟清理过期会话
const defaultChatCleanupInterval = 16 * time.Minute

type Host struct {
	Id     string
	AiType string
	Host   string
	Key    string
}

type Client struct {
	client types.IClient
	host   Host
}

type IClientManager interface {
	CreateClient(h Host) error
	DeleteClient(h Host)
	UpdateClient(h Host) error
	CreateChat(conversationID, aiType string) (bool, types.IChat, error)
	DelChat(conversationID string)
}

type ClientManager struct {
	chatExpiration      time.Duration
	chatCleanupInterval time.Duration

	scheduleType scheduleType
	Scheduler    *SchedulerFactory

	clients map[string][]*Client
	sync.RWMutex

	chatManager ChatManager
}

type Opt func(*ClientManager)

func WithScheduleType(t scheduleType) Opt {
	return func(cm *ClientManager) {
		cm.scheduleType = t
	}
}

func WhithChatExpiration(t time.Duration) Opt {
	return func(cm *ClientManager) {
		cm.chatExpiration = t
	}
}

func WhithChatCleanupInterval(t time.Duration) Opt {
	return func(cm *ClientManager) {
		cm.chatCleanupInterval = t
	}
}

func NewClientManager(opts ...Opt) *ClientManager {
	c := &ClientManager{
		Scheduler:           NewSchedulerFactory(),
		chatCleanupInterval: defaultChatCleanupInterval,
		chatExpiration:      defaultChatExpiration,
		clients:             make(map[string][]*Client),
	}
	for _, o := range opts {
		o(c)
	}
	c.chatManager = *NewChatManager(c.chatExpiration, c.chatCleanupInterval)
	return c
}

func createIClient(h Host) (client *Client, err error) {
	var ilient types.IClient
	switch h.AiType {
	case "gemini":
		ilient, err = gemini.NewClient(context.Background(),
			h.Key, gemini.WithProxy(config.Config.AiClient.Proxy))

	case "openai":
		ilient, err = gpt.NewClient(context.Background(), h.Key,
			gpt.WithProxy(config.Config.AiClient.Proxy))

	default:
		return nil, errors.Errorf("unsupported AI type: %s", h.AiType)
	}
	if err != nil {
		return nil, err

	}
	return &Client{client: ilient, host: h}, err
}

func (c *ClientManager) CreateClient(h Host) error {
	c.Lock()
	defer c.Unlock()
	return c.createClient(h)
}

func (c *ClientManager) createClient(h Host) error {
	client, err := createIClient(h)
	if err != nil {
		return err
	}
	c.clients[h.AiType] = append(c.clients[h.AiType], client)
	return nil
}

func (c *ClientManager) DeleteClient(h Host) {
	c.Lock()
	defer c.Unlock()
	var newData []*Client
	data := c.clients[h.AiType]
	for _, x := range data {
		if x.host.Id != h.Id {
			newData = append(data, x)
		}
	}
	c.clients[h.AiType] = newData
}

func (c *ClientManager) UpdateClient(h Host) (err error) {
	c.Lock()
	defer c.Unlock()
	var newData []*Client
	data := c.clients[h.AiType]
	if len(data) == 0 {
		return c.createClient(h)
	}
	for _, x := range data {
		if x.host.Id != h.Id {
			newData = append(data, x)
			continue
		}
		if x.host.Host == h.Host && x.host.Key == h.Key {
			newData = append(newData, x)
			continue
		}
		client, rawErr := createIClient(h)
		if rawErr != nil {
			err = rawErr
			continue
		}
		newData = append(newData, client)
	}
	c.clients[h.AiType] = newData
	return
}

func (c *ClientManager) selectIndex(count int, scheduler Scheduler) (int, error) {
	if count == 0 {
		return 0, errors.Errorf("no clients available")
	}
	return scheduler.NextIndex(count), nil
}

func (c *ClientManager) getOnceClient(k string) (types.IClient, error) {
	data := c.clients[k]

	scheduler := c.Scheduler.GetOrCreateScheduler(k, c.scheduleType)
	index, err := c.selectIndex(len(data), scheduler)
	if err != nil {
		return nil, err
	}
	return data[index].client, nil
}

func (c *ClientManager) DelChat(conversationID string) {
	c.chatManager.Del(conversationID)
}

// bool值为false说明chat已经存在
func (c *ClientManager) CreateChat(conversationID, aiType string) (bool, types.IChat, error) {
	c.RLock()
	chat, ok, err := c.getAndSetChat(conversationID, aiType)
	c.RUnlock()
	if ok {
		return false, chat, err
	}
	c.Lock()
	defer c.Unlock()
	chat, ok, err = c.getAndSetChat(conversationID, aiType)
	if ok {
		return false, chat, err
	}
	chat, err = c.createIChat(conversationID, aiType)
	if err != nil {
		return false, nil, err
	}
	return true, chat, c.setChat(chat, aiType)
}

// checkIfChatNeedsCreation 辅助函数，判断是否需要创建新的 chat
func (c *ClientManager) getAndSetChat(conversationID, aiType string) (types.IChat, bool, error) {
	chat, ok := c.chatManager.Get(conversationID)
	if !ok || chat.GetChatType() != aiType {
		return nil, false, nil
	}

	return chat, true, c.setChat(chat, aiType)
}

func (c *ClientManager) createIChat(conversationID, aiType string) (types.IChat, error) {
	var chat types.IChat
	switch aiType {
	case "gemini":
		chat = gemini.NewChat()
	case "openai":
		chat = gpt.NewChat()
	default:
		return nil, errors.Errorf("unknow client type: %v", aiType)
	}
	c.chatManager.Set(conversationID, chat)
	return chat, nil
}

func (c *ClientManager) setChat(ch types.IChat, clientType string) error {
	client, err := c.getOnceClient(clientType)
	if err != nil {
		return err
	}
	return ch.SetClient(client)
}
