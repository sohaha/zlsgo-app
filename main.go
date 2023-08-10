package main

import (
	"app/internal"
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

	var c *service.Conf
	err := zutil.TryCatch(func() (err error) {
		di := internal.InitDI()

		err = zcli.LaunchServiceRun(zcli.Name, "", func() {
			c, err = internal.Init(di, true)
			common.Fatal(err)
			common.Fatal(internal.Start(di))
		})

		_, _ = di.Invoke(internal.Stop)
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
