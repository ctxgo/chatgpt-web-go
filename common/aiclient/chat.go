package aiclient

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"time"

	"github.com/patrickmn/go-cache"
)

type ChatManager struct {
	*cache.Cache
}

type timeOpts func(*time.Duration)

func SetTimeOut(d time.Duration) timeOpts {
	return func(_d *time.Duration) {
		*_d = d
	}

}

func NewChatManager(expiration time.Duration, cleanupInterval time.Duration) *ChatManager {
	return &ChatManager{
		Cache: cache.New(expiration, cleanupInterval),
	}
}

func (m *ChatManager) Set(key string, chat types.IChat, timeOpts ...timeOpts) {
	timeOut := cache.DefaultExpiration
	for _, t := range timeOpts {
		t(&timeOut)
	}
	m.Cache.Set(key, chat, timeOut)
}

func (m *ChatManager) Get(key string) (types.IChat, bool) {
	data, found := m.Cache.Get(key)
	if !found {
		return nil, false
	}

	ChatSession, ok := data.(types.IChat)
	if !ok {
		return nil, false
	}
	return ChatSession, true
}

func (m *ChatManager) Del(key string) {
	m.Cache.Delete(key)
}
