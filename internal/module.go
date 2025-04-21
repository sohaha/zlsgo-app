package internal

import (
	"app/module/example"

	"github.com/zlsgo/app_core/service"
)

func RegModule() []service.Module {
	return []service.Module{
		example.New(),
	}
}
