package service

import (
	"github.com/zlsgo/wechat"
)

type Wechat struct {
	Mp    *wechat.Engine
	Open  *wechat.Engine
	Qy    *wechat.Engine
	Weapp *wechat.Engine
	Pay   *wechat.Pay
}

func InitWechat(c *Conf) *Wechat {
	wx := &Wechat{}

	if len(c.Wechat.Mp.AppSecret) != 0 {
		wx.Mp = wechat.New(&c.Wechat.Mp)
	}

	if len(c.Wechat.Qy.Secret) != 0 {
		wx.Qy = wechat.New(&c.Wechat.Qy)
	}

	if len(c.Wechat.Open.AppSecret) != 0 {
		wx.Open = wechat.New(&c.Wechat.Open)
	}

	if len(c.Wechat.Weapp.AppSecret) != 0 {
		wx.Weapp = wechat.New(&c.Wechat.Weapp)
	}

	if len(c.Wechat.Pay.MchId) != 0 {
		p := c.Wechat.Pay
		wx.Pay = wechat.NewPay(wechat.Pay{
			MchId:    p.MchId,
			Key:      p.Key,
			CertPath: p.CertPath,
			KeyPath:  p.KeyPath,
		})
		if p.Sandbox {
			wx.Pay.Sandbox(true)
		}
	}

	if c.Wechat.Debug {
		wechat.Debug()
	}
	return wx
}
