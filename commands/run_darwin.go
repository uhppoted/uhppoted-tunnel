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
	"github.com/uhppoted/uhppoted-lib/eventlog"
	"github.com/uhppoted/uhppoted-tunnel/tunnel"
)

var RUN = Run{
	console:     false,
	conf:        "/usr/local/etc/com.github.uhppoted/uhppoted-tunnel.conf",
	workdir:     "/usr/local/var/com.github.uhppoted",
	logFile:     fmt.Sprintf("/usr/local/var/com.github.uhppoted/logs/%s.log", SERVICE),
	logFileSize: 10,
}

func (cmd *Run) FlagSet() *flag.FlagSet {
	return cmd.flags()
}

func (cmd *Run) Execute(args ...interface{}) error {
	log.Printf("%s service %s - %s (PID %d)\n", SERVICE, uhppote.VERSION, "MacOS", os.Getpid())

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

	if !cmd.console {
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
