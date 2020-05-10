package main

import (
	"fmt"
	"github.com/JerryZhou343/cctool/internal/app"
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	var (
		application = app.NewApplication()
		rootCmd     = cobra.Command{
			Args: cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				return
			},
		}
		translateCmd = cobra.Command{
			Use:   "translate",
			Short: "翻译字幕",
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				if len(flags.SrcFiles) == 0 {
					err = status.ErrSourceFileNotEnough
					return
				}
				err = application.LoadTranslateTools()
				if err != nil {
					log.Printf("%+v", err)
				}

				application.Run()
				fmt.Println("s1")
				for _, itr := range flags.SrcFiles {
					task := app.NewTranslateTask(itr, flags.From, flags.To, flags.Merge)
					application.AddTranslateTask(task)
				}

				Console(application)
				application.Destroy()
				return nil
			},
		}

		mergeCmd = cobra.Command{
			Use:   "merge",
			Short: "合并字幕",
			RunE: func(cmd *cobra.Command, args []string) error {
				//check
				if len(flags.SrcFiles) != 2 {
					return status.ErrSourceFileNotEnough
				}
				if flags.DstFile == "" {
					return status.ErrDstFile
				}

				err := application.Merge()
				if err != nil {
					log.Printf("%+v", err)
				}
				return nil
			},
		}

		generateCmd = cobra.Command{
			Use:   "generate",
			Short: "生成字幕",
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				if len(flags.SrcFiles) == 0 {
					err = status.ErrSourceFileNotEnough
					return
				}
				err = conf.Load()
				if err != nil {
					err = status.ErrInitConfigFileFailed
					return
				}

				for _, itr := range flags.SrcFiles {
					err = application.GenerateSrt(itr, flags.AudioChannelId)
					if err != nil {
						return
					}
				}
				return nil
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
	mergeCmd.PersistentFlags().StringVar(&flags.MergeStrategy, "strategy", flags.StrategyTimeline,
		fmt.Sprintf("merge strategy：[%s:以第一个源文件的序号主,%s: 以第一个源文件的时间轴为主]",
			flags.StrategySequence, flags.StrategyTimeline))

	//生成字幕
	generateCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "源文件")
	generateCmd.PersistentFlags().IntVarP(&flags.AudioChannelId, "channel", "c", 0, "音频声道")

	rootCmd.AddCommand(&translateCmd)
	rootCmd.AddCommand(&mergeCmd)
	rootCmd.AddCommand(&generateCmd)
	rootCmd.Execute()

}

func Show(application *app.Application) {
	for {
		fmt.Printf("%s\n", strings.Trim(application.GetRunningMsg(), "\n"))
	}
}

func Console(application *app.Application) {
	go Show(application)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
