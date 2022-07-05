package wechat

import (
	"zlsapp/logic/code"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/zlsgo/wechat"
)

func (w *Mp) GetPay(c *znet.Context) {
	openid := "oj3r9s0kqWLyZJRQ2FTim7QnYmGI"
	order := wechat.NewPayOrder(openid, 1, c.GetClientIP(), "商品")
	notifyUrl := "https://mac.hw.73zls.com/wechat/pay/notify"

	prepayID, err := w.Wechat.Pay.UnifiedOrder(w.Conf.Wechat.Pay.AppID, order, notifyUrl)
	if err != nil {
		code.InvalidInput.ApiResult(c, err)
		return
	}

	sign := w.Wechat.Pay.JsSign(w.Conf.Wechat.Pay.AppID, prepayID)

	json, _ := zjson.Marshal(sign)
	c.Template(200, html, map[string]interface{}{
		"sign": zstring.Bytes2String(json),
	})
}

const html = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Pay</title>
  </head>
  <body>
  {{.sign}}
    <script>
      var conf = "{{.sign}}";
      function onBridgeReady() {
        WeixinJSBridge.invoke(
          "getBrandWCPayRequest",
          JSON.parse(conf),
          function (res) {
            alert(res.err_msg);
            if (res.err_msg == "get_brand_wcpay_request:ok") {
              // 使用以上方式判断前端返回,微信团队郑重提示：
              //res.err_msg将在用户支付成功后返回ok，但并不保证它绝对可靠。
            }
          }
        );
      }
      if (typeof WeixinJSBridge == "undefined") {
        if (document.addEventListener) {
          document.addEventListener(
            "WeixinJSBridgeReady",
            onBridgeReady,
            false
          );
        } else if (document.attachEvent) {
          document.attachEvent("WeixinJSBridgeReady", onBridgeReady);
          document.attachEvent("onWeixinJSBridgeReady", onBridgeReady);
        }
      } else {
        onBridgeReady();
      }
    </script>
  </body>
</html>
`
