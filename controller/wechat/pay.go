package wechat

import (
	"strings"
	"zlsapp/logic/code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/zlsgo/wechat"
)

type Pay struct {
	service.App
}

func (w *Pay) Init(r *znet.Engine) {
}

func (w *Pay) GetJsSign(c *znet.Context) {
	openid := "oj3r9s0kqWLyZJRQ2FTim7QnYmGI"
	order := wechat.NewPayOrder(openid, 101, c.GetClientIP(), "Body")
	notifyUrl := strings.Replace(c.Request.URL.Path, "js-sign", "notify", 1)
	prepayID, err := w.Wechat.Pay.UnifiedOrder(w.Conf.Wechat.Pay.AppID, order, notifyUrl)
	if err != nil {
		code.InvalidInput.ApiResult(c, err)
		return
	}
	sign := w.Wechat.Pay.JsSign(w.Conf.Wechat.Pay.AppID, prepayID)
	code.Success.ApiResult(c, sign)
}

func (w *Pay) AnyNotify(c *znet.Context) {
	raw, _ := c.GetDataRaw()

	data, err := w.Wechat.Pay.Notify(raw)
	if err != nil {
		code.InvalidInput.ApiResult(c, err)
		return
	}

	switch data.Type {
	case wechat.PayNotify:
		c.Log.Success("支付成功", data.Data)
	case wechat.RefundNotify:
		c.Log.Success("退款成功", data.Data)
	}

	c.Byte(200, data.Response)
}
