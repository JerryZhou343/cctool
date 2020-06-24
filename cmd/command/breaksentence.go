package command

import (
	"github.com/JerryZhou343/cctool/internal/app"
	"github.com/JerryZhou343/cctool/internal/console"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/spf13/cobra"
)

var (
	BreakCmd = cobra.Command{
		Use:   "break",
		Short: "重新断句",
		Args:  cobra.OnlyValidArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(flags.SrcFiles) == 0 {
				err = status.ErrSourceFileNotEnough
				return
			}
			application.Run()
			for _, itr := range flags.SrcFiles {
				task := app.NewCleanTask(itr)
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
	BreakCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "单个或多个源文件")
}
