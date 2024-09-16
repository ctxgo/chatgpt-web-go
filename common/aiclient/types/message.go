package types

import (
	"io"

	"github.com/sashabaranov/go-openai"
)

const MessageRoleUser = "user"
const MessageRoleModel = "model"

type Options struct {
	AiType           string
	FrequencyPenalty float32
	MaxTokens        *int32
	Model            string
	PresencePenalty  float32
	Temperature      float32
	AiInstruction    string
}

type IMessage interface {
	GetPrompt() string
	GetOptions() Options
}

type TextMessage struct {
	Options Options
	Prompt  string
}

func (m TextMessage) GetPrompt() string {
	return m.Prompt
}
func (m TextMessage) GetOptions() Options {
	return m.Options
}

type StreamData struct {
	Data string
	Err  error
}

type ImageMessage struct {
	TextMessage
	Data io.Reader
}

type ImageSize string
type ResponseFormat string

type ImageCreater struct {
	TextMessage
	Size           ImageSize
	ResponseFormat ResponseFormat
}

type ImageData openai.ImageResponseDataInner

const (
	CreateImageResponseFormatURL     ResponseFormat = "url"
	CreateImageResponseFormatB64JSON ResponseFormat = "b64_json"
)

const (
	CreateImageSize256x256   ImageSize = "256x256"
	CreateImageSize512x512   ImageSize = "512x512"
	CreateImageSize1024x1024 ImageSize = "1024x1024"
	// dall-e-3 supported only.
	CreateImageSize1792x1024 ImageSize = "1792x1024"
	CreateImageSize1024x1792 ImageSize = "1024x1792"
)
