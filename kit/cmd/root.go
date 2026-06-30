package cmd

import (
	"fmt"
	"os"

	"github.com/blabtm/kit/cmd/build"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kit",
	Short: "Collider management toolkit.",
}

func init() {
	rootCmd.AddCommand(build.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
