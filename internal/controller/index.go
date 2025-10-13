package controller

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	"app/internal/errcode"

	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
)

type Index struct {
	service.App
}

var _ = reflect.TypeOf(&Index{})

func (h *Index) Init(r *znet.Engine) error {
	r.Static("/static/", zfile.RealPath("./static"))

	return nil
}

func (h *Index) GET(c *znet.Context) (ztype.Map, error) {
	return ztype.Map{"hello": c.GetClientIP()}, nil
}

func (h *Index) HEAD(c *znet.Context) {
	c.SetStatus(http.StatusNoContent)
}

func (h *Index) GETHealth(c *znet.Context) (ztype.Map, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return ztype.Map{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   zcli.Version,
		"system": ztype.Map{
			"go_version":   runtime.Version(),
			"go_routines":  runtime.NumGoroutine(),
			"memory_alloc": memStats.Alloc,
			"memory_total": memStats.TotalAlloc,
			"memory_sys":   memStats.Sys,
			"gc_runs":      memStats.NumGC,
		},
		"config": ztype.Map{
			"debug":    h.Conf.Base.Debug,
			"timezone": h.Conf.Base.Zone,
		},
	}, nil
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
