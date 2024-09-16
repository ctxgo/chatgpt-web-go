package persona

import (
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/dao"
	"chatgpt-web-new-go/model"
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

var requestGroup singleflight.Group

func PersonaList(ctx context.Context) (personaList []*model.Persona, err error) {
	dp := dao.Q.Persona

	personaList, err = dp.WithContext(ctx).Find()
	if err != nil {
		logs.Error("persona find error: %v", err)
		return
	}
	return
}

func PersonaInfo(ctx context.Context, id int64) (*model.Persona, error) {
	key := generateCacheKey(id)
	persona, err := getPersonaFromCache(ctx, key)
	if err == nil {
		return persona, err
	}
	if err != redis.Nil {
		return nil, errors.Errorf("redis err: %v", err)
	}
	// 使用 singleflight 防止缓存击穿
	result, err, _ := requestGroup.Do(key, func() (interface{}, error) {
		du := dao.Q.Persona
		resultInfo, err := du.WithContext(ctx).Where(du.ID.Eq(id)).First()
		if err != nil {
			return "", errors.WithStack(err)
		}
		err = setPersonaCache(ctx, key, resultInfo)
		return resultInfo, errors.WithStack(err)
	})

	if err != nil {
		logs.Logger.Errorf("%+v", err)
		return nil, err
	}

	persona, ok := result.(*model.Persona)
	if !ok {
		return nil, errors.Errorf("type assertion to string failed")
	}

	return persona, nil
}
