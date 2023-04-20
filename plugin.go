package main

import (
	"zlsapp/plugins/demo"
	"zlsapp/service"
)

func RegPlugin() []service.Plugin {
	return []service.Plugin{
		&demo.Plugin{},
	}
}
