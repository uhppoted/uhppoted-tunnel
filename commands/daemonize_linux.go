package commands

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type usergroup string

type info struct {
	Description   string
	Documentation string
	Executable    string
	In            string
	Out           string
	PID           string
	User          string
	Group         string
	Uid           int
	Gid           int
	LogFiles      []string
}

const serviceTemplate = `[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=syslog.target network.target

[Service]
Type=simple
ExecStart={{.Executable}} --lockfile {{.PID}} --in {{.In}} --out {{.Out}}
PIDFile={{.PID}}
User={{.User}}
Group={{.Group}}

[Install]
WantedBy=multi-user.target
`

const logRotateTemplate = `{{range .LogFiles}}{{. }} {{end}}{
    daily
    rotate 30
    compress
        compresscmd /bin/bzip2
        compressext .bz2
        dateext
    missingok
    notifempty
    su uhppoted uhppoted
    postrotate
       /usr/bin/killall -HUP uhppoted-tunnel
    endscript
}
`

var DAEMONIZE = Daemonize{
	usergroup: "uhppoted:uhppoted",
	workdir:   "/var/uhppoted/tunnel",
	logdir:    "/var/log/uhppoted",
	config:    "/etc/uhppoted/uhppoted.conf",
	etc:       "/etc/uhppoted/tunnel",
	service:   SERVICE,
}

var replacer *strings.Replacer

type Daemonize struct {
	usergroup usergroup
	workdir   string
	logdir    string
	config    string
	html      string
	etc       string
	label     string
	in        string
	out       string
	service   string
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("daemonize", flag.ExitOnError)

	flagset.StringVar(&cmd.in, "in", "", "tunnel connection that accepts requests e.g. udp/listen:0.0.0.0:60000 or tcp/client:101.102.103.104:54321")
	flagset.StringVar(&cmd.out, "out", "", "tunnel connection that dispatches received requests e.g. udp/broadcast:255.255.255.255:60000 or tcp/server:0.0.0.0:54321")
	flagset.StringVar(&cmd.label, "label", "", "Identifying label for the service (to distinguish multiple tunnels running on the same machine)")
	flagset.Var(&cmd.usergroup, "user", "user:group for uhppoted-tunnel service")

	return flagset
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Daemonizes %s as a service/daemon", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return "daemonize [--user <user:group>] --in <connection> --out <connection> [--label <label>]"
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize [--user <user:group>] --in <connection> --out <connection> [--label <label>]\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Registers %s as a systemd service/daemon that runs on startup.\n", SERVICE)
	fmt.Println("      Defaults to the user:group uhppoted:uhppoted unless otherwise specified")
	fmt.Println("      with the --user option")
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Daemonize) Execute(args ...interface{}) error {
	r := bufio.NewReader(os.Stdin)

	if cmd.label == "" {
		fmt.Println()
		fmt.Printf("     **** WARNING: running daemonize without the --label option will overwrite any existing uhppoted-tunnel service.\n")
		fmt.Println()
		fmt.Printf("     Enter 'yes' to continue with the installation: ")

		text, err := r.ReadString('\n')
		if err != nil || strings.TrimSpace(text) != "yes" {
			fmt.Println()
			fmt.Printf("     -- installation cancelled --")
			fmt.Println()
			return nil
		}
	} else {
		cmd.service = fmt.Sprintf("%v-%v", SERVICE, cmd.label)
	}

	// ... check --in connection
	switch {
	case cmd.in == "":
		return fmt.Errorf("--in argument is required")

	case
		strings.HasPrefix(cmd.in, "udp/listen:"),
		strings.HasPrefix(cmd.in, "tcp/client:"),
		strings.HasPrefix(cmd.in, "tcp/server:"),
		strings.HasPrefix(cmd.in, "tls/client:"),
		strings.HasPrefix(cmd.in, "tls/server:"),
		strings.HasPrefix(cmd.in, "http/"),
		strings.HasPrefix(cmd.in, "https/"):
	// OK

	default:
		return fmt.Errorf("Invalid --in argument (%v)", cmd.in)
	}

	// ... check --out connection
	switch {
	case cmd.out == "":
		return fmt.Errorf("--out argument is required")

	case
		strings.HasPrefix(cmd.out, "udp/broadcast:"),
		strings.HasPrefix(cmd.out, "tcp/client:"),
		strings.HasPrefix(cmd.out, "tcp/server:"),
		strings.HasPrefix(cmd.out, "tls/client:"),
		strings.HasPrefix(cmd.out, "tls/server:"):

	default:
		return fmt.Errorf("Invalid --out argument (%v)", cmd.out)
	}

	dir := filepath.Dir(cmd.config)

	fmt.Println()
	fmt.Printf("     **** PLEASE MAKE SURE YOU HAVE A BACKUP COPY OF THE CONFIGURATION INFORMATION AND KEYS IN %s ***\n", dir)
	fmt.Println()
	fmt.Printf("     Enter 'yes' to continue with the installation: ")

	text, err := r.ReadString('\n')
	if err != nil || strings.TrimSpace(text) != "yes" {
		fmt.Println()
		fmt.Printf("     -- installation cancelled --")
		fmt.Println()
		return nil
	}

	return cmd.execute()
}

func (cmd *Daemonize) execute() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	uid, gid, err := getUserGroup(string(cmd.usergroup))
	if err != nil {
		fmt.Println()
		fmt.Printf("     **** PLEASE CREATE uid:gid %v (OR SPECIFY A DIFFERENT uid:gid WITH the --user OPTION) ***\n", cmd.usergroup)
		fmt.Println()
		return err
	}

	username := "uhppoted"
	if u, err := user.LookupId(fmt.Sprintf("%v", uid)); err == nil {
		username = u.Username
	}

	fmt.Println()
	fmt.Println("   ... daemonizing")

	i := info{
		Description:   "UHPPOTE UTO311-L0x access card controllers UDP tunnel service/daemon ",
		Documentation: "https://github.com/uhppoted/uhppoted-tunnel",
		Executable:    executable,
		In:            cmd.in,
		Out:           cmd.out,
		PID:           fmt.Sprintf("/var/uhppoted/%v.pid", cmd.service),
		User:          "uhppoted",
		Group:         "uhppoted",
		Uid:           uid,
		Gid:           gid,
		LogFiles: []string{
			fmt.Sprintf("/var/log/uhppoted/%s.log", cmd.service),
		},
	}

	chown := func(path string, info fs.DirEntry, err error) error {
		if err == nil {
			err = os.Chown(path, uid, gid)
		}
		return err
	}

	if err := cmd.systemd(&i); err != nil {
		return err
	}

	if err := cmd.mkdirs(&i); err != nil {
		return err
	}

	if err := cmd.logrotate(&i); err != nil {
		return err
	}

	if err := filepath.WalkDir(cmd.etc, chown); err != nil {
		return err
	}

	if err := filepath.WalkDir(cmd.workdir, chown); err != nil {
		return err
	}

	// .. get network addresses for UFW
	var udp *net.UDPAddr

	switch {
	case strings.HasPrefix(cmd.in, "udp/listen:"):
		udp, _ = net.ResolveUDPAddr("udp", cmd.in[11:])

	case strings.HasPrefix(cmd.out, "udp/broadcast:"):
		udp, _ = net.ResolveUDPAddr("udp", cmd.out[14:])
	}

	fmt.Printf("   ... %s registered as a systemd service\n", cmd.service)
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Printf(`     > sudo systemctl start  "%s"`, cmd.service)
	fmt.Println()
	fmt.Printf(`     > sudo systemctl status "%s"`, cmd.service)
	fmt.Println()
	fmt.Println()

	if udp != nil {
		fmt.Println("   The firewall may need additional rules to allow UDP broadcast e.g. for UFW:")
		fmt.Println()
		fmt.Printf("     > sudo ufw allow 60000/udp\n")
		fmt.Println()
	}

	fmt.Printf("   The installation can be verified by running the %v service in 'console' mode:\n", SERVICE)
	fmt.Println()
	fmt.Printf("     > sudo su %v\n", username)
	fmt.Printf("     > ./%v --debug --console --in %v --out %v\n", SERVICE, cmd.in, cmd.out)
	fmt.Println()
	fmt.Println()

	return nil
}

func (cmd *Daemonize) systemd(i *info) error {
	path := filepath.Join("/etc/systemd/system", fmt.Sprintf("%v.service", cmd.service))
	t := template.Must(template.New(cmd.service).Parse(serviceTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, i)
}

func (cmd *Daemonize) mkdirs(i *info) error {
	directories := []string{
		"/var/uhppoted",
		"/var/uhppoted/tunnel",
		"/var/log/uhppoted",
		"/etc/uhppoted",
		"/etc/uhppoted/tunnel",
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}

		if err := os.Chown(dir, i.Uid, i.Gid); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Daemonize) logrotate(i *info) error {
	path := filepath.Join("/etc/logrotate.d", cmd.service)
	t := template.Must(template.New(fmt.Sprintf("%s.logrotate", cmd.service)).Parse(logRotateTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, i)
}

// usergroup::flag.Value
func getUserGroup(s string) (int, int, error) {
	match := regexp.MustCompile(`(\w+?):(\w+)`).FindStringSubmatch(s)
	if match == nil {
		return 0, 0, fmt.Errorf("Invalid user:group '%s'", s)
	}

	u, err := user.Lookup(match[1])
	if err != nil {
		return 0, 0, err
	}

	g, err := user.LookupGroup(match[2])
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

func (f *usergroup) String() string {
	if f == nil {
		return "uhppoted:uhppoted"
	}

	return string(*f)
}

func (f *usergroup) Set(s string) error {
	_, _, err := getUserGroup(s)
	if err != nil {
		return err
	}

	*f = usergroup(strings.TrimSpace(s))

	return nil
}
