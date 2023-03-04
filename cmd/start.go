package cmd

import (
	"github.com/internet24/shadowsocks24/internal/app"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start",
	Run: startFunc,
}

func startFunc(_ *cobra.Command, _ []string) {
	a, err := app.New(configPath)
	if err != nil {
		panic(err)
	}
	defer a.Shutdown()

	go a.Coordinator.Run()
	go a.HttpServer.Run()

	a.Wait()
}
