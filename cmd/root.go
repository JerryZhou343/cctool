package cmd

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = cobra.Command{
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
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
