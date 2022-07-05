package wechat

import (
	"strconv"
	"sync"
	"time"

	"github.com/sohaha/zlsgo/zstring"
	"zlsapp/logic/code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/znet"

	"github.com/zlsgo/wechat"
)

type Open struct {
	service.App
}

func (w *Open) Init(_ *znet.Engine) {
}

// AnyNotification 开放平台推送消息授权事件接收
func (w *Open) AnyNotification(c *znet.Context) {
	wx := w.App.Wechat.Open

	data, _ := c.GetDataRaw()
	_, err := wx.ComponentVerifyTicket(data)
	if err != nil {
		c.Log.Warn(err.Error())
	}
	// 需要返回 success 给微信
	c.String(200, "success")
}

// FullAnyReceiveMessage 接收消息与事件
func (w *Open) FullAnyReceiveMessage(c *znet.Context) {
	wx := w.App.Wechat.Open

	body, _ := c.GetDataRaw()
	reply, err := wx.Reply(c.GetAllQueryMaps(), zstring.String2Bytes(body))
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

// GetApiQueryAuth 公众号授权
func (w *Open) GetApiQueryAuth(c *znet.Context) {
	wx := w.App.Wechat.Open

	authCode, ok := c.GetQuery("auth_code")
	if !ok {
		c.Log.Debug("需要跳转")
	}

	// TODO 本地使用内网穿透做测试，所以需要固定回调域名
	wx.SetOptions(wechat.WithRedirectDomain("http://mac.hw.73zls.com/"))

	res, redirect, err := wx.ComponentApiQueryAuth(c, authCode)
	if err != nil {
		if err != wechat.ErrOpenJumpAuthorization {
			code.InvalidInput.ApiResult(c, wechat.ErrorMsg(err))
			return
		}
		// JS 发起跳转
		c.Template(200, `<html lang='zh'><head><title>Loading...
</title></head><body><script type='text/javascript'>
referLink=document.createElement('a');referLink.href="{{.redirect}}";referLink.
click()</script></body></html>`, map[string]string{"redirect": redirect})
		return
	}
	code.Success.ApiResult(c, zjson.Parse(res).Value())
}

func (w *Open) GetTicket(c *znet.Context) {
	wx := w.App.Wechat.Open

	ticket, err := wx.GetConfig().(*wechat.Open).GetComponentTicket()

	if err != nil {
		code.ServerError.ApiResult(c, wechat.ErrorMsg(err))
		return
	}

	code.Success.ApiResult(c, ticket)
}

var o sync.Once

func (w *Open) GetAccessToken(c *znet.Context) {
	wx := w.App.Wechat.Open

	o.Do(func() {
		// todo 因为开放平台的是授权的时候获取的，这里可以强制手动设置
		// wx.GetConfig().(*wechat.Open).SetAuthorizerAccessToken("", "", "", 4)
	})

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

func (w *Open) GetJsapiTicket(c *znet.Context) {
	wx := w.App.Wechat.Open

	jsapiTicket, err := wx.GetJsapiTicket()
	if err != nil {
		code.ServerError.ApiResult(c, err)
		return
	}

	url := c.Host(true)
	jsSign, err := wx.GetJsSign(url)
	if err != nil {
		code.ServerError.ApiResult(c, err)
		return
	}

	code.Success.ApiResult(c, map[string]interface{}{
		"jsapiTicket": jsapiTicket,
		"jsSign":      jsSign,
		"url":         url,
	})
}
