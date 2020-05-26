package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/go/internal/version"
)



func main() {
	RootCmd.Version = fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)
	RootCmd.Execute()
}
