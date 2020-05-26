package version


import (
	"fmt"
	"github.com/JerryZhou343/cctool/go/internal/version"
	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
)

const (
	channelName = "cctool/go/version"
	getVersion  = "getVersion"
)

type VersionPlugin struct{}

var _ flutter.Plugin = &VersionPlugin{}

func (VersionPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc(getVersion, getVersionFunc)
	return nil;
}

func getVersionFunc(arguments interface{}) (reply interface{}, err error) {
	return fmt.Sprintf("v%d.%d.%d",version.Major,version.Minor,version.Patch), nil
}
