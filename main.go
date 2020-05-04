package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/internal/app"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var (
		application  = app.NewApplication()
		rootCmd      = cobra.Command{}
		translateCmd = cobra.Command{
			Use:   "translate",
			Short: "翻译字幕",
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				if len(flags.SrcFiles) == 0 {
					err = status.ErrSourceFileNotEnough
					return
				}
				err = conf.Init()
				if err != nil {
					err = status.ErrInitConfigFileFailed
					return
				}
				err = application.Translate()
				if err != nil {
					log.Printf("%+v", err)
				}

				return nil
			},
		}

		mergeCmd = cobra.Command{
			Use:   "merge",
			Short: "合并字幕",
			Run: func(cmd *cobra.Command, args []string) {
				//check
				if len(flags.SrcFiles) != 2 {
					return
				}
				if flags.DstFile == "" {
					return
				}

				err := application.Merge()
				if err != nil {
					log.Printf("%+v", err)
				}
			},
		}
	)

	//翻译
	translateCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "源文件")
	translateCmd.PersistentFlags().StringVarP(&flags.From, "from", "f", "en", "源语言")
	translateCmd.PersistentFlags().StringVarP(&flags.To, "to", "t", "zh", "目标语言")
	translateCmd.PersistentFlags().StringVarP(&flags.TransTool, "transtool", "", "google",
		fmt.Sprintf("翻译工具: %s,%s,%s", flags.TransTool_Baidu, flags.TransTool_Google, flags.TransTool_Tencent))
	translateCmd.PersistentFlags().BoolVarP(&flags.Merge, "merge", "m", false, "双语字幕")

	//合并
	mergeCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "源文件")
	mergeCmd.PersistentFlags().StringVarP(&flags.DstFile, "destination", "d", "", "目标文件")
	mergeCmd.PersistentFlags().StringVar(&flags.MergeStrategy, "strategy", flags.StrategySequence,
		fmt.Sprintf("merge strategy：[%s:以第一个源文件的序号主,%s: 以第一个源文件的时间轴为主]",
			flags.StrategySequence, flags.StrategyTimeline))

	rootCmd.AddCommand(&translateCmd)
	rootCmd.AddCommand(&mergeCmd)

	rootCmd.Execute()
}
