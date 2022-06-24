package service

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/zutil"
)

type (
	// App 控制器关联对象
	App struct {
		Di   zdi.Injector
		Conf *Conf
		Log  *zlog.Logger
	}
	// Router 控制器函数
	Router interface {
		Init(r *znet.Engine)
	}
)

// InitWeb 初始化 WEB
func InitWeb(app *App, middlewares []znet.Handler) *znet.Engine {
	r := znet.New()
	r.Log = app.Log
	r.BindStructSuffix = ".go"
	r.BindStructDelimiter = "-"
	r.SetAddr(app.Conf.Base.Port)

	isDebug := app.Conf.Base.Debug
	if isDebug {
		r.SetMode(znet.DebugMode)
	}

	r.Use(znet.Recovery(func(c *znet.Context, err error) {
		if isDebug {
			c.Log.Track(err.Error(), 20)
		}
		c.ApiJSON(500, err.Error(), struct{}{})
	}))

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	return r
}

func RunWeb(r *znet.Engine, app *App, controllers []Router) {
	zlog.Tips("启动 web 服务...")

	for _, c := range controllers {
		err := zutil.TryCatch(func() error {
			typeOf := reflect.TypeOf(c).Elem()
			controller := strings.TrimPrefix(typeOf.String(), "controller.")
			api := -1
			for i := 0; i < typeOf.NumField(); i++ {
				if typeOf.Field(i).Type.String() == "service.App" {
					api = i
					break
				}
			}
			if api == -1 {
				return fmt.Errorf("%s not a legitimate controller", controller)
			}
			reflect.ValueOf(c).Elem().Field(api).Set(reflect.ValueOf(*app))

			name := zstring.CamelCaseToSnakeCase(controller, "/")
			if name == "home" {
				name = ""
			}

			return r.BindStruct(name, c)
		})
		zerror.Panic(err)
	}

	znet.Run()
}

func StopWeb(_ *znet.Engine, _ *App) {
	znet.SetShutdown(func() {

	})
}
