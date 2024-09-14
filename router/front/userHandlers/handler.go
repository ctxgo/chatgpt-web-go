package userHandlers

import (
	authGlobal "chatgpt-web-new-go/common/auth"

	"chatgpt-web-new-go/common/auth"
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/model"
	"chatgpt-web-new-go/router/base"
	"chatgpt-web-new-go/service/user"

	"github.com/gin-gonic/gin"
)

func UserInfoHandler(c *gin.Context) {
	uFromCtx, found := c.Get(auth.GinCtxKey)
	if !found {
		logs.Error("cannot get auth user key:%v", auth.GinCtxKey)
		base.Fail(c, "cannot get auth user key:"+auth.GinCtxKey)
		return
	}

	u := uFromCtx.(*model.User)

	userModel, err := user.GetUserInfo(c, u.Account)
	if err != nil {
		logs.Error("user GetUserInfo bizError: %v", err)
		base.Fail(c, "user GetUserInfo bizError")
		return
	}
	base.Success(c, userModel)
}

func UserRecordHandler(c *gin.Context) {
	var request user.Request
	err := c.ShouldBindQuery(&request)
	if err != nil {
		base.Fail(c, err.Error())
		return
	}

	uFromCtx, found := c.Get(authGlobal.GinCtxKey)
	if !found {
		logs.Error("cannot get auth user key:%v", authGlobal.GinCtxKey)
		base.Fail(c, "cannot get auth user key:"+authGlobal.GinCtxKey)
		return
	}

	u := uFromCtx.(*model.User)

	data, count, err := user.GetUserRecord(c, u.ID, &request)
	if err != nil {
		logs.Error("message list error: %v", err)
		base.Fail(c, "查询记录列表异常："+err.Error())
		return
	}

	base.Success(c, base.PageResponse{
		Count: count,
		Rows:  data,
	})
}
