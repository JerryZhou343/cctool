package command

import (
	"github.com/JerryZhou343/cctool/internal/flags"
	"github.com/JerryZhou343/cctool/internal/status"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	BreakCmd = cobra.Command{
		Use:   "break",
		Short: "重新断句",
		Args:  cobra.OnlyValidArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			//check
			if len(flags.SrcFiles) != 1 {
				return status.ErrSourceFileNotEnough
			}
			if flags.DstFile == "" {
				return status.ErrDstFile
			}

			err = application.BreakSentence()
			if err != nil{
				color.Red("%+v",err)
			}
			return nil
		},
	}
)

func init() {
	BreakCmd.PersistentFlags().StringSliceVarP(&flags.SrcFiles, "source", "s", []string{}, "单个或多个源文件")
	BreakCmd.PersistentFlags().StringVarP(&flags.DstFile, "destination", "d", "", "目标文件")
}
