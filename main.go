package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	cli = "cli"
	gui = "gui"
)

type config struct {
	path string
	args string
}

func main() {
	cfg := new(config)

	// Extract the -- flag first
	for i, arg := range os.Args {
		if arg == "--" {
			cfg.args = strings.Join(os.Args[i+1:], " ")
			os.Args = os.Args[:i]
			break
		}
	}

	// Parse the rest of the flags
	flag.StringVar(&cfg.path, "path", "~/.config/dmenual", "path to config dir")
	flag.Parse()

	fmt.Println(cfg)
}
