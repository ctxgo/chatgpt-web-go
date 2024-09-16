package persona

import (
	"chatgpt-web-new-go/common/config"
	"chatgpt-web-new-go/common/logs"
	"chatgpt-web-new-go/model"
	"chatgpt-web-new-go/pkgs/stringx"
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

func generateCacheKey(id int64) string {
	return "persona-" + strconv.FormatInt(id, 10)
}

func getPersonaFromCache(ctx context.Context, key string) (*model.Persona, error) {
	data, err := config.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var persona model.Persona
	err = json.Unmarshal([]byte(data), &persona)
	if err != nil {
		return nil, errors.Errorf("JSON decoding failed: %v", err)
	}
	return &persona, nil
}

func setPersonaCache(ctx context.Context, key string, v *model.Persona) (err error) {
	defer func() {
		if err != nil {
			logs.Error(err.Error())
		}
	}()
	data, err := stringx.ConvertToString(v)
	if err != nil {
		return err
	}
	err = config.Redis.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func delPersonaCache(ctx context.Context, key string) error {
	err := config.Redis.Del(ctx, key).Err()
	if err != nil {
		logs.Error(err.Error())
	}
	return err
}
