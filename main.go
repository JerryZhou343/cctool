package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/internal/app"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/flags"
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
			Run: func(cmd *cobra.Command, args []string) {
				conf.Init()
				err := application.Translate()
				if err != nil {
					log.Printf("%+v", err)
				}
			},
		}
	)

	translateCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "源文件")
	translateCmd.PersistentFlags().StringVarP(&flags.From, "from", "f", "en", "源语言")
	translateCmd.PersistentFlags().StringVarP(&flags.To, "to", "t", "zh", "目标语言")
	translateCmd.PersistentFlags().StringVarP(&flags.TransTool, "transtool", "", "google",
		fmt.Sprintf("翻译工具: %s,%s,%s", flags.TransTool_Baidu, flags.TransTool_Google, flags.TransTool_Tencent))
	translateCmd.PersistentFlags().BoolVarP(&flags.Merge, "merge", "m", false, "双语字幕")

	rootCmd.AddCommand(&translateCmd)
	rootCmd.Execute()
}
