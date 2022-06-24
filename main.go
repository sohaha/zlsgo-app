package main

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
)

var di zdi.Injector

func main() {
	var c *service.Conf

	zcli.Name = "ZlsApp"
	zcli.Logo = `
_____                   
/  _  \  ______  ______  
/  /_\  \ \____ \ \____ \ 
/    |    \|  |_> >|  |_> >
\____|__  /|   __/ |   __/ 
	\/ |__|    |__|     `
	zcli.Version = "1.0.0"
	zcli.EnableDetach = true

	err := zutil.TryCatch(func() (err error) {
		di = InitDI()

		zerror.Panic(zerror.With(di.Resolve(&c), "配置读取失败"))

		zcli.Run(func() {
			_, err = di.Invoke(service.RunWeb)
			if err != nil {
				err = zerror.With(err, "服务启动失败")
			} else {
				_, _ = di.Invoke(service.StopWeb)
			}
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
