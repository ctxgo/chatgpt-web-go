package message

type Request struct {
	ParentMessageId string `form:"parent_message_id"`
}
