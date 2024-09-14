package pluginHandlers

import (
	"chatgpt-web-new-go/model"
	"chatgpt-web-new-go/router/base"

	"github.com/gin-gonic/gin"
)

func PluginHandler(c *gin.Context) {
	// pluginList, err := plugin.PluginList(c)
	// if err != nil {
	// 	logs.Error("plugin list error: %v", err)
	// 	base.Fail(c, "查询异常："+err.Error())
	// 	return
	// }

	base.Success(c, []*model.Plugin{})
}
