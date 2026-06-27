package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	v1 "github.com/blabtm/v2k/platform/sm/rest/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

var Tag string

var config string
var remoteHost string
var remotePort int
var message string

var Root = &cobra.Command{
	Use:   "sm",
	Short: "Service manager for the VEPP-2000 collider platform.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(cmd)
	},
}

var Version = &cobra.Command{
	Use:   "version",
	Short: "Print version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sm version %s\n", Tag)
	},
}

var Ls = &cobra.Command{
	Use:   "ls",
	Short: "List registered services.",
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := v1.NewClient(
			viper.GetString("remote.host"),
			viper.GetInt("remote.port"),
		).Ls()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}

		for _, s := range services {
			fmt.Printf("* %s\n", s)
		}

		return nil
	},
}

var Ps = &cobra.Command{
	Use:   "ps [name]",
	Short: "Get the current operational status of the service.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := ""

		if len(args) == 2 {
			id = args[1]
		}

		status, err := v1.NewClient(
			viper.GetString("remote.host"),
			viper.GetInt("remote.port"),
		).Ps(args[0], id)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}

		fmt.Printf("%s:\n", args[0])
		fmt.Printf("  State: %s\n", status.State)

		return nil
	},
}

var Get = &cobra.Command{
	Use:   "get [name] {id}",
	Short: "Get runtime configuration of the service.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := ""

		if len(args) == 2 {
			id = args[1]
		}

		raw, err := v1.NewClient(
			viper.GetString("remote.host"),
			viper.GetInt("remote.port"),
		).Get(args[0], id)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}

		yml, err := yaml.JSONToYAML(raw)
		if err != nil {
			return err
		}

		fmt.Println(string(yml))

		return nil
	},
}

var Set = &cobra.Command{
	Use:   "set [name]",
	Short: "Set runtime configuration of the service.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		err = v1.NewClient(
			viper.GetString("remote.host"),
			viper.GetInt("remote.port"),
		).Set(args[0], raw)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}

		fmt.Println("OK")

		return nil
	},
}

var Fork = &cobra.Command{
	Use:   "fork [name]",
	Short: "Create new oneshot instance.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		id, err := v1.NewClient(
			viper.GetString("remote.host"),
			viper.GetInt("remote.port"),
		).Fork(args[0], raw, message)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}

		fmt.Printf("OK: %s\n", id)

		return nil
	},
}

func init() {
	Root.AddCommand(Version)
	Root.AddCommand(Ls)
	Root.AddCommand(Ps)
	Root.AddCommand(Get)
	Root.AddCommand(Set)
	Root.AddCommand(Fork)

	Root.PersistentFlags().StringVarP(&config, "config", "c", "", "configuration file (default locations are: $HOME/.config/sm/config.yaml, /etc/sm/config.yaml)")
	Root.PersistentFlags().StringVar(&remoteHost, "remote.host", "", "remote host")
	Root.PersistentFlags().IntVar(&remotePort, "remote.port", 0, "remote port")

	Fork.Flags().StringVarP(&message, "message", "m", "", "commit message")
}

func initConfig(cmd *cobra.Command) error {
	if config != "" {
		viper.SetConfigFile(config)
	} else {
		home, err := os.UserHomeDir()

		if err != nil {
			return err
		}

		viper.AddConfigPath(home + "/.config/sm")
		viper.AddConfigPath("/etc/sm")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return err
		}
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	return nil
}
