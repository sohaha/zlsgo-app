package demo

import (
	"reflect"

	"github.com/zlsgo/app_core/service"
)

type Plugin struct {
	service.App
	service.Pluginer
}

var (
	_                = reflect.TypeOf(&Plugin{})
	_ service.Plugin = &Plugin{}
)

// 插件名称，非必须
func (p *Plugin) Name() string {
	return "Example"
}
