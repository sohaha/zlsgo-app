package example

import (
	"reflect"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/app_core/service"
)

type Plugin struct {
	log *zlog.Logger

	// di  zdi.Invoker

	// service.App
	service.Pluginer
}

var (
	_                = reflect.TypeOf(&Plugin{})
	_ service.Plugin = &Plugin{}
)

// Name 插件名称，非必须
func (p *Plugin) Name() string {
	return "Example"
}
