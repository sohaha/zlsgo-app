package internal

import (
	"context"
	"net/http"
	"strings"

	"app/internal/errcode"

	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/zutil"
	"github.com/zlsgo/conf"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
)

func InitDI(ctx ...context.Context) zdi.Injector {
	di := zdi.New()

	if len(ctx) > 0 {
		di.Map(ctx[0])
	} else {
		di.Map(context.Background())
	}

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))

	di.Provide(service.NewConf(func(o conf.Options) conf.Options {
		o.AutoCreate = true
		return o
	}))

	di.Provide(service.NewApp(func(o service.BaseConf) service.BaseConf {
		// Can be set via environment variables or configuration file
		o.Port = zutil.Getenv("PORT", "8181")
		o.Debug = strings.ToLower(zutil.Getenv("DEBUG", "false")) == "true"
		return o
	}))

	di.Provide(service.NewWeb())

	di.Provide(RegMiddleware)
	di.Provide(RegRouter)
	di.Provide(RegRouterBefore)
	di.Provide(RegModule)
	di.Provide(RegTasks)
	di.Provide(RegErrHandler)

	return di
}

// ErrorMapping 错误映射配置
type ErrorMapping struct {
	HTTPStatus int
	ErrorCode  int32
}

// 默认错误映射配置
var defaultErrorMappings = map[zerror.TagKind]ErrorMapping{
	zerror.Internal:         {http.StatusInternalServerError, int32(errcode.ServerError)},
	zerror.InvalidInput:     {http.StatusBadRequest, int32(errcode.InvalidInput)},
	zerror.PermissionDenied: {http.StatusForbidden, int32(errcode.PermissionDenied)},
	zerror.Unauthorized:     {http.StatusUnauthorized, int32(errcode.Unauthorized)},
}

func RegErrHandler(app *service.App) znet.ErrHandlerFunc {
	return func(c *znet.Context, err error) {
		statusCode, errorCode, errMsg := processError(err, app.Conf.Base.Debug)

		c.JSON(int32(statusCode), map[string]interface{}{
			"code": errorCode,
			"msg":  errMsg,
		})
	}
}

// processError 处理错误并返回HTTP状态码、错误码和错误消息
func processError(err error, debug bool) (statusCode int, errorCode int32, errMsg string) {
	tag := zerror.GetTag(err)

	// 尝试从预定义映射中获取错误信息
	if mapping, ok := defaultErrorMappings[tag]; ok {
		statusCode = mapping.HTTPStatus
		errorCode = mapping.ErrorCode
	} else {
		// 处理自定义错误码
		if code, hasCode := zerror.UnwrapCode(err); hasCode && code != 0 {
			errorCode = int32(code)
		} else {
			errorCode = int32(errcode.ServerError)
		}

		// 处理自定义HTTP状态码
		statusCode = http.StatusInternalServerError
		if tag != zerror.None {
			if customStatus := ztype.ToInt(string(tag)); customStatus > 0 {
				statusCode = customStatus
			}
		}
	}

	// 处理错误消息
	allErr := zerror.UnwrapErrors(err)
	errMsg = strings.Join(allErr, ": ")
	if errMsg == "" {
		errMsg = "unknown error"
	}

	// 调试模式下记录详细错误
	if debug && len(allErr) > 1 {
		zlog.Error(err)
	}

	return
}

func Init(di zdi.Injector, loadModule bool) (c *service.Conf, err error) {
	if loadModule {
		err = di.InvokeWithErrorOnly(service.InitModule)
		if err != nil {
			return nil, zerror.With(err, "failed to init module")
		}
	}

	err = di.Resolve(&c)
	if err != nil {
		return nil, zerror.With(err, "failed to init config")
	}

	ztime.SetTimeZone(int(c.Base.Zone))
	return
}

func Start(di zdi.Injector) error {
	err := di.InvokeWithErrorOnly(service.InitTask)
	if err != nil {
		return zerror.With(err, "timed task launch failed")
	}

	err = di.InvokeWithErrorOnly(service.RunWeb)
	if err != nil {
		err = zerror.With(err, "service startup failed")
	}
	return err
}

func Stop(di zdi.Invoker, ps []service.Module) {
}
