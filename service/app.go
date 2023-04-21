package service

import (
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
)

// App 控制器关联对象
type App struct {
	DI   zdi.Injector
	Conf *Conf
	Log  *zlog.Logger
}

var Global *App

func NewApp(conf *Conf, di zdi.Injector) *App {
	Global = &App{
		DI:   di,
		Conf: conf,
		Log:  initLog(conf),
	}
	return Global
}

func initLog(c *Conf) *zlog.Logger {
	log := zlog.Log
	log.SetPrefix(LogPrefix)

	logFlags := zlog.BitLevel | zlog.BitTime
	if c.Base.LogPosition {
		logFlags = logFlags | zlog.BitLongFile
	}
	log.ResetFlags(logFlags)

	if c.Base.LogDir != "" {
		log.SetSaveFile(zfile.RealPath(c.Base.LogDir, true)+"app.log", true)
	}

	if c.Base.Debug {
		log.SetLogLevel(zlog.LogDump)
	} else {
		log.SetLogLevel(zlog.LogSuccess)
	}

	return log
}

func PrintLog(tip string, v ...interface{}) {
	d := []interface{}{
		zlog.ColorTextWrap(zlog.ColorLightMagenta, zstring.Pad(tip, 6, " ", zstring.PadLeft)),
	}
	d = append(d, v...)
	zlog.Debug(d...)
}
