package redis

import (
	"chatgpt-web-new-go/common/config"
	"chatgpt-web-new-go/pkgs/retry"
	"context"

	"github.com/redis/go-redis/v9"
)

func Init() {
	redisConfig := config.Config.Redis
	config.Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password, // Redis 服务器没有设置密码
		DB:       redisConfig.DB,       // 使用默认数据库
	})
	err := retry.Retry(func() error {
		return connect(config.Redis)
	})
	if err != nil {
		panic(err)
	}
}

// connect connect test
func connect(r *redis.Client) error {
	var err error
	_, err = r.Ping(context.TODO()).Result()
	return err
}
