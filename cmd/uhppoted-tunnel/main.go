package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/pelletier/go-toml/v2"

	core "github.com/uhppoted/uhppote-core/uhppote"
	lib "github.com/uhppoted/uhppoted-lib/command"

	"github.com/uhppoted/uhppoted-tunnel/commands"
)

var cli = []lib.Command{
	&commands.DAEMONIZE,
	&commands.UNDAEMONIZE,
	&version,
}

var version = lib.Version{
	Application: commands.SERVICE,
	Version:     core.VERSION,
}

var help = lib.NewHelp(commands.SERVICE, cli, &commands.RUN)

func main() {
	var cmd lib.Command = &commands.RUN

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case commands.DAEMONIZE.Name():
			cmd = &commands.DAEMONIZE

		case commands.UNDAEMONIZE.Name():
			cmd = &commands.UNDAEMONIZE

		case version.Name():
			cmd = &version

		case help.Name():
			cmd = help
		}
	}

	flagset := cmd.FlagSet()
	if flagset == nil {
		panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", cmd.Name(), cmd))
	}

	flagset.Parse(args)

	config := map[string]any{}
	visited := map[string]bool{}

	flagset.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
		if f.Name == "config" {
			if c, err := configure(f.Value.String()); err != nil {
				fmt.Printf("\nERROR  %v\n\n", err)
				os.Exit(1)
			} else {
				config = c
			}
		}
	})

	flagset.VisitAll(func(f *flag.Flag) {
		if v, ok := config[f.Name]; ok && !visited[f.Name] {
			flagset.Set(f.Name, fmt.Sprintf("%v", v))
		}
	})

	if err := cmd.Execute(); err != nil {
		fmt.Printf("\nERROR: %v\n\n", err)
		os.Exit(1)
	}
}

func configure(configuration string) (map[string]any, error) {
	file := configuration
	section := ""
	if match := regexp.MustCompile("(.*?)(?:::|#)(.*)").FindStringSubmatch(configuration); match != nil {
		file = match[1]
		section = match[2]
	}

	config := map[string]any{}

	if bytes, err := os.ReadFile(file); err != nil {
		return nil, err
	} else {
		c := map[string]any{}
		if err := toml.Unmarshal(bytes, &c); err != nil {
			return nil, err
		}

		if m, ok := c["defaults"]; ok {
			if defaults, ok := m.(map[string]any); ok {
				for k, v := range defaults {
					config[k] = v
				}
			}
		}

		if m, ok := c[section]; ok {
			if tunnel, ok := m.(map[string]any); ok {
				for k, v := range tunnel {
					config[k] = v
				}
			}
		}
	}

	return config, nil
}
