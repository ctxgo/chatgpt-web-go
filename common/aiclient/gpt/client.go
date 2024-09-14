package gpt

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"chatgpt-web-new-go/pkgs/httpclient"
	"context"
	"time"

	"github.com/sashabaranov/go-openai"
)

// Default values
const (
	defaultMaxHistory = 30
)

type Client struct {
	*openai.Client
	key      string
	proxyUrl string
	options  openai.ClientConfig
}

type opt func(*Client)

func WithProxy(u string) opt {
	return func(gc *Client) {
		gc.proxyUrl = u
	}
}

func NewClient(ctx context.Context, apiKey string, opts ...opt) (*Client, error) {
	c := &Client{
		options: openai.DefaultConfig(apiKey),
	}
	for _, o := range opts {
		o(c)
	}
	if c.proxyUrl != "" {
		httpClient, err := httpclient.NewHttpClient(
			httpclient.SetProxy(c.proxyUrl),
			httpclient.SetTimeout(300*time.Second),
		)
		if err != nil {
			return nil, err
		}
		c.options.HTTPClient = httpClient
	}
	c.Client = openai.NewClientWithConfig(c.options)
	c.key = apiKey
	return c, nil
}

func (c *Client) GetClient() types.IClient {
	return c
}

func (c *Client) GetClientKey() string {
	return c.key
}
