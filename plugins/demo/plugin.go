package demo

import (
	"reflect"

	"github.com/zlsgo/app_core/service"
)

type Plugin struct {
	service.App
}

var (
	_                = reflect.TypeOf(&Index{})
	_ service.Plugin = &Plugin{}
)

func (p *Plugin) Name() string {
	// 插件名称
	return "Sample"
}

func (p *Plugin) Controller() []service.Controller {
	// 定义控制器
	return []service.Controller{
		&Index{},
	}
}

func (p *Plugin) Tasks() []service.Task {
	// 定义定时任务
	return []service.Task{}
}

func (p *Plugin) Load() error {
	// 配置解析完成后执行
	return nil
}

func (p *Plugin) Start() error {
	// 全部插件加载完成后执行
	return nil
}

func (p *Plugin) Done() error {
	// 全部插件启动后执行
	return nil
}
