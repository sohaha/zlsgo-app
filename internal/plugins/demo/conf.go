package demo

import (
	"github.com/zlsgo/app_core/service"
)

// Conf 插件配置
type Conf struct {
	Dev bool
}

// ConfKey 配置文件key
func (Conf) ConfKey() string {
	return "demo"
}

// New 实例化插件
func New() *Plugin {
	// 配置文件默认值
	service.DefaultConf = append(service.DefaultConf, Conf{
		Dev: true,
	})
	return &Plugin{}
}
