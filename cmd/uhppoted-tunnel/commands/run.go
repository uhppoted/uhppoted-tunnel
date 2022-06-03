package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/tunnel"
)

type Run struct {
	console       bool
	udp           string
	pipe          string
	maxRetries    int
	maxRetryDelay time.Duration
	debug         bool
	workdir       string
	logFile       string
	logFileSize   int
}

const MAX_RETRIES = -1
const MAX_RETRY_DELAY = 5 * time.Minute

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.udp, "udp", "", "UDP connection e.g. listen:0.0.0.0:60000 or broadcast:255.255.255.255:60000")
	flagset.StringVar(&r.pipe, "pipe", "", "TCP pipe connection e.g. tcp/server:0.0.0.0:54321 or tcp/client:101.102.103.104:54321")

	flagset.IntVar(&r.maxRetries, "max-retries", MAX_RETRIES, "Maximum number of times to retry failed connection. Defaults to -1 (retry forever)")
	flagset.DurationVar(&r.maxRetryDelay, "max-retry-delay", MAX_RETRY_DELAY, "Maximum delay between retrying failed connections")
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
	var udp tunnel.UDP
	var pipe tunnel.TCP
	var mode = tunnel.ModeNormal

	if strings.HasPrefix(cmd.udp, "broadcast:") && strings.HasPrefix(cmd.pipe, "tcp/server:") {
		mode = tunnel.ModeReverse
	}

	if strings.HasPrefix(cmd.udp, "listen:") && strings.HasPrefix(cmd.pipe, "tcp/client:") {
		mode = tunnel.ModeReverse
	}

	// ... create UDP packet handler
	switch {
	case cmd.udp == "":
		return fmt.Errorf("--udp argument is required")

	case strings.HasPrefix(cmd.udp, "listen:"):
		if u, err := tunnel.NewUDPListen(cmd.udp[7:]); err != nil {
			return err
		} else {
			udp = u
		}

	case strings.HasPrefix(cmd.udp, "broadcast:"):
		if u, err := tunnel.NewUDPBroadcast(cmd.udp[10:]); err != nil {
			return err
		} else {
			udp = u
		}

	default:
		return fmt.Errorf("Invalid --udp argument (%v)", cmd.udp)
	}

	// ... create TCP/IP pipe
	switch {
	case cmd.pipe == "":
		return fmt.Errorf("--pipe argument is required")

	case strings.HasPrefix(cmd.pipe, "tcp/client:"):
		if tcp, err := tunnel.NewTCPClient(cmd.pipe[11:], cmd.maxRetries, cmd.maxRetryDelay, mode); err != nil {
			return err
		} else {
			pipe = tcp
		}

	case strings.HasPrefix(cmd.pipe, "tcp/server:"):
		if tcp, err := tunnel.NewTCPServer(cmd.pipe[11:], mode); err != nil {
			return err
		} else {
			pipe = tcp
		}

	default:
		return fmt.Errorf("Invalid --pipe argument (%v)", cmd.pipe)
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

	t := tunnel.NewTunnel(udp, pipe)

	f(t)

	return nil
}

func (cmd *Run) run(t *tunnel.Tunnel, interrupt chan os.Signal) {
	t.Run(interrupt)
}
