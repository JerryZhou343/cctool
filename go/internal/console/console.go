package console

import (
	"fmt"
	"github.com/JerryZhou343/cctool/go/internal/app"
	"strings"
)

func Show(application *app.Application) {
	for {
		fmt.Printf("%s\n", strings.Trim(application.GetRunningMsg(), "\n"))
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
