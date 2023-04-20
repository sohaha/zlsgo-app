package service

import (
	"zlsapp/internal/utils"

	"github.com/sohaha/zlsgo/zerror"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zlog"
)

type Plugin interface {
	Name() string
	Init(di zdi.Injector) error
	Tasks() []Task
	Before() error
	After() error
	Controller() []Controller
}

func InitPlugin(ps []Plugin, di zdi.Injector) (err error) {
	for _, p := range ps {
		if err := p.Init(di); err != nil {
			return zerror.With(err, p.Name())
		}

		if err := p.Before(); err != nil {
			return zerror.With(err, p.Name())
		}

		di.Map(p)
	}

	return utils.InvokeErr(di.Invoke(func(a *App, tasks *[]Task, controller *[]Controller, r *Web) error {
		for _, p := range ps {
			*tasks = append(*tasks, p.Tasks()...)
			*controller = append(*controller, p.Controller()...)

			if err := p.After(); err != nil {
				return zerror.With(err, p.Name())
			}

			PrintLog("Plugin", zlog.Log.ColorTextWrap(zlog.ColorLightGreen, p.Name()))
		}

		return nil
	}))
}
