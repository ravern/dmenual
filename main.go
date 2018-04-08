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

var files = map[string]func(config, string){
	"cli": func(cfg config, arg string) {
		cmd := exec.Command(cfg.term, "-e", arg)
		err := cmd.Run()
		check(err)
	},
	"gui": func(cfg config, arg string) {
		cmd := exec.Command(arg)
		err := cmd.Run()
		check(err)
	},
}

type config struct {
	path string
	term string
	args []string
}

func main() {
	cfg := config{}

	// Parse the configuration
	cfg.args = extractArgs()
	err := parseFlags(&cfg)
	check(err)

	var (
		cmdsStr = &bytes.Buffer{}     // pass into dmenu
		cmdsLbl = map[string]string{} // label to actual name
		cmdsIdx = map[string]string{} // which type of launch to use
	)

	// Read content files and add to command list
	for f := range files {
		s, err := ioutil.ReadFile(cfg.path + "/" + f)
		check(err)
		ls := strings.Split(string(s), "\n")

		// Add each one with type of launch
		for _, l := range ls {
			ll := strings.Split(l, ":")
			if len(ll) == 2 {
				ll[0] = strings.TrimSpace(ll[0])
				ll[1] = strings.TrimSpace(ll[1])
				cmdsStr.WriteString(ll[0] + "\n")
				cmdsLbl[ll[0]] = ll[1]
				cmdsIdx[ll[1]] = f
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
	appStr := cmdsLbl[strings.TrimSpace(string(app))]
	files[cmdsIdx[appStr]](cfg, appStr)
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
	flag.StringVar(&cfg.term, "term", "i3-sensible-terminal", "terminal to execute cli")
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
