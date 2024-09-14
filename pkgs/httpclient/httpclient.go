package httpclient

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// HttpClientOpts holds configurable options for the HTTP client.
type HttpClientOpts struct {
	Timeout        time.Duration
	KeepAlive      time.Duration
	ConnectTimeout time.Duration
	ProxyURL       string
}

// ClientOption defines a function type that modifies HttpClientOpts.
type ClientOption func(*HttpClientOpts)

// NewHttpClient creates and returns a new HTTP client with the given options.
func NewHttpClient(opts ...ClientOption) (*http.Client, error) {
	options := HttpClientOpts{
		KeepAlive:      30 * time.Second,
		ConnectTimeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(&options)
	}

	dialer := &net.Dialer{
		Timeout:   options.ConnectTimeout,
		KeepAlive: options.KeepAlive,
	}

	transport := getDefaultTransport()
	transport.DialContext = dialer.DialContext
	if err := configureProxy(transport, options.ProxyURL, dialer); err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: transport,
		Timeout:   options.Timeout, // Set the overall request timeout
	}, nil
}

func getDefaultTransport() *http.Transport {
	return &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

}

// configureProxy sets the proxy configuration for the HTTP client.
// It supports HTTP, HTTPS, and SOCKS5 proxies.
func configureProxy(transport *http.Transport, proxyURL string, dialer *net.Dialer) error {
	if proxyURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("parsing proxy URL failed: %w", err)
	}

	switch parsedURL.Scheme {
	case "http", "https":
		// For HTTP and HTTPS proxies, set the Proxy field in the transport.
		transport.Proxy = http.ProxyURL(parsedURL)
	case "socks5", "socks5h":
		// For SOCKS5 proxies, configure the DialContext with proxy dialer support.
		proxyDialer, err := proxy.FromURL(parsedURL, dialer)
		if err != nil {
			return fmt.Errorf("setting up socks proxy failed: %w", err)
		}
		contextDialer, ok := proxyDialer.(proxy.ContextDialer)
		if !ok {
			return errors.New("proxy does not support DialContext")
		}
		transport.DialContext = contextDialer.DialContext
	default:
		return fmt.Errorf("unsupported proxy scheme: %s", parsedURL.Scheme)
	}

	return nil
}

// SetTimeout sets the overall timeout for the HTTP client.
func SetTimeout(timeout time.Duration) ClientOption {
	return func(opts *HttpClientOpts) {
		opts.Timeout = timeout
	}
}

// SetKeepAlive sets the keep-alive duration for the HTTP client's connections.
func SetKeepAlive(keepAlive time.Duration) ClientOption {
	return func(opts *HttpClientOpts) {
		opts.KeepAlive = keepAlive
	}
}

// SetConnectTimeout sets the connection timeout for the HTTP client.
func SetConnectTimeout(connectTimeout time.Duration) ClientOption {
	return func(opts *HttpClientOpts) {
		opts.ConnectTimeout = connectTimeout
	}
}

// SetProxy sets the proxy URL for the HTTP client.
func SetProxy(proxyURL string) ClientOption {
	return func(opts *HttpClientOpts) {
		opts.ProxyURL = proxyURL
	}
}
