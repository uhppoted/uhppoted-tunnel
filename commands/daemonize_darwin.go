package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	xpath "github.com/uhppoted/uhppoted-lib/encoding/plist"
)

type info struct {
	Label      string
	Executable string
	StdLogFile string
	ErrLogFile string
}

type plist struct {
	Label             string
	Program           string
	WorkingDirectory  string
	ProgramArguments  []string
	KeepAlive         bool
	RunAtLoad         bool
	StandardOutPath   string
	StandardErrorPath string
}

const newsyslog = `#logfilename                                       [owner:group]  mode  count  size   when  flags [/pid_file]  [sig_num]
{{range .}}{{.LogFile}}  :              644   30     10000  @T00  J     {{.PID}}
{{end}}`

var DAEMONIZE = Daemonize{
	plist:   fmt.Sprintf("com.github.uhppoted.%s.plist", SERVICE),
	conf:    "/usr/local/etc/com.github.uhppoted/uhppoted-tunnel.toml",
	workdir: "/usr/local/var/com.github.uhppoted/tunnel",
	logdir:  "/usr/local/var/com.github.uhppoted/logs",
	etc:     "/usr/local/etc/com.github.uhppoted/tunnel",
}

var replacer *strings.Replacer

type Daemonize struct {
	plist   string
	workdir string
	logdir  string
	etc     string
	conf    string
	in      string
	out     string
	label   string
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("daemonize", flag.ExitOnError)

	flagset.StringVar(&cmd.conf, "config", cmd.conf, "tunnel TOML configuration file. Defaults to /usr/local/etc/com.github.uhppoted/uhppoted-tunnel.toml")
	flagset.StringVar(&cmd.in, "in", "", "tunnel connection that accepts requests e.g. udp/listen:0.0.0.0:60000 or tcp/client:101.102.103.104:54321")
	flagset.StringVar(&cmd.out, "out", "", "tunnel connection that dispatches received requests e.g. udp/broadcast:255.255.255.255:60000 or tcp/server:0.0.0.0:54321")
	flagset.StringVar(&cmd.label, "label", "", "(optional) Identifying label for the service to distinguish multiple tunnels running on the same machine")

	return flagset
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Daemonizes %s as a service/daemon", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return ""
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize [--config <TOML file>] [--in <connection>] [--out <connection>] [--label <label>]\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Daemonizes %s as a service/daemon that runs on startup\n", SERVICE)
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

func (cmd *Daemonize) configuration(flagset *flag.FlagSet) string {
	config := ""
	flagset.Visit(func(f *flag.Flag) {
		if f.Name == "config" {
			config = f.Value.String()
		}
	})

	file := config
	section := ""
	if match := regexp.MustCompile("(.*?)((?:::|#).*)").FindStringSubmatch(config); match != nil {
		file = match[1]
		section = match[2]
	}

	if file != "" {
		return config
	} else if f := flagset.Lookup("config"); f != nil && f.DefValue != "" {
		return f.DefValue + section
	}

	return ""
}

func (cmd *Daemonize) Execute(args ...any) error {
	// ... validate configuration
	in := ""
	out := ""
	label := ""

	if configuration, err := configure(cmd.conf); err != nil {
		return err
	} else {
		if v, ok := configuration["in"]; ok {
			if u, ok := v.(string); ok {
				in = u
			}
		}

		if v, ok := configuration["out"]; ok {
			if u, ok := v.(string); ok {
				out = u
			}
		}

		if v, ok := configuration["label"]; ok {
			if u, ok := v.(string); ok {
				label = u
			}
		}
	}

	if cmd.in != "" {
		in = cmd.in
	}

	if cmd.out != "" {
		out = cmd.out
	}

	if cmd.label != "" {
		label = cmd.label
	}

	// ... verify IN connector
	switch {
	case in == "":
		return fmt.Errorf("A valid IN connector is required")

	case
		strings.HasPrefix(in, "udp/listen:"),
		strings.HasPrefix(in, "tcp/client:"),
		strings.HasPrefix(in, "tcp/server:"),
		strings.HasPrefix(in, "tls/client:"),
		strings.HasPrefix(in, "tls/server:"),
		strings.HasPrefix(in, "http/"),
		strings.HasPrefix(in, "https/"):
	// OK

	default:
		return fmt.Errorf("Invalid IN connector (%v)", in)
	}

	// ... verify OUT connector
	switch {
	case out == "":
		return fmt.Errorf("A valid OUT connector is required")

	case
		strings.HasPrefix(out, "udp/broadcast:"),
		strings.HasPrefix(in, "tcp/client:"),
		strings.HasPrefix(in, "tcp/server:"),
		strings.HasPrefix(in, "tls/client:"),
		strings.HasPrefix(in, "tls/server:"):
	// OK

	default:
		return fmt.Errorf("Invalid OUT connector (%v)", out)
	}

	// ... verify label

	if label == "" {
		fmt.Println()
		fmt.Printf("     **** WARNING: running daemonize without the --label option will overwrite any existing uhppoted-tunnel daemon.\n")
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

	// ... install daemon
	if label != "" {
		cmd.plist = fmt.Sprintf("com.github.uhppoted.%v-%v.plist", SERVICE, label)
	}

	return cmd.execute(label)
}

func (cmd *Daemonize) execute(label string) error {
	fmt.Println()
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	l := fmt.Sprintf("com.github.uhppoted.%s", SERVICE)
	if label != "" {
		l = fmt.Sprintf("com.github.uhppoted.%s-%v", SERVICE, label)
	}

	i := info{
		Label:      l,
		Executable: executable,
		StdLogFile: filepath.Join(cmd.logdir, fmt.Sprintf("%s.log", SERVICE)),
		ErrLogFile: filepath.Join(cmd.logdir, fmt.Sprintf("%s.err", SERVICE)),
	}

	if err := cmd.launchd(&i); err != nil {
		return err
	}

	if err := cmd.mkdirs(); err != nil {
		return err
	}

	if err := cmd.logrotate(&i); err != nil {
		return err
	}

	if err := cmd.firewall(i); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a LaunchDaemon\n", i.Label)
	fmt.Println()
	fmt.Printf("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:\n")
	fmt.Println()
	fmt.Printf("   sudo launchctl load /Library/LaunchDaemons/%v.plist\n", l)
	fmt.Println()
	fmt.Println()

	return nil
}

func (cmd *Daemonize) launchd(i *info) error {
	path := filepath.Join("/Library/LaunchDaemons", cmd.plist)
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	args := []string{
		path, // ref. https://apple.stackexchange.com/questions/110644/getting-launchd-to-read-program-arguments-correctly
	}

	if cmd.conf != "" {
		if file, err := resolve(cmd.workdir, cmd.conf); err != nil {
			return err
		} else {
			args = append(args, "--config", file)
		}
	}

	if cmd.in != "" {
		args = append(args, "--in", cmd.in)
	}

	if cmd.out != "" {
		args = append(args, "--out", cmd.out)
	}

	pl := plist{
		Label:             i.Label,
		Program:           i.Executable,
		WorkingDirectory:  cmd.workdir,
		ProgramArguments:  args,
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   i.StdLogFile,
		StandardErrorPath: i.ErrLogFile,
	}

	if !os.IsNotExist(err) {
		current, err := cmd.parse(path)
		if err != nil {
			return err
		}

		pl.WorkingDirectory = current.WorkingDirectory
		pl.ProgramArguments = current.ProgramArguments
		pl.KeepAlive = current.KeepAlive
		pl.RunAtLoad = current.RunAtLoad
		pl.StandardOutPath = current.StandardOutPath
		pl.StandardErrorPath = current.StandardErrorPath
	}

	return cmd.daemonize(path, pl)
}

func (cmd *Daemonize) parse(path string) (*plist, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	p := plist{}
	decoder := xpath.NewDecoder(f)
	err = decoder.Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (cmd *Daemonize) daemonize(path string, p interface{}) error {
	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := xpath.NewEncoder(f)
	if err = encoder.Encode(p); err != nil {
		return err
	}

	return nil
}

func (cmd *Daemonize) mkdirs() error {
	directories := []string{
		cmd.workdir,
		cmd.logdir,
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Daemonize) logrotate(i *info) error {
	pid := filepath.Join(cmd.workdir, fmt.Sprintf("%s.pid", SERVICE))
	logfiles := []struct {
		LogFile string
		PID     string
	}{
		{
			LogFile: i.StdLogFile,
			PID:     pid,
		},
		{
			LogFile: i.ErrLogFile,
			PID:     pid,
		},
	}

	t := template.Must(template.New("logrotate.conf").Parse(newsyslog))
	path := filepath.Join("/etc/newsyslog.d", fmt.Sprintf("%s.conf", SERVICE))

	fmt.Printf("   ... creating '%s'\n", path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, logfiles)
}

func (cmd *Daemonize) firewall(i info) error {
	fmt.Println()
	fmt.Printf("   ***\n")
	fmt.Printf("   *** WARNING: adding '%s' to the application firewall and unblocking incoming connections\n", SERVICE)
	fmt.Printf("   ***\n")
	fmt.Println()

	path := i.Executable

	command := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := command.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to retrieve application firewall global state (%v)", err)
	}

	if strings.Contains(string(out), "State = 1") {
		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to disable the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to add 'uhppoted-tunnel' to the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to unblock 'uhppoted-tunnel' on the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to re-enable the application firewall (%v)", err)
		}

		fmt.Println()
	}

	return nil
}

func resolve(base string, cfg string) (string, error) {
	file := cfg
	section := ""

	if match := regexp.MustCompile("(.*?)((?:::|#).*)").FindStringSubmatch(cfg); match != nil {
		file = match[1]
		section = match[2]
	}

	if strings.HasPrefix(file, ".") {
		if abs, err := filepath.Abs(file); err != nil {
			return cfg, err
		} else {
			return fmt.Sprintf("%v%v", abs, section), nil
		}
	}

	return cfg, nil
}
