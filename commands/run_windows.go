package commands

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"golang.org/x/sys/windows/svc"
	syslog "golang.org/x/sys/windows/svc/eventlog"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/eventlog"

	"github.com/uhppoted/uhppoted-tunnel/tunnel"
)

var RUN = Run{
	console:     false,
	conf:        filepath.Join(workdir(), "uhppoted-tunnel.conf"),
	workdir:     workdir(),
	logFile:     filepath.Join(workdir(), "logs", fmt.Sprintf("%s.log", SERVICE)),
	logFileSize: 10,
}

type service struct {
	name   string
	conf   config.Config
	cmd    *Run
	tunnel *tunnel.Tunnel
	ctx    context.Context
	cancel context.CancelFunc
}

func (cmd *Run) FlagSet() *flag.FlagSet {
	flagset := cmd.flags()

	flagset.StringVar(&cmd.label, "label", "", "(optional) Identifying label for the service to distinguish multiple tunnels running on the same machine")

	return flagset
}

func (cmd *Run) Execute(args ...interface{}) error {
	name := SERVICE
	if cmd.label != "" {
		name = fmt.Sprintf("%v-%v", SERVICE, cmd.label)
	}

	log.Printf("%s service %s - %s (PID %d)\n", name, uhppote.VERSION, "Microsoft Windows", os.Getpid())

	f := func(t *tunnel.Tunnel, ctx context.Context, cancel context.CancelFunc) {
		cmd.start(t, ctx, cancel)
	}

	return cmd.execute(f)
}

func (cmd *Run) start(t *tunnel.Tunnel, ctx context.Context, cancel context.CancelFunc) {
	if cmd.console {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags)

		interrupt := make(chan os.Signal, 1)

		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		cmd.run(t, ctx, cancel, interrupt)
		return
	}

	name := SERVICE
	if cmd.label != "" {
		name = fmt.Sprintf("%v-%v", SERVICE, cmd.label)
	}

	if eventlogger, err := syslog.Open(name); err != nil {
		events := eventlog.Ticker{Filename: cmd.logFile, MaxSize: cmd.logFileSize}

		log.SetOutput(&events)
	} else {
		defer eventlogger.Close()

		log.SetOutput(&EventLog{eventlogger})
	}

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	log.Printf("%s service - start\n", name)

	uhppoted := service{
		name:   name,
		cmd:    cmd,
		tunnel: t,
		ctx:    ctx,
		cancel: cancel,
	}

	log.Printf("%s service - starting\n", name)

	if err := svc.Run(name, &uhppoted); err != nil {
		fmt.Printf("   Unable to execute ServiceManager.Run request (%v)\n", err)
		fmt.Println()
		fmt.Printf("   To run %s as a command line application, type:\n", SERVICE)
		fmt.Println()
		fmt.Printf("     > %s --console\n", SERVICE)
		fmt.Println()

		log.Printf("   Unable to execute ServiceManager.Run request (%v)\n", err)
		log.Println()
		log.Printf("   To run %s as a command line application, type:\n", SERVICE)
		log.Println()
		log.Printf("     > %s --console\n", SERVICE)
		log.Println()

		log.Panicf("Error executing ServiceManager.Run request: %v", err)

		return
	}

	log.Printf("%s daemon - started\n", name)
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {
	log.Printf("%s service - Execute\n", s.name)

	const commands = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	interrupt := make(chan os.Signal, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.cmd.run(s.tunnel, s.ctx, s.cancel, interrupt)

		log.Printf("exit\n")
	}()

	status <- svc.Status{State: svc.Running, Accepts: commands}

loop:
	for {
		select {
		case c := <-r:
			log.Printf("%s service - select: %v  %v\n", s.name, c.Cmd, c.CurrentStatus)
			switch c.Cmd {
			case svc.Interrogate:
				log.Printf("%s service - svc.Interrogate %v\n", s.name, c.CurrentStatus)
				status <- c.CurrentStatus

			case svc.Stop:
				interrupt <- syscall.SIGINT
				log.Printf("%s service- svc.Stop\n", s.name)
				break loop

			case svc.Shutdown:
				interrupt <- syscall.SIGTERM
				log.Printf("%s service - svc.Shutdown\n", s.name)
				break loop

			default:
				log.Printf("%s service - svc.????? (%v)\n", s.name, c.Cmd)
			}
		}
	}

	log.Printf("%s service - stopping\n", s.name)
	status <- svc.Status{State: svc.StopPending}
	wg.Wait()
	status <- svc.Status{State: svc.Stopped}
	log.Printf("%s service - stopped\n", s.name)

	return false, 0
}
