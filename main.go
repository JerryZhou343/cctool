package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/cmd"
)

var (
	Major = 0
	Minor = 1
	Patch = 0
)

func main() {
	cmd.RootCmd.Version = fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	cmd.RootCmd.Execute()
}
