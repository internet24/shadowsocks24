package cmd

import (
	"fmt"
	"github.com/internet24/shadowsocks24/internal/config"
	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use: "shadowsocks24",
}

func init() {
	cobra.OnInitialize(func() { fmt.Println(config.AppName) })

	rootCmd.PersistentFlags().StringVarP(
		&configPath, "config", "c", "configs/config.json", "Config file path",
	)

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
