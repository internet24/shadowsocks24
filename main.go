package main

import (
	"fmt"
	"github.com/internet24/shadowsocks24/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
