package conf

import "github.com/zlsgo/wechat"

type (
	Wechat struct {
		Debug bool
		Mp    wechat.Mp
		Open  wechat.Open
		Qy    wechat.Qy
		Weapp wechat.Weapp
		Pay   wechatPay
	}

	wechatPay struct {
		Sandbox  bool
		AppID    string
		MchId    string
		Key      string
		CertPath string
		KeyPath  string
	}
)

func init() {
	DefaultSet = append(DefaultSet, Wechat{
		Debug: false,
		Mp: wechat.Mp{
			AppID:     "",
			AppSecret: "",
		},
		Qy: wechat.Qy{
			CorpID:         "",
			AgentID:        "",
			Secret:         "",
			Token:          "",
			EncodingAesKey: "",
		},
		Open: wechat.Open{
			AppID:          "",
			AppSecret:      "",
			EncodingAesKey: "",
			Token:          "",
		},
		Weapp: wechat.Weapp{
			AppID:          "",
			AppSecret:      "",
			EncodingAesKey: "",
		},
		Pay: wechatPay{
			Sandbox:  false,
			AppID:    "",
			MchId:    "",
			Key:      "",
			CertPath: "",
			KeyPath:  "",
		},
	})
}
