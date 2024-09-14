package user

import (
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/dao"
	"chatgpt-web-new-go/model"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUserRecord(ctx *gin.Context, uid int64, r *Request) (d any, c int64, e error) {
	switch r.Type {
	case "invitation_records":
		return GetUserInviteRecordRecord(ctx, uid, r.Page.Page, r.PageSize)
	default:
		return nil, 0, errors.Errorf("暂不支持")
	}
}

func GetUserInviteRecordRecord(ctx *gin.Context, uid int64, page, size int) (d []*model.InviteRecord, c int64, e error) {
	dt := dao.Q.InviteRecord
	recorDao := dt.WithContext(ctx)

	d, c, e = recorDao.Where(dt.UserID.Eq(uid), dt.IsDelete.Eq(0)).
		FindByPage((page-1)*size, size)
	if e != nil {
		logs.Error("userDao get error: %v", e)
		return
	}

	return
}
