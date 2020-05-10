//合并字幕文件

package cmd

import (
	"fmt"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/spf13/cobra"
	"log"
)

var (
	mergeCmd = cobra.Command{
		Use:   "merge",
		Short: "合并字幕",
		Args:  cobra.OnlyValidArgs,
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
)

func init() {
	//合并
	mergeCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "源文件")
	mergeCmd.PersistentFlags().StringVarP(&flags.DstFile, "destination", "d", "", "目标文件")
	mergeCmd.PersistentFlags().StringVar(&flags.MergeStrategy, "strategy", flags.StrategyTimeline,
		fmt.Sprintf("merge strategy:[%s:以第一个源文件的序号主,%s: 以第一个源文件的时间轴为主]",
			flags.StrategySequence, flags.StrategyTimeline))

}
