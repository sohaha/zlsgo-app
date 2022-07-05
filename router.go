package main

import (
	"zlsapp/controller"
	"zlsapp/controller/wechat"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

func InitRouter(_ *service.Conf) []service.Router {
	return []service.Router{
		&controller.Home{},
		&wechat.Pay{},
		&wechat.Mp{},
		&wechat.Open{},
		&wechat.Qy{},
	}
}

func InitMiddleware(conf *service.Conf, app *service.App) []znet.Handler {
	return []znet.Handler{
		cors.Default(),
	}
}
