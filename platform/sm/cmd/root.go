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

var config string
var remoteHost string
var remotePort int

var Root = &cobra.Command{
	Use:   "sm",
	Short: "Service manager for the VEPP-2000 Collider platform.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(cmd)
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
			return err
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
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := v1.NewClient(
			viper.GetString("remote.host"),
			viper.GetInt("remote.port"),
		).Ps(args[0])

		if err != nil {
			return err
		}

		fmt.Printf("%s:\n", args[0])
		fmt.Printf("\tState: %s\n", status.State)

		return nil
	},
}

var Get = &cobra.Command{
	Use:   "get [name]",
	Short: "Get runtime configuration of the service.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, err := v1.NewClient(
			viper.GetString("remote.host"),
			viper.GetInt("remote.port"),
		).Get(args[0])

		if err != nil {
			return err
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
			return err
		}

		fmt.Println("OK")

		return nil
	},
}

func init() {
	Root.AddCommand(Ls)
	Root.AddCommand(Ps)
	Root.AddCommand(Get)
	Root.AddCommand(Set)

	Root.PersistentFlags().StringVarP(&config, "config", "c", "", "configuration file (default locations are: $HOME/.sm/config.yaml, /etc/sm/config.yaml)")
	Root.PersistentFlags().StringVar(&remoteHost, "remote.host", "", "remote host")
	Root.PersistentFlags().IntVar(&remotePort, "remote.port", 0, "remote port")
}

func initConfig(cmd *cobra.Command) error {
	if config != "" {
		viper.SetConfigFile(config)
	} else {
		home, err := os.UserHomeDir()

		if err != nil {
			return err
		}

		viper.AddConfigPath(home + "/.sm")
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
