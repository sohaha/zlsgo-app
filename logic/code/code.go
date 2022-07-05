package code

import (
	"github.com/sohaha/zlsgo/zerror"
)

type ErrCode zerror.ErrCode

const (
	Success      ErrCode = 0
	ServerError  ErrCode = 10000
	InvalidInput ErrCode = 20000
)
