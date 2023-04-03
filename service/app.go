package service

import (
	"zlsapp/conf"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
)

// App 控制器关联对象
type App struct {
	Di   zdi.Injector
	Conf *Conf
	Log  *zlog.Logger
}

var Global *App

func RegApp(conf *Conf, di zdi.Injector) *App {
	Global = &App{
		Di:   di,
		Conf: conf,
		Log:  initLog(conf),
	}
	return Global
}

func initLog(c *Conf) *zlog.Logger {
	log := zlog.Log
	log.SetPrefix(conf.LogPrefix)
	logFlags := zlog.BitLevel
	if c.Base.LogPosition {
		logFlags = logFlags | zlog.BitLongFile
	}
	log.ResetFlags(logFlags)
	if c.Base.LogDir != "" {
		log.SetSaveFile(zfile.RealPath(c.Base.LogDir, true)+"app.log", true)
	}
	return log
}

func PrintLog(tip string, v ...interface{}) {
	d := []interface{}{
		zlog.ColorTextWrap(zlog.ColorLightMagenta, zstring.Pad(tip, 6, " ", zstring.PadLeft)),
	}
	d = append(d, v...)
	zlog.Tips(d...)
}
