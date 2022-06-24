package main

import (
	"zlsapp/controller"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

func InitRouter(_ *service.Conf) []service.Router {
	return []service.Router{
		&controller.Home{},
	}
}

func InitMiddleware(conf *service.Conf, app *service.App) []znet.Handler {
	return []znet.Handler{
		cors.Default(),
	}
}
