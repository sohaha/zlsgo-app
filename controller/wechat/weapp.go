package wechat

import (
	"zlsapp/logic/code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
)

type Weapp struct {
	service.App
}

func (w *Weapp) Init(_ *znet.Engine) {
}

// GetAuth 网页授权
func (w *Weapp) GetCodeAuth(c *znet.Context) {
	wx := w.App.Wechat.Weapp
	cod, ok := c.GetQuery("code")
	if !ok {
		code.InvalidInput.ApiResult(c, "code 不能为空")
		return
	}
	res, err := wx.GetAuthInfo(cod)
	if err != nil {
		code.InvalidInput.ApiResult(c, err)
		return
	}
	zlog.Success("授权成功", res)
	code.Success.ApiResult(c, res.Get("openid").String())
}
