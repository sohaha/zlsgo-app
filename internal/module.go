package internal

import (
	"app/internal/module/example"

	"github.com/zlsgo/app_core/service"
)

func RegModule() []service.Module {
	return []service.Module{
		example.New(),
	}
}
