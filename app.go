package main

import (
	"net/http"
	"strings"
	"zlsapp/internal/errcode"

	"github.com/zlsgo/app_core/service"
	"github.com/zlsgo/app_core/utils"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
)

var (
	c *service.Conf
)

func InitDI() zdi.Injector {
	di := zdi.New()

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))

	di.Provide(service.NewConf(func(o *conf.Option) {
		o.AutoCreate = true
	}))
	di.Provide(service.NewApp)
	di.Provide(service.NewWeb)

	di.Provide(RegMiddleware)
	di.Provide(RegRouter)
	di.Provide(RegRouterBefore)
	di.Provide(RegPlugin)
	di.Provide(RegTasks)
	di.Provide(RegErrHandler)

	return di
}

func RegErrHandler(app *service.App) znet.ErrHandlerFunc {
	return func(c *znet.Context, err error) {
		var (
			code       int32
			statusCode = http.StatusInternalServerError
			tag        = zerror.GetTag(err)
		)

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
		if app.Conf.Base.Debug && len(allErr) > 1 {
			zlog.Error(err)
		}
		if errMsg == "" {
			errMsg = "unknown error"
		}

		c.JSON(int32(statusCode), map[string]interface{}{
			"code": code,
			"msg":  errMsg,
		})
	}
}
func Start(di zdi.Injector) error {
	err := utils.InvokeErr(di.Invoke(service.InitPlugin))
	if err != nil {
		return zerror.With(err, "初始化插件失败")
	}

	err = di.Resolve(&c)
	if err != nil {
		return zerror.With(err, "初始化配置失败")
	}

	ztime.SetTimeZone(int(c.Base.Zone))

	err = utils.InvokeErr(di.Invoke(service.InitTask))
	if err != nil {
		return zerror.With(err, "定时任务启动失败")
	}

	err = utils.InvokeErr(di.Invoke(service.RunWeb))
	if err != nil {
		err = zerror.With(err, "服务启动失败")
	}
	return err
}

func Stop() {
}
