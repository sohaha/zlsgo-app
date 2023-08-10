package internal

import (
	"net/http"
	"strings"

	"app/internal/errcode"

	"github.com/sohaha/zlsgo/ztime"
	"github.com/zlsgo/conf"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
)

func InitDI() zdi.Injector {
	di := zdi.New()

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))

	di.Provide(service.NewConf(func(o *conf.Option) {
		o.AutoCreate = true
	}))

	di.Provide(service.NewApp(func(o *service.BaseConf) {
		o.Port = "8181"
	}))

	di.Provide(service.NewWeb())

	di.Provide(RegMiddleware)
	di.Provide(RegRouter)
	di.Provide(RegRouterBefore)
	di.Provide(RegPlugin)
	di.Provide(RegTasks)
	di.Provide(RegErrHandler)

	return di
}

func RegErrHandler(app *service.App) znet.ErrHandlerFunc {
	var tagMap = map[zerror.TagKind]int{
		zerror.Internal:         http.StatusInternalServerError,
		zerror.InvalidInput:     http.StatusBadRequest,
		zerror.PermissionDenied: http.StatusForbidden,
		zerror.Unauthorized:     http.StatusUnauthorized,
	}

	return func(c *znet.Context, err error) {
		var (
			code       int32
			statusCode = http.StatusInternalServerError
			tag        = zerror.GetTag(err)
		)
		if val, ok := tagMap[tag]; ok {
			statusCode = val
			code = int32(errcode.ServerError)
		} else {
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

func Init(di zdi.Injector, loadPlugin bool) (c *service.Conf, err error) {
	if loadPlugin {
		err = di.InvokeWithErrorOnly(service.InitPlugin)
		if err != nil {
			return nil, zerror.With(err, "初始化插件失败")
		}
	}

	err = di.Resolve(&c)
	if err != nil {
		return nil, zerror.With(err, "初始化配置失败")
	}

	ztime.SetTimeZone(int(c.Base.Zone))
	return
}

func Start(di zdi.Injector) error {
	err := di.InvokeWithErrorOnly(service.InitTask)
	if err != nil {
		return zerror.With(err, "定时任务启动失败")
	}

	err = di.InvokeWithErrorOnly(service.RunWeb)
	if err != nil {
		err = zerror.With(err, "服务启动失败")
	}
	return err
}

func Stop() {
}
