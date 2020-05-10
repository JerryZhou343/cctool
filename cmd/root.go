package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	RootCmd = cobra.Command{
		Args:    cobra.NoArgs,
		Version: fmt.Sprintf("v0.0.0"),
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
}
