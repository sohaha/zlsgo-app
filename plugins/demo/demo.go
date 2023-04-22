package demo

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
)

type Plugin struct {
	DI   zdi.Invoker
	Conf *service.Conf
}

var _ service.Plugin = &Plugin{}

func New() *Plugin {
	return &Plugin{}
}

func (d *Plugin) Name() string {
	return "示例插件"
}

func (d *Plugin) Tasks() []service.Task {
	return []service.Task{}
}

func (d *Plugin) Before() error {
	return nil
}

func (d *Plugin) After() error {
	return nil
}

func (d *Plugin) Controller() []service.Controller {
	return []service.Controller{
		&Index{},
	}
}
