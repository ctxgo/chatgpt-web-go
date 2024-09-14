package messageHandlers

import (
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/router/base"
	"chatgpt-web-new-go/service/message"

	"github.com/gin-gonic/gin"
)

func MessageList(c *gin.Context) {
	request := &base.Page{}
	if err := c.Bind(request); err != nil {
		logs.Error("json bind error: %v", err)
		base.Fail(c, "参数异常！"+err.Error())
		return
	}

	result, count, err := message.AdminMessageList(c, request.Page, request.PageSize)
	if err != nil {
		logs.Error("message list error: %v", err)
		base.Fail(c, "查询会话列表异常："+err.Error())
		return
	}

	base.Success(c, base.PageResponse{
		Count: count,
		Rows:  result,
	})

}
