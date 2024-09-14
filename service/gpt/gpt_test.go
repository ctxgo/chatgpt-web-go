package gpt

import (
	"chatgpt-web-new-go/common/aiclient"
	"chatgpt-web-new-go/common/config"
	"chatgpt-web-new-go/common/logs"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// config init
	config.InitConfig()

	// log init
	logs.LogInit()

	// aiClient init
	aiclient.InitClients()

	code := m.Run()
	os.Exit(code)
}
