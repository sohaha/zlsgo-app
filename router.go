package main

import (
	"zlsapp/controller"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

// RegRouter 注册路由
func RegRouter(_ *service.Conf) *[]service.Controller {
	return &[]service.Controller{
		&controller.Index{},
	}
}

// RegMiddleware 注册全局中间件
func RegMiddleware(_ *service.Conf, _ *service.App) []znet.Handler {
	return []znet.Handler{
		cors.New(&cors.Config{
			ExposeHeaders: []string{"Authorization", "Re-Token"},
		}),
	}
}

// RegRouterBefore 注册路由前置处理
func RegRouterBefore(_ *service.Conf, _ *service.App) service.RouterBeforeProcess {
	return func(r *service.Web, app *service.App) {
	}
}
