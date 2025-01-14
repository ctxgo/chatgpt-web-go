package config

import (
	"chatgpt-web-new-go/common/env"
	"context"
	"flag"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitConfig() {
	configFile := InitConfigWithFlag() // 使用配置文件路径

	v := viper.New()

	// env initializer
	e := env.GetEnv()

	//设置配置文件的名字
	v.SetConfigName(e)

	//添加配置文件所在的路径,注意在Linux环境下%GOPATH要替换为$GOPATH
	v.SetConfigFile(configFile)

	//设置配置文件类型
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// command line args
	commandLineConfig(v)

	Config = &Configuration{
		Db:          new(dbConfig),
		Redis:       new(redisConfig),
		AiClient:    new(gptConfig),
		EmailServer: new(emailServerConfig),
	}
	err := v.Unmarshal(Config)
	if err != nil {
		panic(err)
	}

	log.Printf("global config: %v \n", Config)
	go watchConfigChange(v)
}

func InitConfigWithFlag() string {
	var configfile = pflag.StringP("config", "c", "", "Use -c to set your config file")

	pflag.Parse()

	if pflag.NArg() != 0 || *configfile == "" {
		pflag.Usage()
		os.Exit(1)
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(*configfile); err != nil {
		panic(err)
	}

	return *configfile
}

// 监听配置文件的修改和变动
func watchConfigChange(v *viper.Viper) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("watchConfigChange panic recover: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	v.WatchConfig()
	//监听回调函数
	watch := func(e fsnotify.Event) {
		log.Printf("Config file is changed: %s \n", e.String())
		cancel()
	}
	v.OnConfigChange(watch)
	<-ctx.Done()
}

func commandLineConfig(v *viper.Viper) {
	pflag.String("token", "", "please input the token")
	pflag.Int("adminUserId", 0, "please input the admin user id")
	pflag.Bool("debug", false, "please input the debug flag")
	//获取标准包的flag
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	//BindFlag
	err := v.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
}
