package main

import (
	"github.com/zlsgo/app_core/common"
	"github.com/zlsgo/app_core/service"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
)

func main() {
	zlog.ResetFlags(zlog.BitLevel)

	zcli.Name = service.AppName
	zcli.Version = "1.0.0"

	zcli.EnableDetach = true

	err := zutil.TryCatch(func() (err error) {
		di := InitDI()

		err = zcli.LaunchServiceRun(zcli.Name, "", func() {
			common.Fatal(Start(di))
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
