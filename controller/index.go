package controller

import (
	"net/http"

	"app/internal/errcode"

	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
)

type Index struct {
	service.App
}

func (h *Index) Init(r *znet.Engine) error {
	// 开放静态资源目录
	r.Static("/static/", zfile.RealPath("./static"))

	// 处理不存在的路由请求
	r.NotFoundHandler(func(c *znet.Context) {
		c.JSON(http.StatusNotFound,
			znet.ApiData{
				Code: 404,
				Msg:  "Not Found",
				Data: struct{}{},
			})
	})
	return nil
}

func (h *Index) GetError(r *znet.Context) error {
	return errcode.InvalidInput.WrapText("test Error")
}
