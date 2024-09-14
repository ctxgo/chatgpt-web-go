package types

import (
	"context"
)

type IClient interface {
	GetClient() IClient
	GetClientKey() string
}

type IChat interface {
	SetClient(c IClient) error
	GetClient() IClient
	SetMaxHistory(int)
	GetMaxHistory() int
	AddHistory(msg, role string)
	SetEnableChatContext(en bool)
	GetEnableChatContext() bool
	SendStreamMessage(ctx context.Context, msg IMessage) (<-chan StreamData, error)
	SendMessage(ctx context.Context, msg IMessage) (string, error)
	CreateImage(ctx context.Context, msg IMessage) (*ImageData, error)
}
