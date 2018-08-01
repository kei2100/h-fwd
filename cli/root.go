package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd for CLI
var RootCmd = &cobra.Command{
	Use:   "fwxy",
	Short: "fwxy is a simple HTTP forward proxy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test")
	},
}
