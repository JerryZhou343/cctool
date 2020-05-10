//从视频生成字幕

package cmd

import (
	"github.com/JerryZhou343/cctool/internal/conf"
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/status"
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

func init() {
	//生成字幕
	generateCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "源文件")
	generateCmd.PersistentFlags().IntVarP(&flags.AudioChannelId, "channel", "c", 0, "音频声道")
}
