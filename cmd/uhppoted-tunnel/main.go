package main

import (
	"fmt"
	"os"

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
	cmd, err := lib.Parse(cli, &commands.RUN, help)
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Printf("\nERROR: %v\n\n", err)
		os.Exit(1)
	}
}
