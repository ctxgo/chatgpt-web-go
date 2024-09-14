package aiclient

import (
	"chatgpt-web-new-go/common/aiclient/types"
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/dao"
	"chatgpt-web-new-go/model"
	"sort"
	"strconv"
)

var (
	aiClient IClientManager
)

func InitClients(opts ...Opt) error {
	aiClient = NewClientManager(opts...)
	return initClients()
}

func initClients() error {
	dk := dao.Q.Aikey
	aiKeys, err := dk.Where(dk.IsDelete.Eq(0), dk.Status.Eq(1)).Find()
	if err != nil {
		return err
	}

	for _, key := range aiKeys {
		err = CreateClient(key)
		if err != nil {
			RevokeAikey(key)
			logs.Logger.Error(err)
			continue
		}
	}
	return nil
}

func GetChat(conversationID string, aiType string) (chat types.IChat, err error) {
	created, chat, err := aiClient.CreateChat(conversationID, aiType)
	if err != nil {
		return nil, err
	}
	if created && chat.GetEnableChatContext() {
		dm := dao.Q.Message
		data, err := dm.Where(dm.ParentMessageID.Eq(conversationID), dm.IsDelete.Eq(0)).
			Order(dm.CreateTime.Desc()).
			Limit(chat.GetMaxHistory()).
			Find()
		if err != nil {
			return nil, err
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i].CreateTime.Before(data[j].CreateTime)
		})
		for _, m := range data {
			chat.AddHistory(m.Content, m.Role)
		}

	}
	return chat, nil
}

func DelChat(conversationID string) {
	aiClient.DelChat(conversationID)
}

func RevokeAikey(key *model.Aikey) {
	dk := dao.Q.Aikey
	resultInfo, err := dk.Where(dk.ID.Eq(key.ID)).Update(dk.Status, 0)
	if err != nil {
		return
	}
	if resultInfo.RowsAffected < 1 {
		logs.Error("aiKey update fail: RowsAffected < 1")
		return
	}
}

func CreateClient(key *model.Aikey) (err error) {
	defer func() {
		if err != nil {
			RevokeAikey(key)
		}
	}()
	return aiClient.CreateClient(buildHost(key))
}

func UpdateClient(key *model.Aikey) (err error) {
	defer func() {
		if err != nil {
			RevokeAikey(key)
		}
	}()
	h := buildHost(key)
	if key.Status != 0 {
		return aiClient.UpdateClient(h)
	}
	aiClient.DeleteClient(h)
	return nil
}

func DeleteClient(key *model.Aikey) {
	aiClient.DeleteClient(buildHost(key))
}

func buildHost(key *model.Aikey) Host {
	return Host{
		Id:     strconv.FormatInt(key.ID, 10),
		Key:    key.Key,
		Host:   key.Host,
		AiType: key.Type,
	}
}
