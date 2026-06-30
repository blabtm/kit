package build

import (
	"fmt"

	"github.com/blabtm/kit/build"
	"github.com/blabtm/kit/build/dagger"
	"github.com/spf13/cobra"
)

var builder build.Builder

func init() {
	dag, err := dagger.New()

	if err != nil {
		panic(err)
	}

	builder = dag
}

var Cmd = &cobra.Command{
	Use:   "build {name}",
	Short: "Build the module.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mod, err := build.OpenModule(args[0])

		if err != nil {
			return fmt.Errorf("open: %w", err)
		}

		return builder.Build(mod)
	},
}
