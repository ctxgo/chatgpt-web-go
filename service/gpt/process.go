package gpt

import (
	"chatgpt-web-new-go/common/aiclient"
	"chatgpt-web-new-go/common/bizError"
	"chatgpt-web-new-go/common/goUtil"
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/common/types"
	"chatgpt-web-new-go/service/persona"
	"context"

	aiTypes "chatgpt-web-new-go/common/aiclient/types"
	"chatgpt-web-new-go/dao"
	"chatgpt-web-new-go/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	coinPerMessage = -1
)

func Process(ctx *gin.Context, r *Request, uid int64) (stream <-chan aiTypes.StreamData, err error) {
	// Deduction of points
	du := dao.Q.User
	updateInfo, err := du.WithContext(ctx).Where(du.ID.Eq(uid)).UpdateSimple(du.Integral.Add(coinPerMessage))
	if err != nil {
		logs.Error("user integral update error: %v", err)
		return
	}
	if updateInfo.RowsAffected < 1 {
		logs.Error("user integral update error: updateInfo.RowsAffected <= 0")
		err = bizError.IntegralNoneError
		return
	}

	// cnf.Model 是否在 chatModels 中
	chat, err := aiclient.GetChat(r.ParentMessageId, r.Options.AiType)
	if err != nil {
		logs.Logger.Errorf("get client err: %v!", err)
		return nil, bizError.AiKeyNoneUsefullError
	}
	var personaContext string
	if r.PersonaId != 0 {
		personaContext, err = getPersonaContext(ctx, r.PersonaId)
		if err != nil {
			logs.Logger.Errorf("get persona  err: %v!", err)
			return nil, bizError.PersonaGetError

		}
	}

	stream, err = chat.SendStreamMessage(ctx, aiTypes.TextMessage{
		Options: aiTypes.Options{
			Model:            r.Options.Model,
			MaxTokens:        r.Options.MaxTokens,
			Temperature:      r.Options.Temperature,
			PresencePenalty:  r.Options.PresencePenalty,
			FrequencyPenalty: r.Options.FrequencyPenalty,
			AiInstruction:    personaContext,
		},
		Prompt: r.Prompt,
	})
	if err != nil {
		logs.Error("aiClient client.CreateChatCompletion bizError:%v", err)
		aiKey := chat.GetClient().GetClientKey()
		goUtil.New(func() {
			refreshKey(aiKey)
		})
		return nil, err
	}

	// insert message
	goUtil.New(func() {
		addMessage(ctx, r, uid)
	})

	// insert record
	goUtil.New(func() {
		insertTurnOverRecord(ctx, r, uid)
	})

	return
}

func refreshKey(aiKey string) {

	da := dao.Q.Aikey
	resultInfo, err := da.Where(da.Key.Eq(aiKey)).Update(da.Status, 0)
	if err != nil {
		logs.Error("aiKey update error: %v", err)
		return
	}
	if resultInfo.RowsAffected < 1 {
		logs.Error("aiKey update fail: RowsAffected < 1")
		return
	}

}

func insertTurnOverRecord(ctx *gin.Context, r *Request, uid int64) {
	t := &model.Turnover{
		UserID:   uid,
		Describe: "对话(" + r.Options.Model + ")",
		Value:    fmt.Sprintf("%v积分", coinPerMessage),
	}

	dr := dao.Q.Turnover
	err := dr.WithContext(ctx).Create(t)
	if err != nil {
		logs.Error("turnover create error: %v", err)
	}
}

func addMessage(ctx *gin.Context, r *Request, uid int64) {
	msg := &model.Message{
		UserID:           uid,
		Content:          r.Prompt,
		PersonaID:        types.InterfaceToInt64(r.PersonaId),
		Role:             aiTypes.MessageRoleUser,
		FrequencyPenalty: int32(r.Options.FrequencyPenalty),
		MaxTokens:        *r.Options.MaxTokens,
		Model:            r.Options.Model,
		PresencePenalty:  int32(r.Options.PresencePenalty),
		Temperature:      int32(r.Options.Temperature),
		ParentMessageID:  r.ParentMessageId,
	}

	dm := dao.Q.Message
	err := dm.WithContext(ctx).Create(msg)
	if err != nil {
		logs.Error("message create error: %v", err)
	}
}

func getPersonaContext(ctx context.Context, id int64) (string, error) {
	persona, err := persona.PersonaInfo(ctx, id)
	if err != nil {
		return "", err
	}
	return persona.Context, nil
}
