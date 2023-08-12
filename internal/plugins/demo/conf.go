package demo

import (
	"github.com/zlsgo/app_core/service"
)

// Conf 插件配置
type Conf struct {
	Dev  bool
	Text string
}

// ConfKey 配置文件key
func (Conf) ConfKey() string {
	return "demo"
}

// 可以禁止配置自动写入配置文件
// func (Conf) DisableWrite() bool {
// 	return true
// }

// New 实例化插件
func New() *Plugin {
	// 配置文件默认值
	defaultConf := Conf{
		Dev:  true,
		Text: "这是一个插件配置",
	}
	service.DefaultConf = append(service.DefaultConf, defaultConf)

	p := &Plugin{}
	p.Pluginer = service.Pluginer{
		OnLoad: func() error {
			// 配置解析完成后执行
			return p.DI.InvokeWithErrorOnly(func(conf *service.Conf) error {
				p.Log.Debug("插件配置：", conf.Get(defaultConf.ConfKey()))
				return nil
			})
		},
		OnStart: func() error {
			// 全部插件加载完成后执行
			return nil
		},
		OnDone: func() error {
			// 全部插件启动后执行
			return nil
		},
		Service: &service.PluginService{
			Controllers: []service.Controller{&Index{}},
			Tasks: []service.Task{
				{Name: "demo task", Cron: "1 * * * * * *", Run: func() {
					p.Log.Debug("定时执行任务")
				}},
			},
		},
	}

	return p
}
