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
	&lib.Version{
		Application: commands.SERVICE,
		Version:     core.VERSION,
	},
}

var help = lib.NewHelp(commands.SERVICE, cli, &commands.RUN)

func main() {
	conf := flag.String("config", "", "(optional) tunnel TOML configuration file")

	flag.Parse()

	if conf != nil && *conf != "" {
		if config, err := configure(*conf); err != nil {
			fmt.Printf("\nERROR  %v\n\n", err)
			os.Exit(1)
		} else if config != nil {
			fmt.Printf(">>>>>> %v\n", config)
		}
	}

	cmd, err := lib.Parse(cli, &commands.RUN, help)
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Printf("\nERROR: %v\n\n", err)
		os.Exit(1)
	}
}

func configure(configuration string) (map[string]any, error) {
	file := configuration
	section := ""
	if match := regexp.MustCompile("(.*?)::(.*)").FindStringSubmatch(configuration); match != nil {
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
