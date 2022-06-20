package commands

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/log"
	"github.com/uhppoted/uhppoted-tunnel/tunnel"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/tcp"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/tls"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/udp"
)

type Run struct {
	label         string
	portal        string
	pipe          string
	maxRetries    int
	maxRetryDelay time.Duration
	udpTimeout    time.Duration
	lockfile      string
	logFile       string
	logFileSize   int
	logLevel      string
	workdir       string
	debug         bool
	console       bool
}

const MAX_RETRIES = -1
const MAX_RETRY_DELAY = 5 * time.Minute
const UDP_TIMEOUT = 5 * time.Second
const CACERT = "ca.cert"

func (r *Run) flags() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.portal, "portal", "", "UDP connection e.g. udp/listen:0.0.0.0:60000 or udp/broadcast:255.255.255.255:60000")
	flagset.StringVar(&r.pipe, "pipe", "", "TCP pipe connection e.g. tcp/server:0.0.0.0:54321 or tcp/client:101.102.103.104:54321")
	flagset.StringVar(&r.lockfile, "lockfile", "", "(optional) name of lockfile used to prevent running multiple copies of the service. A default lockfile name is generated if none is supplied")
	flagset.IntVar(&r.maxRetries, "max-retries", MAX_RETRIES, "Maximum number of times to retry failed connection. Defaults to -1 (retry forever)")
	flagset.DurationVar(&r.maxRetryDelay, "max-retry-delay", MAX_RETRY_DELAY, "Maximum delay between retrying failed connections")
	flagset.DurationVar(&r.udpTimeout, "udp-timeout", UDP_TIMEOUT, "Time limit to wait for UDP replies")
	flagset.StringVar(&r.logLevel, "log-level", "info", "Sets the log level (debug, info, warn or error)")
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
	return "uhppoted-tunnel [--debug] [--console] [--lockfile <PID filepath>] --portal <UDP connection> --pipe <TCP connection>"
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

func (cmd *Run) execute(f func(t *tunnel.Tunnel)) (err error) {
	var portal tunnel.UDP
	var pipe tunnel.TCP

	// ... create UDP packet handler
	switch {
	case cmd.portal == "":
		err = fmt.Errorf("--portal argument is required")
		return

	case strings.HasPrefix(cmd.portal, "udp/listen:"):
		if portal, err = udp.NewUDPListen(cmd.portal[11:]); err != nil {
			return
		}

	case strings.HasPrefix(cmd.portal, "udp/broadcast:"):
		if portal, err = udp.NewUDPBroadcast(cmd.portal[14:], cmd.udpTimeout); err != nil {
			return
		}

	default:
		err = fmt.Errorf("Invalid --portal argument (%v)", cmd.portal)
		return
	}

	// ... create TCP/IP pipe
	switch {
	case cmd.pipe == "":
		err = fmt.Errorf("--pipe argument is required")
		return

	case strings.HasPrefix(cmd.pipe, "tcp/client:"):
		if pipe, err = tcp.NewTCPClient(cmd.pipe[11:], cmd.maxRetries, cmd.maxRetryDelay); err != nil {
			return
		}

	case strings.HasPrefix(cmd.pipe, "tcp/server:"):
		if pipe, err = tcp.NewTCPServer(cmd.pipe[11:]); err != nil {
			return
		}

	case strings.HasPrefix(cmd.pipe, "tls/client:"):
		if pipe, err = tls.NewTLSClient(cmd.pipe[11:], CACERT, cmd.maxRetries, cmd.maxRetryDelay); err != nil {
			return
		}

	case strings.HasPrefix(cmd.pipe, "tls/server:"):
		if pipe, err = tls.NewTLSServer(cmd.pipe[11:]); err != nil {
			return
		}

	default:
		err = fmt.Errorf("Invalid --pipe argument (%v)", cmd.pipe)
		return
	}

	// ... create lockfile

	if err := os.MkdirAll(cmd.workdir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", cmd.workdir, err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())
	lockfile := cmd.lockfile

	if lockfile == "" {
		hash := sha1.Sum([]byte(cmd.portal + cmd.pipe))
		lockfile = filepath.Join(cmd.workdir, fmt.Sprintf("%s-%x.pid", SERVICE, hash))
	}

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

	t := tunnel.NewTunnel(portal, pipe)

	f(t)

	return nil
}

func (cmd *Run) run(t *tunnel.Tunnel, interrupt chan os.Signal) {
	log.SetDebug(cmd.debug)
	log.SetLevel(cmd.logLevel)

	t.Run(interrupt)
}
