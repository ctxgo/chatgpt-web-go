package gemini

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"chatgpt-web-new-go/pkgs/httpclient"
	"context"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/option"
)

type Client struct {
	*genai.Client
	key           string
	proxyUrl      string
	geminiOptions []option.ClientOption
}

type ClientOpt func(*Client)

func WithProxy(url string) ClientOpt {
	return func(gc *Client) {
		gc.proxyUrl = url
	}
}

func WithGeminOption(opts ...option.ClientOption) ClientOpt {
	return func(gc *Client) {
		gc.geminiOptions = append(gc.geminiOptions, opts...)
	}
}

func NewClient(ctx context.Context, apiKey string, opts ...ClientOpt) (*Client, error) {
	c := &Client{}
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
		httpClient.Transport = &transport.APIKey{
			Transport: httpClient.Transport,
			Key:       apiKey,
		}
		c.geminiOptions = append(c.geminiOptions, option.WithHTTPClient(httpClient))
	}

	c.geminiOptions = append(c.geminiOptions, option.WithAPIKey(apiKey))
	client, err := genai.NewClient(ctx, c.geminiOptions...)
	if err != nil {
		return nil, err
	}
	c.key = apiKey
	c.Client = client
	return c, nil
}

func (c *Client) GetClient() types.IClient {
	return c
}

func (c *Client) GetClientKey() string {
	return c.key
}
