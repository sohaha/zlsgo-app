package service

import (
	"reflect"
	"zlsapp/internal/utils"

	"github.com/sohaha/zlsgo/zerror"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zlog"
)

type Plugin interface {
	Name() string
	Tasks() []Task
	Before() error
	After() error
	Controller() []Controller
}

func InitPlugin(ps []Plugin, di zdi.Injector) (err error) {

	for _, p := range ps {
		pdi := reflect.Indirect(reflect.ValueOf(p)).FieldByName("DI")
		if pdi.IsValid() {
			switch pdi.Type().String() {
			case "zdi.Invoker", "zdi.Injector":
				pdi.Set(reflect.ValueOf(di))
			}
		}
		err := zerror.TryCatch(func() error {
			return p.Before()
		})
		if err != nil {
			return zerror.With(err, p.Name())
		}

		di.Map(p)
	}

	return utils.InvokeErr(di.Invoke(func(app *App, tasks *[]Task, controller *[]Controller, r *Web) error {
		for _, p := range ps {
			*tasks = append(*tasks, p.Tasks()...)
			*controller = append(*controller, p.Controller()...)

			conf := reflect.Indirect(reflect.ValueOf(p)).FieldByName("Conf")
			if conf.IsValid() && conf.Type().String() == "*service.Conf" {
				conf.Set(reflect.ValueOf(app.Conf))
			}
			err := zerror.TryCatch(func() error {
				return p.After()
			})
			if err != nil {
				return zerror.With(err, p.Name())
			}

			PrintLog("Plugin", zlog.Log.ColorTextWrap(zlog.ColorLightGreen, p.Name()))
		}

		return nil
	}))
}
