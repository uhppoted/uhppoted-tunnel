package commands

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/eventlog"

	"github.com/uhppoted/uhppoted-tunnel/tunnel"
)

var RUN = Run{
	in:                "",
	out:               "",
	maxRetries:        MAX_RETRIES,
	maxRetryDelay:     MAX_RETRY_DELAY,
	udpTimeout:        UDP_TIMEOUT,
	caCertificate:     "ca.cert",
	certificate:       "",
	key:               "",
	requireClientAuth: false,
	html:              "./html",
	lockfile: config.Lockfile{
		File:   DefaultLockfile,
		Remove: false,
	},
	logLevel: "info",
	debug:    false,
	console:  false,
	daemon:   false,

	conf:        "",
	workdir:     "/var/uhppoted",
	logFile:     fmt.Sprintf("/var/log/uhppoted/%s.log", SERVICE),
	logFileSize: 10,

	rateLimit:  1,
	burstLimit: 120,

	controllers: map[uint32]string{},
}

func (cmd *Run) FlagSet() *flag.FlagSet {
	return cmd.flags()
}

func (cmd *Run) Execute(args ...interface{}) error {
	log.Printf("%s service %s - %s (PID %d)\n", SERVICE, uhppote.VERSION, "Linux", os.Getpid())

	f := func(t *tunnel.Tunnel, ctx context.Context, cancel context.CancelFunc) {
		cmd.exec(t, ctx, cancel)
	}

	return cmd.execute(f)
}

func (cmd *Run) exec(t *tunnel.Tunnel, ctx context.Context, cancel context.CancelFunc) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	if !cmd.console || cmd.daemon {
		events := eventlog.Ticker{Filename: cmd.logFile, MaxSize: cmd.logFileSize}

		log.SetOutput(&events)
		log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

		rotate := make(chan os.Signal, 1)

		signal.Notify(rotate, syscall.SIGHUP)

		go func() {
			for {
				<-rotate
				log.Printf("Rotating %s log file '%s'\n", SERVICE, cmd.logFile)
				events.Rotate()
			}
		}()
	}

	cmd.run(t, ctx, cancel, interrupt)
}
