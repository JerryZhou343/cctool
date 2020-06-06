package ui

import (
	"github.com/JerryZhou343/cctool/go/internal/ui/plugin/file_picker"
	"github.com/JerryZhou343/cctool/go/internal/ui/plugin/version"
	"github.com/go-flutter-desktop/go-flutter"
)

var options = []flutter.Option{
	flutter.WindowInitialDimensions(600, 800),
	flutter.AddPlugin(version.VersionPlugin{}),
	flutter.AddPlugin(&file_picker.FilePickerPlugin{}),
}
