package service

import (
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztime/cron"
)

type Task struct {
	Name string
	Cron string
	Run  func(app *App)
}

func InitTask(tasks *[]Task, app *App) (err error) {
	t := cron.New()

	for i := range *tasks {
		task := &(*tasks)[i]
		if task.Cron == "" || task.Run == nil {
			continue
		}
		_, err = t.Add(task.Cron, func() {
			task.Run(app)
		})

		if err != nil {
			return
		}

		next, _ := cron.ParseNextTime(task.Cron)
		PrintLog("Cron", "Register: "+zlog.Log.ColorTextWrap(zlog.ColorLightGreen, task.Name)+zlog.ColorTextWrap(zlog.ColorLightWhite, " ["+task.Cron+"] -> ["+ztime.FormatTime(next)+"]"))
	}

	t.Run()
	return nil
}
