package main

import (
	"github.com/sohaha/zlsgo/ztime"
	"zlsapp/internal/utils"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
)

var (
	c *service.Conf
)

func InitDI() zdi.Injector {
	di := zdi.New()

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))

	di.Provide(service.NewConf)
	di.Provide(service.NewApp)
	di.Provide(service.NewWeb)

	di.Provide(RegMiddleware)
	di.Provide(RegRouter)
	di.Provide(RegRouterBefore)
	di.Provide(RegPlugin)
	di.Provide(RegTasks)

	return di
}

func Start(di zdi.Injector) error {
	err := utils.InvokeErr(di.Invoke(service.InitPlugin))
	if err != nil {
		return zerror.With(err, "初始化插件失败")
	}

	err = di.Resolve(&c)
	if err != nil {
		return zerror.With(err, "初始化配置失败")
	}

	ztime.SetTimeZone(int(c.Base.Zone))

	err = utils.InvokeErr(di.Invoke(service.InitTask))
	if err != nil {
		return zerror.With(err, "定时任务启动失败")
	}

	err = utils.InvokeErr(di.Invoke(service.RunWeb))
	if err != nil {
		err = zerror.With(err, "服务启动失败")
	} else {
		_, _ = di.Invoke(service.StopWeb)
	}
	return err
}
