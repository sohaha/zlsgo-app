package internal

import (
	"app/internal/plugins/example"

	"github.com/zlsgo/app_core/service"
)

func RegPlugin() []service.Plugin {
	return []service.Plugin{
		example.New(),
	}
}
