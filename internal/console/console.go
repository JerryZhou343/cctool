package console

import (
	"github.com/JerryZhou343/cctool/internal/app"
	"github.com/fatih/color"
)

func Show(application *app.Application) {
	for {
		//fmt.Printf("%s\n", strings.Trim(application.GetRunningMsg(), "\n"))
		task := application.GetRunningMsg()
		if task.GetState() == app.TaskStateFailed{
			color.Red("%s",task)
		}else{
			color.Green("%s",task)
		}
	}
}

func Console(application *app.Application) {
	go Show(application)

	/*
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			s := <-c
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	*/
}
