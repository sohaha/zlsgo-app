package service

import (
	"zlsapp/conf"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
)

var Global *App

func InitApp(conf *Conf, di zdi.Injector, wechat *Wechat) *App {
	Global = &App{
		Di:     di,
		Conf:   conf,
		Wechat: wechat,
		Log:    initLog(conf),
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
