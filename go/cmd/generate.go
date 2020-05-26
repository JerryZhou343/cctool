//从视频生成字幕

package main

import (
	"github.com/JerryZhou343/cctool/go/internal/app"
	"github.com/JerryZhou343/cctool/go/internal/console"
	"github.com/JerryZhou343/cctool/go/internal/flags"
	"github.com/JerryZhou343/cctool/go/internal/status"
	"github.com/spf13/cobra"
)

var (
	generateCmd = cobra.Command{
		Use:   "generate",
		Short: "生成字幕",
		Args:  cobra.OnlyValidArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(flags.SrcFiles) == 0 {
				err = status.ErrSourceFileNotEnough
				return
			}
			err = application.LoadSrtGenerator()
			if err != nil {
				return
			}

			application.Run()
			for _, itr := range flags.SrcFiles {
				task := app.NewGenerateTask(itr)
				application.AddTask(task)
			}

			console.Console(application)
			application.CheckTask()
			application.Destroy()
			return nil
		},
	}
)

func init() {
	//生成字幕
	generateCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "单个或多个源文件")
}
