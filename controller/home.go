package controller

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
)

type Home struct {
	service.App
}

func (h *Home) Init(r *znet.Engine) {
	// 静态资源目录，常用于放上传的文件
	r.Static("/static/", zfile.RealPathMkdir("./resource/static"))

	r.NotFoundHandler(func(c *znet.Context) {
		c.ApiJSON(404, "此路不通", nil)
	})
}

func (h *Home) Get(c *znet.Context) {
	c.ApiJSON(200, "Success", map[string]interface{}{
		"name": h.Conf.Base.Name,
	})
}
