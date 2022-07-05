package code

import (
	"github.com/sohaha/zlsgo/znet"
)

type ApiData struct {
	Code ErrCode     `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func (code ErrCode) ApiResult(c *znet.Context, data interface{}) {
	var d interface{} = struct{}{}

	if code == Success {
		if data != nil {
			d = data
		}
		c.JSON(200, ApiData{code, "", d})
		return
	}

	m := Text(code)
	switch v := data.(type) {
	case error:
		m = v.Error()
	case string:
		if len(v) > 0 {
			m = v
		}
	default:
		m = Text(code)
	}

	c.JSON(200, ApiData{code, m, struct{}{}})
}

func Text(code ErrCode) string {
	msg, ok := zhCNText[code]
	if !ok {
		return "no error code defined"
	}
	return msg
}
