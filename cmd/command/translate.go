//翻译字幕文件

package command

import (
	"github.com/JerryZhou343/cctool/internal/app"
	"github.com/JerryZhou343/cctool/internal/console"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/spf13/cobra"
)

var (
	translateCmd = cobra.Command{
		Use:   "translate",
		Short: "翻译字幕",
		Args:  cobra.OnlyValidArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(flags.SrcFiles) == 0 {
				err = status.ErrSourceFileNotEnough
				return
			}
			err = application.LoadTranslateTools()
			if err != nil {
				return
			}

			application.Run()
			for _, itr := range flags.SrcFiles {
				task := app.NewTranslateTask(itr, flags.From, flags.To, flags.Merge)
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
	//翻译
	translateCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "源文件")
	translateCmd.PersistentFlags().StringVarP(&flags.From, "from", "f", "en", "源语言")
	translateCmd.PersistentFlags().StringVarP(&flags.To, "to", "t", "zh", "目标语言")
	translateCmd.PersistentFlags().BoolVarP(&flags.Merge, "merge", "m", false, "双语字幕")
}
