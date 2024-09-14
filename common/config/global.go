package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var (
	Config      *Configuration
	DB          *gorm.DB // DB instance
	Redis       *redis.Client
	Gcron       *cron.Cron // cron
	EmailDialer *gomail.Dialer
)
