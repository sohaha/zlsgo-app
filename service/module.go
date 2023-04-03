package service

import (
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zlog"
)

type Module struct {
	di         zdi.Injector
	Name       string
	Before     func() error
	After      func() error
	Tasks      []Task
	Controller []Controller
}

func InitModule(modules []*Module, di zdi.Injector) (err error) {
	for i := range modules {
		plugin := (modules)[i]
		_, err = di.Invoke(plugin.Reg)
		if err != nil {
			return err
		}
	}

	for i := range modules {
		if err = (modules)[i].Init(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Module) Reg(di zdi.Injector) {
	ndi := zdi.New(di)
	PrintLog("Plugin", "Register: "+zlog.Log.ColorTextWrap(zlog.ColorLightGreen, d.Name))
	ndi.Map(ndi)
	d.di = ndi
	if d.Before != nil {
		d.Before()
	}
}

func (d *Module) Init() (err error) {
	_, err = d.di.Invoke(func(tasks *[]Task, controller *[]Controller, app *App, r *Web) {
		if len(d.Tasks) > 0 {
			*tasks = append(*tasks, d.Tasks...)
		}

		if len(d.Controller) > 0 {
			*controller = append(*controller, d.Controller...)
		}

		if d.After != nil {
			d.After()
		}
	})

	return
}

func (d *Module) DI() zdi.Injector {
	return d.di
}
