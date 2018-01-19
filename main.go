package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

var files = map[string][]string{
	"cli": []string{"i3-sensible-terminal", "-e"},
	"gui": []string{},
}

type config struct {
	path string
	args []string
}

func main() {
	cfg := config{}

	// Parse the configuration
	cfg.args = extractArgs()
	err := parseFlags(&cfg)
	check(err)

	var (
		cmdsStr = &bytes.Buffer{}
		cmdsIdx = map[string]string{}
	)

	// Read content files and add to command list
	for f := range files {
		s, err := ioutil.ReadFile(cfg.path + "/" + f)
		check(err)
		ls := strings.Split(string(s), "\n")

		// Add each one with type of file
		for _, l := range ls {
			if l != "" {
				cmdsStr.WriteString(l + "\n")
				cmdsIdx[l] = f
			}
		}
	}

	// Run the dmenu command
	dmenu := exec.Command("dmenu", cfg.args...)
	dmenu.Env = os.Environ()
	dmenu.Stdin = cmdsStr
	app, err := dmenu.Output()
	check(err)

	// Execute the chosen app
	fmt.Println(string(app))
}

// extractArgs returns a string of everything after a -- flag and
// removes them from os.Args.
func extractArgs() []string {
	for i, arg := range os.Args {
		if arg == "--" {
			defer func() {
				os.Args = os.Args[:i]
			}()
			return os.Args[i+1:]
		}
	}
	return nil
}

// parseFlags parses os.Args and returns a config object, along
// with a possible error.
func parseFlags(cfg *config) error {
	var tmp string // for help text, is a no-op
	flag.StringVar(&tmp, "-", "", "args to pass to dmenu")
	flag.StringVar(&cfg.path, "path", "~/.config/dmenual", "path to config dir")
	flag.Parse()

	// Expand the possible tilde in the path
	path, err := tilde.Expand(cfg.path)
	if err != nil {
		return err
	}
	cfg.path = path

	return nil
}

// check prints the error and exists if the error is non-nil.
func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
