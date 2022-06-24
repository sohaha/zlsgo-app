package main

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
)

func InitDI() zdi.Injector {
	di := zdi.New()

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))

	di.Provide(service.InitConf)
	di.Provide(service.InitApp)
	di.Provide(service.InitWeb)

	di.Provide(InitMiddleware)
	di.Provide(InitRouter)

	return di
}
