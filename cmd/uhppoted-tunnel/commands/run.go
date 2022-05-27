package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/uhppoted/uhppoted-tunnel/tunnel"
)

type Run struct {
	console     bool
	in          string
	out         string
	debug       bool
	workdir     string
	logFile     string
	logFileSize int
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.in, "in", "", "IN connection e.g. udp/listen:0.0.0.0:60000, udp/broadcast:255.255.255.255:60000, tcp/listen:0.0.0.0:54321 or tcp/connect:101.102.103.104:54321")
	flagset.StringVar(&r.out, "out", "", "OUT connection e.g. udp/255.255.255.255:60000 or tcp/bind:0.0.0.0:54321 or tcp/connect:101.102.103.104:54321")
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
	return "uhppoted-tunnel [--debug] [--in <connection>] [--out <connection>] [--logfile <file>] [--logfilesize <bytes>] [--pid <file>]"
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

func (cmd *Run) execute(f func(t *tunnel.Tunnel)) error {
	var in tunnel.In
	var out tunnel.Out

	// ... create 'in' connection
	switch {
	case cmd.in == "":
		return fmt.Errorf("--in argument is required")

	case strings.HasPrefix(cmd.in, "udp/listen:"):
		if udp, err := tunnel.NewUDPIn(cmd.in[11:]); err != nil {
			return err
		} else {
			in = udp
		}

	case strings.HasPrefix(cmd.in, "tcp/connect:"):
		if tcp, err := tunnel.NewTCPIn(cmd.in[12:]); err != nil {
			return err
		} else {
			in = tcp
		}

	default:
		return fmt.Errorf("Invalid --in argument (%v)", cmd.in)
	}

	// ... create 'out' connection
	switch {
	case cmd.out == "":
		return fmt.Errorf("--out argument is required")

	case strings.HasPrefix(cmd.out, "udp/broadcast:"):
		if udp, err := tunnel.NewUDPOut(cmd.out[14:]); err != nil {
			return err
		} else {
			out = udp
		}

	case strings.HasPrefix(cmd.out, "tcp/listen:"):
		if tcp, err := tunnel.NewTCPOutHost(cmd.out[11:]); err != nil {
			return err
		} else {
			out = tcp
		}

	default:
		return fmt.Errorf("Invalid --out argument (%v)", cmd.out)
	}

	// // ... create lockfile
	// if err := os.MkdirAll(cmd.workdir, os.ModeDir|os.ModePerm); err != nil {
	// 	return fmt.Errorf("Unable to create working directory '%v': %v", cmd.workdir, err)
	// }
	//
	// pid := fmt.Sprintf("%d\n", os.Getpid())
	// lockfile := filepath.Join(cmd.workdir, fmt.Sprintf("%s.pid", SERVICE))
	//
	// if _, err := os.Stat(lockfile); err == nil {
	// 	return fmt.Errorf("PID lockfile '%v' already in use", lockfile)
	// } else if !os.IsNotExist(err) {
	// 	return fmt.Errorf("Error checking PID lockfile '%v' (%v)", lockfile, err)
	// }
	//
	// if err := os.WriteFile(lockfile, []byte(pid), 0644); err != nil {
	// 	return fmt.Errorf("Unable to create PID lockfile: %v", err)
	// }
	//
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Fatalf("%-5s %v\n", "FATAL", err)
	// 	}
	// }()
	//
	// defer os.Remove(lockfile)

	t := tunnel.NewTunnel(in, out)

	f(t)

	return nil
}

func (cmd *Run) run(t *tunnel.Tunnel, interrupt chan os.Signal) {

	t.Run(interrupt)
}
