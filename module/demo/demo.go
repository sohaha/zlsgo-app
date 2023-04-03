package demo

import (
	"zlsapp/service"
)

func New() (s *service.Module) {
	s = &service.Module{
		Name:  "示例插件",
		Tasks: []service.Task{},
		Before: func() error {
			return nil
		},
		After: func() error {
			return nil
		},
		Controller: []service.Controller{
			&Index{},
		},
	}
	return
}
