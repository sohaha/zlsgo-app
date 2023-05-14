package demo

import (
	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Index struct {
	service.App
}

func (h *Index) Init(r *znet.Engine) {
	// 注册中间件
	r.Use(func(c *znet.Context) {
		c.Next()
	})
}

func (h *Index) GET(c *znet.Context) string {
	return "ok"
}

func (h *Index) IDGET(c *znet.Context) (ztype.Map, error) {
	return ztype.Map{
		"id": c.GetParam("id"),
	}, nil
}
