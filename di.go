package main

import (
	"zlsapp/conf"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
)

func InitDI() zdi.Injector {
	di := zdi.New()

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))
	di.Map(service.InitConf(conf.DefaultSet))

	di.Provide(service.InitApp)
	di.Provide(service.InitWeb)

	di.Provide(service.InitDB)
	di.Provide(service.InitWechat)

	di.Provide(InitMiddleware)
	di.Provide(InitRouter)

	return di
}
