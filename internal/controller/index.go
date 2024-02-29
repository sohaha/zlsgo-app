package controller

import (
	"app/internal/errcode"

	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
)

type Index struct {
	service.App
}

func (h *Index) Init(r *znet.Engine) error {
	// 开放静态资源目录
	r.Static("/static/", zfile.RealPath("./static"))

	return nil
}

func (h *Index) GET(r *znet.Context) (ztype.Map, error) {
	return ztype.Map{"hello": "world"}, nil
}

type Body struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string
}

func (h *Index) POST(r *znet.Context) (ztype.Map, error) {
	var body Body

	valid := r.ValidRule().Required()
	m := map[string]zvalid.Engine{
		"name":    valid.MinUTF8Length(2, "姓名最短两个字").SetAlias("姓名"),
		"age":     valid.IsNumber("年龄必须是整数").SetAlias("年龄"),
		"Address": valid.IsChinese().SetAlias("地址"),
	}
	if err := r.BindValid(&body, m); err != nil {
		return nil, errcode.InvalidInput.WrapErr(err)
	}

	return ztype.Map{"body": body}, nil
}
