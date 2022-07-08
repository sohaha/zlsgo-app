package wechat

import (
	"strconv"
	"time"

	"zlsapp/logic/code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"

	"github.com/zlsgo/wechat"
)

type Mp struct {
	service.App
}

func (w *Mp) Init(_ *znet.Engine) {
}

// GetAuth 网页授权
func (w *Mp) GetAuth(c *znet.Context) {
	wx := w.App.Wechat.Mp

	// TODO 本地使用内网穿透做测试，所以需要固定回调域名
	wx.SetOptions(wechat.WithRedirectDomain("http://mac.hw.73zls.com/"))

	json, ok, err := wx.Auth(c, "xxx", wechat.ScopeUserinfo)
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

	// 获取用户信息
	openid := json.Get("openid").String()
	token := json.Get("access_token").String()
	user, err := wx.GetAuthUserInfo(openid, token)
	if err != nil {
		code.ServerError.ApiResult(c, wechat.ErrorMsg(err))
		return
	}

	code.Success.ApiResult(c, user.Value())
}

// GetCodeAuth 通过 Code 获取用户 Openid
func (w *Mp) GetCodeAuth(c *znet.Context) {
	// https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx6a24b584b45b6791&redirect_uri=encodeURIComponent编码回调地址&response_type=code&scope=snsapi_userinfo&connect_redirect=1#wechat_redirect
	cod, ok := c.GetQuery("code")
	if !ok {
		code.InvalidInput.ApiResult(c, "code 不能为空")
	}
	_ = cod
	wx := w.App.Wechat.Mp
	res, err := wx.GetAuthInfo(cod)
	if err != nil {
		code.InvalidInput.ApiResult(c, err)
		return
	}
	zlog.Success("授权成功", res)
	code.Success.ApiResult(c, res.Get("openid").String())
}

// GetAccessToken 获取 AccessToken
func (w *Mp) GetAccessToken(c *znet.Context) {
	wx := w.App.Wechat.Mp

	token, err := wx.GetAccessToken()
	if err != nil {
		code.ServerError.ApiResult(c, wechat.ErrorMsg(err))
		return
	}

	code.Success.ApiResult(c, map[string]interface{}{
		"time":        wx.GetAccessTokenExpiresInCountdown(),
		"accessToken": token,
	})
}

// GetJsapiTicket 获取 JsapiTicket
func (w *Mp) GetJsapiTicket(c *znet.Context) {
	wx := w.App.Wechat.Mp

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

func (w *Mp) AnyReceiveMessage(c *znet.Context) {
	wx := w.App.Wechat.Mp

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
	replyXml := received.ReplyCustom(func(r *wechat.ReplySt) (xml string) {
		xml, _ = wechat.FormatMap2XML(map[string]string{
			"Content":      "收到:" + r.MsgType + "|" + r.Content,
			"CreateTime":   strconv.FormatInt(time.Now().Unix(), 10),
			"ToUserName":   r.FromUserName,
			"FromUserName": r.ToUserName,
			"MsgType":      "text",
		})
		return
	})

	c.String(200, replyXml)
}
