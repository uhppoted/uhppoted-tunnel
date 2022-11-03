package main

import (
	_ "embed"
	"fmt"
	"os"

	// "github.com/pelletier/go-toml/v2"

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

// //go:embed uhppoted-tunnel.toml
// var configuration []byte

func main() {
	// config := map[string]any{}
	// if err := toml.Unmarshal(configuration, &config); err != nil {
	// 	fmt.Printf(">>> Error unmarshalling TOML configuration (%v)\n", err)
	// } else {
	// 	fmt.Printf(">>> TOML: %v\n", config)
	// }

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
