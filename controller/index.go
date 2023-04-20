package controller

import (
	"net/http"

	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
)

type Index struct {
	service.App
	DI zdi.Injector
}

func (h *Index) Init(r *znet.Engine) {
	var web *service.Web
	h.DI.Resolve(&web)

	// 开放静态资源目录
	r.Static("/static/", zfile.RealPath("./static"))

	// 处理不存在的路由请求
	r.NotFoundHandler(func(c *znet.Context) {
		hijacked := web.GetHijack()

		for i := range hijacked {
			if hijacked[i](c) {
				return
			}
		}

		path := c.Request.URL.Path
		if path == "/" {
			c.JSON(http.StatusOK,
				znet.ApiData{
					Code: 0,
					Data: ztype.Map{
						"App": service.AppName,
						"now": ztime.Now(),
					},
				})
			return
		}
		c.JSON(http.StatusNotFound,
			znet.ApiData{
				Code: 404,
				Msg:  "Not Found",
				Data: struct{}{},
			})
	})
}
