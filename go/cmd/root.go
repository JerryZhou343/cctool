package main

import (
	"github.com/JerryZhou343/cctool/go/internal/ui"
	"github.com/spf13/cobra"
)

var (
	RootCmd = cobra.Command{
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ui.UI()
			return
		},
	}
)

func init() {
	RootCmd.AddCommand(&translateCmd)
	RootCmd.AddCommand(&mergeCmd)
	RootCmd.AddCommand(&generateCmd)
	RootCmd.AddCommand(&convertCmd)
	RootCmd.AddCommand(&cleanCmd)
}
