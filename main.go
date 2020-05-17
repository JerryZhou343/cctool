package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/cmd"
)

var (
	Major = 1
	Minor = 0
	Patch = 1
)

func main() {
	cmd.RootCmd.Version = fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	cmd.RootCmd.Execute()
}
