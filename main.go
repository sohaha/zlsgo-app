package main

import (
	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
	"github.com/zlsgo/app_core/common"
	"github.com/zlsgo/app_core/service"
)

func main() {
	zlog.ResetFlags(zlog.BitLevel)

	zcli.Name = service.AppName
	zcli.EnableDetach = true
	zcli.Version = "1.0.0"

	err := zutil.TryCatch(func() (err error) {
		di := InitDI()

		err = zcli.LaunchServiceRun(zcli.Name, "", func() {
			common.Fatal(Init(di, true))
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
