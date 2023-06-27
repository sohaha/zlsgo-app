package main

import (
	"app/plugins/demo"

	"github.com/zlsgo/app_core/service"
)

func RegPlugin() []service.Plugin {
	return []service.Plugin{
		demo.New(),
	}
}
