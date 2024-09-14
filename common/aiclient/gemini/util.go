package gemini

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

// func IsGoogleapiError(err error) (error, bool) {
// 	var raWerr *googleapi.Error
// 	if errors.As(err, &raWerr) {
// 		return raWerr, true
// 	}
// 	return err, false
// }

func iterMesssage(iter *genai.GenerateContentResponseIterator,
	responseChan chan types.StreamData) error {
	var text string
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			responseChan <- types.StreamData{Err: err}
			return err
		}
		if resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != 0 {
			return errors.New(resp.PromptFeedback.BlockReason.String())
		}
		data := retrieveResponce(resp)
		responseChan <- types.StreamData{Data: data}
		text = text + data
	}
	return nil
}

func retrieveResponce(resp *genai.GenerateContentResponse) string {
	var builder strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Fprint(&builder, part)
			}
		}
	}
	return builder.String()
}

func getImageFormat(data []byte) (string, error) {
	mimeType := http.DetectContentType(data)
	if !strings.HasPrefix(mimeType, "image/") {
		return "", fmt.Errorf("unsupported MIME type: %s", mimeType)
	}

	// 从MIME类型中提取格式，例如从 "image/jpeg" 中提取 "jpeg"
	format := strings.TrimPrefix(mimeType, "image/")
	return format, nil
}
