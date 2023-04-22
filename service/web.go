package service

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"zlsapp/internal/errcode"
	"zlsapp/internal/utils"

	"github.com/arl/statsviz"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zpprof"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
)

type (
	Web struct {
		*znet.Engine
		hijacked []func(c *znet.Context) bool
	}
	// Controller 控制器函数
	Controller interface {
		Init(r *znet.Engine)
	}
	// RouterBeforeProcess 控制器前置处理
	RouterBeforeProcess func(r *Web, app *App)
	Template            struct {
		DIR    string
		Global ztype.Map
	}
)

func (w *Web) AddHijack(fn func(c *znet.Context) bool) {
	if fn == nil {
		return
	}
	w.hijacked = append(w.hijacked, fn)
}

func (w *Web) GetHijack() []func(c *znet.Context) bool {
	return w.hijacked
}

// NewWeb 初始化 WEB
func NewWeb(app *App, middlewares []znet.Handler) (*Web, *znet.Engine) {
	r := znet.New()
	r.Log = app.Log
	zlog.Log = r.Log

	r.BindStructSuffix = ""
	r.BindStructDelimiter = "-"
	r.SetAddr(app.Conf.Base.Port)

	isDebug := app.Conf.Base.Debug
	if isDebug {
		r.SetMode(znet.DebugMode)
	} else {
		r.SetMode(znet.ProdMode)
	}

	if app.Conf.Base.Pprof {
		zpprof.Register(r, app.Conf.Base.PprofToken)
	}

	if app.Conf.Base.Statsviz {
		r.GET(`/debug/statsviz{*:[\S]*}`, func(c *znet.Context) {
			q := c.GetParam("*")
			if q == "" {
				c.Redirect("/debug/statsviz/")
				return
			}
			if q == "/ws" {
				statsviz.Ws(c.Writer, c.Request)
				return
			}
			statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(c.Writer, c.Request)
		})
	}

	r.Use(znet.RewriteErrorHandler(func(c *znet.Context, err error) {
		var code int32
		statusCode := http.StatusInternalServerError
		tag := zerror.GetTag(err)
		switch tag {
		case zerror.Internal:
			statusCode = http.StatusInternalServerError
			code = int32(errcode.ServerError)
		case zerror.InvalidInput:
			statusCode = http.StatusBadRequest
			code = int32(errcode.InvalidInput)
		case zerror.PermissionDenied:
			statusCode = http.StatusForbidden
			code = int32(errcode.PermissionDenied)
		case zerror.Unauthorized:
			statusCode = http.StatusUnauthorized
			code = int32(errcode.Unauthorized)
		default:
			errCode, ok := zerror.UnwrapCode(err)
			if ok && errCode != 0 {
				code = int32(errCode)
			} else {
				code = int32(errcode.ServerError)
			}
			if tag != zerror.None {
				statusCode = ztype.ToInt(string(tag))
			}
		}

		allErr := zerror.UnwrapErrors(err)
		errMsg := strings.Join(allErr, ": ")
		if isDebug && len(allErr) > 1 {
			zlog.Error(err)
		}
		if errMsg == "" {
			errMsg = "unknown error"
		}

		c.JSON(int32(statusCode), map[string]interface{}{
			"code": code,
			"msg":  errMsg,
		})
	}))

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	return &Web{
		Engine: r,
	}, r
}

func RunWeb(r *Web, app *App, controllers *[]Controller) {
	_, err := app.DI.Invoke(func(after RouterBeforeProcess) {
		after(r, app)
	})
	if err != nil && !strings.Contains(err.Error(), "value not found for type service.RouterBeforeProcess") {
		utils.Fatal(err)
	}

	utils.Fatal(initRouter(app, r, *controllers))

	// r.StartUp()
	znet.Run()
}

func initRouter(app *App, _ *Web, controllers []Controller) (err error) {
	_, _ = app.DI.Invoke(func(r *Web) {
		for i := range controllers {
			c := controllers[i]
			err = zutil.TryCatch(func() (err error) {
				typeOf := reflect.TypeOf(c).Elem()
				controller := strings.TrimPrefix(typeOf.String(), "controller.")
				controller = strings.Replace(controller, ".", "/", -1)
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

				cDI := reflect.Indirect(reflect.ValueOf(c)).FieldByName("DI")
				if cDI.IsValid() {
					switch cDI.Type().String() {
					case "zdi.Invoker", "zdi.Injector":
						cDI.Set(reflect.ValueOf(app.DI))
					}
				}

				name := ""
				cName := reflect.Indirect(reflect.ValueOf(c)).FieldByName("Path")

				if cName.IsValid() && cName.String() != "" {
					name = zstring.CamelCaseToSnakeCase(cName.String(), "/")
				} else {
					name = zstring.CamelCaseToSnakeCase(controller, "/")
				}

				lname := strings.Split(name, "/")
				if lname[len(lname)-1] == "index" {
					name = strings.Join(lname[:len(lname)-1], "/")
					name = strings.TrimSuffix(name, "/")
				}
				if name == "" {
					err = r.BindStruct(name, c)
				} else {
					err = r.Group("/").BindStruct(name, c)
				}
				return err
			})

			if err != nil {
				err = fmt.Errorf("初始化路由失败: %w", err)
				return
			}
		}
	})
	return
}
