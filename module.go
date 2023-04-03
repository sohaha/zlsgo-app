package main

import (
	"zlsapp/module/demo"
	"zlsapp/service"
)

func RegPlugin() []*service.Module {
	return []*service.Module{
		// 插件列表
		demo.New(),
	}
}
