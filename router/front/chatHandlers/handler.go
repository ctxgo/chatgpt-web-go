package chatHandlers

import (
	aiTypes "chatgpt-web-new-go/common/aiclient/types"
	authGlobal "chatgpt-web-new-go/common/auth"
	"chatgpt-web-new-go/common/bizError"
	"chatgpt-web-new-go/common/goUtil"
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/common/types"
	"chatgpt-web-new-go/model"
	"chatgpt-web-new-go/router/base"
	"chatgpt-web-new-go/service/gpt"
	"chatgpt-web-new-go/service/message"
	"encoding/json"
	"regexp"

	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ChatListHandler(c *gin.Context) {

}

func ChatCompletionsHandler(c *gin.Context) {
	var request gpt.Request
	err := c.BindJSON(&request)
	if err != nil {
		base.Fail(c, err.Error())
		return
	}
	logs.Info("chat process request: %v", request)

	if len(request.Prompt) == 0 {
		base.Fail(c, "request messages required")
		return
	}

	uFromCtx, found := c.Get(authGlobal.GinCtxKey)
	if !found {
		logs.Error("cannot get auth user key:%v", authGlobal.GinCtxKey)
		base.Fail(c, "cannot get auth user key:"+authGlobal.GinCtxKey)
		return
	}

	u := uFromCtx.(*model.User)
	model, aiType, err := parseModel(request.Options.Model)
	if err != nil {
		logs.Logger.Error(err)
		base.Fail(c, bizError.AiKeyModeParesError.Error())
		return
	}
	request.Options.AiType = aiType
	request.Options.Model = model
	stream, err := gpt.Process(c, &request, u.ID)
	if err != nil {
		logs.Error("aiClient Process bizError: %v", err)
		base.Fail(c, "错误:"+err.Error())
		return
	}

	var msgList []*message.ChatProcessResponse
	c.Stream(func(w io.Writer) bool {
		if data, ok := <-stream; ok {
			msg := &message.ChatProcessResponse{
				Role:            aiTypes.MessageRoleModel,
				Segment:         message.SegmentText,
				DateTime:        time.Now().Format(types.TimeFormatDate),
				Content:         data.Data,
				ParentMessageID: request.ParentMessageId,
			}
			if data.Err != nil {
				logs.Error("aiClient Process bizError: %v", data.Err)
				msg.Content = data.Err.Error()
			}
			msgList = append(msgList, msg)

			marshal, _ := json.Marshal(msg)
			write, err := w.Write(marshal)
			if err != nil {
				logs.Error("w write error: %v, write: %v", err, write)
				return false
			}

			_, _ = w.Write([]byte("\n\n"))
			return true
		}
		return false
	})

	goUtil.New(func() {
		message.MessageAdd(c, u.ID, &request, msgList)
	})
}

// 正则表达式用于匹配 'model(aiType)' 格式
func parseModel(model string) (string, string, error) {
	re := regexp.MustCompile(`^(.*)\((.*)\)$`)

	matches := re.FindStringSubmatch(model)

	if len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	return "", "", errors.Errorf("模型格式无法解析，%v", model)
}
