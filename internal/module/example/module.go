package example

import (
	"reflect"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/app_core/service"
)

type Module struct {
	log *zlog.Logger

	// di  zdi.Invoker

	// service.App
	service.ModuleLifeCycle
}

var (
	_                = reflect.TypeOf(&Module{})
	_ service.Module = &Module{}
)

// Name 插件名称，非必须
func (p *Module) Name() string {
	return "Example"
}
