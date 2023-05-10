package main

import (
	"zlsapp/plugins/demo"

	"github.com/zlsgo/app_core/service"
)

func RegPlugin() []service.Plugin {
	return []service.Plugin{
		demo.New(),
	}
}
