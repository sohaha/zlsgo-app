package main

import (
	"zlsapp/conf"
	"zlsapp/internal/utils"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
)

func main() {
	zcli.Name = conf.AppName
	zcli.Version = conf.AppVersion
	zcli.EnableDetach = true

	err := zutil.TryCatch(func() (err error) {
		utils.Fatal(Init())
		zcli.Run(func() {
			utils.Fatal(Start())
		})

		return err
	})

	if err != nil {
		if c == nil || !c.Base.Debug {
			zcli.Error(err.Error())
		} else {
			zlog.Errorf("%+v\n", err)
		}
	}
}
