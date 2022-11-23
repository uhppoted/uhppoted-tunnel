package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var DAEMONIZE = Daemonize{
	name:        SERVICE,
	description: "UHPPOTE UTO311-L0x access card controllers UDP tunnel service",
	conf:        filepath.Join(workdir(), "uhppoted-tunnel.toml"),
	workdir:     filepath.Join(workdir(), "tunnel"),
	logdir:      filepath.Join(workdir(), "logs"),
	etc:         filepath.Join(workdir(), "tunnel"),
}

type info struct {
	Executable string
	WorkDir    string
	HTML       string
	LogDir     string
	Label      string
}

type Daemonize struct {
	name        string
	description string
	conf        string
	workdir     string
	logdir      string
	etc         string
	in          string
	out         string
	label       string
}

var replacer = strings.NewReplacer(
	"\r\n", "\r\n",
	"\r", "\r\n",
	"\n", "\r\n",
)

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("daemonize", flag.ExitOnError)

	flagset.StringVar(&cmd.conf, "config", cmd.conf, "tunnel TOML configuration file. Defaults to "+DefaultConfig)
	flagset.StringVar(&cmd.in, "in", "", "tunnel connection that accepts requests e.g. udp/listen:0.0.0.0:60000 or tcp/client:101.102.103.104:54321")
	flagset.StringVar(&cmd.out, "out", "", "tunnel connection that dispatches received requests e.g. udp/broadcast:255.255.255.255:60000 or tcp/server:0.0.0.0:54321")
	flagset.StringVar(&cmd.label, "label", "", "(optional) Identifying label for the service to distinguish multiple tunnels running on the same machine")

	return flagset
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Registers %s as a Windows service", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return ""
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize --in <connection> --out <connection> [--label <label>]\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Registers %s as a Windows service that runs on startup.\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Daemonize) ParseCmd(args ...string) error {
	flagset := cmd.FlagSet()
	if flagset == nil {
		panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", cmd.Name(), cmd))
	}

	flagset.Parse(args)

	cmd.conf = cmd.configuration(flagset)

	return nil
}

func (cmd *Daemonize) Execute(args ...interface{}) error {
	// ... validate configuration
	label, err := cmd.validate()
	if err != nil {
		return err
	} else if label == "" {
		fmt.Println()
		fmt.Printf("     **** WARNING: running daemonize without the --label option will overwrite any existing uhppoted-tunnel service.\n")
		fmt.Println()
		fmt.Printf("     Enter 'yes' to continue with the installation: ")

		r := bufio.NewReader(os.Stdin)
		text, err := r.ReadString('\n')
		if err != nil || strings.TrimSpace(text) != "yes" {
			fmt.Println()
			fmt.Printf("     -- installation cancelled --")
			fmt.Println()
			return nil
		}
	}

	// ... install service
	return cmd.execute(label)
}

func (cmd *Daemonize) execute(label string) error {
	fmt.Println()
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	i := info{
		Executable: executable,
		WorkDir:    cmd.workdir,
		LogDir:     cmd.logdir,
		Label:      label,
	}

	if err := cmd.register(&i); err != nil {
		return err
	}

	if err := cmd.mkdirs(&i); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a Windows system service\n", cmd.name)
	fmt.Println()
	fmt.Println("   The service will start automatically on the next system restart. Start it manually from the")
	fmt.Println("   'Services' application or from the command line by executing the following command:")
	fmt.Println()
	fmt.Printf(`     > net start "%s"\n`, cmd.name)
	fmt.Printf(`     > sc query "%s"\n`, cmd.name)
	fmt.Println()

	return nil
}

func (cmd *Daemonize) register(i *info) error {
	// ... initialise service command line args
	args := []string{}

	if cmd.conf != "" {
		args = append(args, "--config", cmd.conf)
	}

	if cmd.in != "" {
		args = append(args, "--in", cmd.in)
	}

	if cmd.out != "" {
		args = append(args, "--out", cmd.out)
	}

	if i.Label != "" {
		args = append(args, "--label", i.Label)
	}

	// ... create service config
	config := mgr.Config{
		DisplayName:      cmd.name,
		Description:      cmd.description,
		StartType:        mgr.StartAutomatic,
		DelayedAutoStart: true,
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(cmd.name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", cmd.Name)
	}

	s, err = m.CreateService(cmd.name, i.Executable, config, args...)
	if err != nil {
		return err
	}

	defer s.Close()

	err = eventlog.InstallAsEventCreate(cmd.name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("InstallAsEventCreate() failed: %v", err)
	}

	return nil
}

func (cmd *Daemonize) mkdirs(i *info) error {
	directories := []string{
		i.WorkDir,
		i.LogDir,
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}
	}

	return nil
}
