package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/cmd"
	"github.com/JerryZhou343/cctool/internal/logger"
	"github.com/sirupsen/logrus"
)

var (
	Major = 1
	Minor = 1
	Patch = 4
)

func main() {
	logger.Init("debug", "log/", "cctool.log")
	cmd.RootCmd.Version = fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	logrus.Info("cctool version:", cmd.RootCmd.Version)
	cmd.RootCmd.Execute()
}
