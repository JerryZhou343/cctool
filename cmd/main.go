package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/cmd/command"
	"github.com/JerryZhou343/cctool/internal/logger"
	"github.com/sirupsen/logrus"
)

var (
	Major = 1
	Minor = 3
	Patch = 0
)

func main() {
	logger.Init("debug", "log/", "cctool.log")
	command.RootCmd.Version = fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	logrus.Info("cctool version:", command.RootCmd.Version)
	command.RootCmd.Execute()
}
