package utils

import (
	"reflect"
	"strings"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
)

func Fatal(err error) {
	if err == nil {
		return
	}
	zlog.Fatal(strings.Join(zerror.UnwrapErrors(err), ": "))
}

func InvokeErr(v []reflect.Value, err error) error {
	if err != nil {
		return err
	}
	if len(v) == 0 {
		return nil
	}
	err, ok := v[0].Interface().(error)
	if !ok {
		return nil
	}
	return err
}
