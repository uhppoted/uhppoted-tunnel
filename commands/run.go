package commands

import (
	"context"
	"crypto/sha1"
	TLS "crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	lib "github.com/uhppoted/uhppoted-lib/lockfile"

	"github.com/uhppoted/uhppoted-tunnel/log"
	"github.com/uhppoted/uhppoted-tunnel/tunnel"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/http"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/tcp"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/tls"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/udp"
)

type Run struct {
	conf              string
	label             string
	in                string
	out               string
	maxRetries        int
	maxRetryDelay     time.Duration
	udpTimeout        time.Duration
	caCertificate     string
	certificate       string
	key               string
	requireClientAuth bool
	html              string
	lockfile          string
	logFile           string
	logFileSize       int
	logLevel          string
	workdir           string
	debug             bool
	console           bool
}

const MAX_RETRIES = 32
const MAX_RETRY_DELAY = 5 * time.Minute
const UDP_TIMEOUT = 5 * time.Second

func (r *Run) flags() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.conf, "config", r.conf, "optional tunnel TOML configuration file")
	flagset.StringVar(&r.in, "in", r.in, "tunnel connection that accepts external requests e.g. udp/listen:0.0.0.0:60000 or tcp/client:101.102.103.104:54321")
	flagset.StringVar(&r.out, "out", r.out, "tunnel connection that dispatches received requests e.g. udp/broadcast:255.255.255.255:60000 or tcp/server:0.0.0.0:54321")
	flagset.StringVar(&r.lockfile, "lockfile", r.lockfile, "(optional) name of lockfile used to prevent running multiple copies of the service. A default lockfile name is generated if none is supplied")
	flagset.IntVar(&r.maxRetries, "max-retries", r.maxRetries, "Maximum number of times to retry failed connection. Defaults to -1 (retry forever)")
	flagset.DurationVar(&r.maxRetryDelay, "max-retry-delay", r.maxRetryDelay, "Maximum delay between retrying failed connections")
	flagset.DurationVar(&r.udpTimeout, "udp-timeout", r.udpTimeout, "Time limit to wait for UDP replies")

	flagset.StringVar(&r.caCertificate, "ca-cert", r.caCertificate, "File path for CA certificate PEM file (defaults to ca.cert)")
	flagset.StringVar(&r.certificate, "cert", r.certificate, "File path for client/server TLS certificate PEM file (defaults to client.cert or server.cert)")
	flagset.StringVar(&r.key, "key", r.key, "File path for client/server TLS key PEM file (defaults to client.key or server.key)")
	flagset.BoolVar(&r.requireClientAuth, "client-auth", r.requireClientAuth, "Requires client authentication for TLS")

	flagset.StringVar(&r.html, "html", r.html, "HTML folder for HTTP/HTTPS connectors")
	flagset.StringVar(&r.logLevel, "log-level", r.logLevel, "Sets the log level (debug, info, warn or error)")
	flagset.BoolVar(&r.console, "console", r.console, "Runs as a console application rather than a service")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Enables detailed debugging logs")

	return flagset
}

func (cmd *Run) Name() string {
	return "run"
}

func (cmd *Run) Description() string {
	return "Runs the uhppoted-tunnel daemon/service until terminated by the system service manager"
}

func (cmd *Run) Usage() string {
	return "uhppoted-tunnel [--debug] [--console] [--lockfile <PID filepath>] --in <connection> --out <connection>"
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

func (cmd *Run) execute(f func(t *tunnel.Tunnel, ctx context.Context, cancel context.CancelFunc)) (err error) {
	var in tunnel.Conn
	var out tunnel.Conn
	var ctx, cancel = context.WithCancel(context.Background())

	defer cancel()

	// ... create request handler
	switch {
	case cmd.in == "":
		err = fmt.Errorf("--in argument is required")
		return

	case
		strings.HasPrefix(cmd.in, "udp/listen:"),
		strings.HasPrefix(cmd.in, "tcp/client:"),
		strings.HasPrefix(cmd.in, "tcp/server:"),
		strings.HasPrefix(cmd.in, "tls/client:"),
		strings.HasPrefix(cmd.in, "tls/server:"),
		strings.HasPrefix(cmd.in, "http/"),
		strings.HasPrefix(cmd.in, "https/"):
		if in, err = cmd.makeConn("--in", cmd.in, ctx); err != nil {
			return
		}

	default:
		err = fmt.Errorf("Invalid --in argument (%v)", cmd.in)
		return
	}

	// ... create request dispatcher
	switch {
	case cmd.out == "":
		err = fmt.Errorf("--out argument is required")
		return

	case
		strings.HasPrefix(cmd.out, "udp/broadcast:"),
		strings.HasPrefix(cmd.out, "tcp/client:"),
		strings.HasPrefix(cmd.out, "tcp/server:"),
		strings.HasPrefix(cmd.out, "tls/client:"),
		strings.HasPrefix(cmd.out, "tls/server:"):
		if out, err = cmd.makeConn("--out", cmd.out, ctx); err != nil {
			return
		}

	default:
		err = fmt.Errorf("Invalid --out argument (%v)", cmd.out)
		return
	}

	// ... create lockfile
	dir := os.TempDir()
	lockfile := cmd.lockfile

	if lockfile == "" {
		hash := sha1.Sum([]byte(cmd.in + cmd.out))
		lockfile = filepath.Join(dir, fmt.Sprintf("%s-%x.pid", SERVICE, hash))
	}

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("%-5s %v\n", "FATAL", err)
		}
	}()

	if lock, err := lib.MakeLockFile(lockfile); err != nil {
		return err
	} else {
		defer func() {
			lock.Release()
		}()

		log.SetFatalHook(func() {
			lock.Release()
		})
	}

	// ... run
	if err := os.MkdirAll(cmd.workdir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", cmd.workdir, err)
	}

	t := tunnel.NewTunnel(in, out, ctx)

	f(t, ctx, cancel)

	return nil
}

func (cmd Run) makeConn(arg, spec string, ctx context.Context) (tunnel.Conn, error) {
	retry := conn.NewBackoff(cmd.maxRetries, cmd.maxRetryDelay, ctx)
	switch {
	case strings.HasPrefix(spec, "udp/listen:"):
		return udp.NewUDPListen(spec[11:], retry, ctx)

	case strings.HasPrefix(spec, "udp/broadcast:"):
		return udp.NewUDPBroadcast(spec[14:], cmd.udpTimeout, ctx)

	case strings.HasPrefix(spec, "tcp/client:"):
		return tcp.NewTCPClient(spec[11:], retry, ctx)

	case strings.HasPrefix(spec, "tcp/server:"):
		return tcp.NewTCPServer(spec[11:], retry, ctx)

	case strings.HasPrefix(spec, "tls/client:"):
		if ca, err := tlsCA(cmd.caCertificate); err != nil {
			return nil, err
		} else if certificate, err := tlsClientKeyPair(cmd.certificate, cmd.key); err != nil {
			return nil, err
		} else {
			return tls.NewTLSClient(spec[11:], ca, certificate, retry, ctx)
		}

	case strings.HasPrefix(spec, "tls/server:"):
		if ca, err := tlsCA(cmd.caCertificate); err != nil {
			return nil, err
		} else if certificate, err := tlsServerKeyPair(cmd.certificate, cmd.key); err != nil {
			return nil, err
		} else {
			return tls.NewTLSServer(spec[11:], ca, *certificate, cmd.requireClientAuth, retry, ctx)
		}

	case strings.HasPrefix(spec, "http/"):
		return http.NewHTTP(spec[5:], cmd.html, retry, ctx)

	case strings.HasPrefix(spec, "https/"):
		if ca, err := tlsCA(cmd.caCertificate); err != nil {
			return nil, err
		} else if certificate, err := tlsServerKeyPair(cmd.certificate, cmd.key); err != nil {
			return nil, err
		} else {
			fmt.Printf("%v\n%v\n%v\n%v\n", cmd.caCertificate, cmd.certificate, cmd.key, cmd.requireClientAuth)
			return http.NewHTTPS(spec[6:], cmd.html, ca, *certificate, cmd.requireClientAuth, retry, ctx)
		}

	default:
		return nil, fmt.Errorf("Invalid %v argument (%v)", arg, spec)
	}
}

func (cmd *Run) run(t *tunnel.Tunnel, ctx context.Context, cancel context.CancelFunc, interrupt chan os.Signal) {
	log.SetDebug(cmd.debug)
	log.SetLevel(cmd.logLevel)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := t.Run(interrupt); err != nil {
			log.Errorf("%-5s %v\n", "FATAL", err)
		}
	}()

	<-interrupt

	cancel()
	wg.Wait()
}

func tlsCA(cacert string) (*x509.CertPool, error) {
	if cacert == "" {
		cacert = "ca.cert"
	}

	ca := x509.NewCertPool()
	if bytes, err := os.ReadFile(cacert); err != nil {
		return nil, err
	} else if !ca.AppendCertsFromPEM(bytes) {
		return nil, fmt.Errorf("unable to parse CA certificate")
	}

	return ca, nil
}

func tlsServerKeyPair(certfile, keyfile string) (*TLS.Certificate, error) {
	if certfile == "" {
		certfile = "server.cert"
	}

	if keyfile == "" {
		keyfile = "server.key"
	}

	certificate, err := TLS.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}

	return &certificate, nil
}

func tlsClientKeyPair(certfile, keyfile string) (*TLS.Certificate, error) {
	if certfile != "" && keyfile != "" {
		certificate, err := TLS.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			return nil, err
		}

		return &certificate, nil
	}

	certificate, err := TLS.LoadX509KeyPair("client.cert", "client.key")
	if err != nil {
		return nil, nil
	}

	return &certificate, nil
}
