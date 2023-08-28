package demo

import (
	"reflect"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Index struct {
	di   zdi.Invoker
	Path string
}

var (
	_ = reflect.TypeOf(&Index{})
)

func (h *Index) Init(r *znet.Engine) error {
	// 注册中间件
	r.Use(func(c *znet.Context) {
		c.Next()
	})
	return nil
}

func (h *Index) GET(c *znet.Context) string {
	return "ok"
}

func (h *Index) IDGET(c *znet.Context) (ztype.Map, error) {
	return ztype.Map{
		"id": c.GetParam("id"),
	}, nil
}
