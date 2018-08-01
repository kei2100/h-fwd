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
		fmt.Println(rewritePaths)
		fmt.Println(len(rewritePaths))
	},
}

var (
	username     string
	password     string
	rewritePaths []string
)

func init() {
	flags := RootCmd.PersistentFlags()

	flags.StringVarP(&username, "username", "u", "", "username for the basic authentication")
	flags.StringVarP(&password, "password", "p", "", "password for the basic authentication")
	flags.StringSliceVarP(&rewritePaths, "rewrite", "r", []string{}, "list for path rewrite (-r /old:/new -r /o:/n OR -r /old:/new,/o:/n)")
}
