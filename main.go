package main

import (
	"zlsapp/internal/utils"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
)

func main() {
	zlog.ResetFlags(zlog.BitLevel)

	zcli.Name = service.AppName
	zcli.EnableDetach = true

	err := zutil.TryCatch(func() (err error) {
		di := InitDI()

		err = zcli.LaunchServiceRun(zcli.Name, "", func() {
			utils.Fatal(Start(di))
		})

		_, _ = di.Invoke(Stop)
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
