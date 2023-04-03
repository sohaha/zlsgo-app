package main

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/zutil"
)

var (
	di zdi.Injector
	c  *service.Conf
)

func Init() error {
	// 全局时区
	ztime.SetTimeZone(8)

	err := zutil.TryCatch(func() (err error) {
		di = InitDI()
		_, _ = di.Invoke(func(_ *service.App, conf *service.Conf) {
			c = conf
		})

		return err
	})

	return err
}

func InitDI() zdi.Injector {
	di = zdi.New()

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))

	di.Provide(service.RegConf)
	di.Provide(service.RegApp)
	di.Provide(service.RegWeb)
	di.Provide(service.RegDB)

	di.Provide(RegMiddleware)
	di.Provide(RegRouter)
	di.Provide(RegRouterBefore)
	di.Provide(RegPlugin)
	di.Provide(RegTasks)

	return di
}

func Start() error {
	_, err := di.Invoke(service.InitModule)
	if err != nil {
		return zerror.With(err, "初始化插件失败")
	}

	_, err = di.Invoke(service.InitTask)
	if err != nil {
		return zerror.With(err, "定时任务启动失败")
	}

	_, err = di.Invoke(service.RunWeb)
	if err != nil {
		err = zerror.With(err, "服务启动失败")
	} else {
		_, _ = di.Invoke(service.StopWeb)
	}

	return err
}
