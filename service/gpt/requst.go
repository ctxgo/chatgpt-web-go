package gpt

type Request struct {
	ParentMessageId string  `json:"parentMessageId"`
	PersonaId       int64   `json:"persona_id"`
	Prompt          string  `json:"prompt"`
	Options         Options `json:"options"`
}

type Options struct {
	AiType           string  `json:"ai_type"`
	FrequencyPenalty float32 `json:"frequencyPenalty"`
	MaxTokens        *int32  `json:"max_tokens"`
	Model            string  `json:"model"`
	PresencePenalty  float32 `json:"presence_penalty"`
	Temperature      float32 `json:"temperature"`
}
