package command

import (
	"github.com/sohaha/zlsgo/zcli"
)

func Reg() {
	zcli.SetLangText("zh", "init", "生成配置")
	zcli.Add("init", zcli.GetLangText("init", "Init config file"), &InitCli{})
}
