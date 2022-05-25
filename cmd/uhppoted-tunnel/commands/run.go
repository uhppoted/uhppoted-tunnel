package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/uhppoted/uhppoted-tunnel/log"
	"github.com/uhppoted/uhppoted-tunnel/tunnel"
)

type Run struct {
	console     bool
	debug       bool
	workdir     string
	logFile     string
	logFileSize int
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.BoolVar(&r.console, "console", false, "Runs as a console application rather than a service")
	flagset.BoolVar(&r.debug, "debug", false, "Enables detailed debugging logs")

	return flagset
}

func (cmd *Run) Name() string {
	return "run"
}

func (cmd *Run) Description() string {
	return "Runs the uhppoted-tunnel daemon/service until terminated by the system service manager"
}

func (cmd *Run) Usage() string {
	return "uhppoted-tunnel [--debug] [--config <file>] [--logfile <file>] [--logfilesize <bytes>] [--pid <file>]"
}

func (cmd *Run) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted-tunnel <options>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	cmd.FlagSet().VisitAll(func(f *flag.Flag) {
		fmt.Printf("    --%-12s %s\n", f.Name, f.Usage)
	})
	fmt.Println()
}

func (cmd *Run) execute(f func()) error {
	// ... create lockfile
	if err := os.MkdirAll(cmd.workdir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", cmd.workdir, err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())
	lockfile := filepath.Join(cmd.workdir, fmt.Sprintf("%s.pid", SERVICE))

	if _, err := os.Stat(lockfile); err == nil {
		return fmt.Errorf("PID lockfile '%v' already in use", lockfile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error checking PID lockfile '%v' (%v)", lockfile, err)
	}

	if err := os.WriteFile(lockfile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create PID lockfile: %v", err)
	}

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("%-5s %v\n", "FATAL", err)
		}
	}()

	defer os.Remove(lockfile)

	f()

	return nil
}

func (cmd *Run) run(interrupt chan os.Signal) {
	t := tunnel.Tunnel{}

	t.Run(interrupt)
}
