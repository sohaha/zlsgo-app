package wechat

import (
	"zlsapp/logic/code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/zlsgo/wechat"
)

type Qy struct {
	service.App
}

func (w *Qy) Init(_ *znet.Engine) {
}

// GetAuth 网页授权
func (w *Qy) GetAuth(c *znet.Context) {
	wx := w.App.Wechat.Qy

	// TODO 本地使用内网穿透做测试，所以需要固定回调域名
	wx.SetOptions(wechat.WithRedirectDomain("http://mac.hw.73zls.com/"))

	json, ok, err := wx.Auth(c, "xxx", wechat.ScopePrivateinfo)
	if !ok {
		// 需要跳转到微信授权页面，所以这里直接返回空
		return
	}

	// 如果接口请求失败
	if err != nil {
		// 获取具体错误说明
		errMsg := wechat.ErrorMsg(err)
		code.ServerError.ApiResult(c, errMsg)
		return
	}

	// 获取授权成功，返回授权信息
	c.Log.Success("授权成功", json.Value())

	// 获取用户信息，企业微信使用 userTicket 获取用户信息
	userTicket := json.Get("user_ticket").String()
	token := json.Get("access_token").String()
	user, err := wx.GetAuthUserInfo(userTicket, token)
	if err != nil {
		code.ServerError.ApiResult(c, wechat.ErrorMsg(err))
		return
	}

	code.Success.ApiResult(c, user.Value())
}

func (w *Qy) GetAccessToken(c *znet.Context) {
	wx := w.App.Wechat.Qy

	token, err := wx.GetAccessToken()
	if err != nil {
		c.ApiJSON(211, wechat.ErrorMsg(err), nil)
		return
	}

	code.Success.ApiResult(c, map[string]interface{}{
		"time":        wx.GetAccessTokenExpiresInCountdown(),
		"accessToken": token,
	})
}

func (w *Qy) GetJsapiTicket(c *znet.Context) {
	wx := w.App.Wechat.Qy

	jsapiTicket, err := wx.GetJsapiTicket()
	if err != nil {
		code.ServerError.ApiResult(c, wechat.ErrorMsg(err))
		return
	}

	url := c.Host(true)
	jsSign, err := wx.GetJsSign(url)
	if err != nil {
		code.ServerError.ApiResult(c, wechat.ErrorMsg(err))
		return
	}

	code.Success.ApiResult(c, map[string]interface{}{
		"jsapiTicket": jsapiTicket,
		"jsSign":      jsSign,
		"url":         url,
	})
}

func (w *Qy) AnyReceiveMessage(c *znet.Context) {
	wx := w.App.Wechat.Qy

	body, _ := c.GetDataRaw()
	reply, err := wx.Reply(c.GetAllQueryMaps(),
		zstring.String2Bytes(body))
	if err != nil {
		c.String(211, err.Error())
		return
	}
	if c.Request.Method == "GET" {
		// Get 请求是响应微信发送的Token验证
		validMsg, err := reply.Valid()
		if err != nil {
			c.String(211, err.Error())
			return
		}
		c.String(200, validMsg)
		return
	}
	received, err := reply.Data()
	if err != nil {
		c.String(211, err.Error())
		return
	}
	c.Log.Info(received)
	replyXml := received.ReplyText("收到消息: " + received.MsgType)

	c.String(200, replyXml)
}
