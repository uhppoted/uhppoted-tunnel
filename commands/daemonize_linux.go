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
	Portal        string
	Pipe          string
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
ExecStart={{.Executable}} --portal {{.Portal}} --pipe {{.Pipe}}
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
	portal    string
	pipe      string
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("daemonize", flag.ExitOnError)

	flagset.StringVar(&cmd.label, "label", "", "Identifying label for the service (to distinguish multiple tunnels running on the same machine)")
	flagset.StringVar(&cmd.portal, "portal", "", "UDP connection e.g. udp/listen:0.0.0.0:60000 or udp/broadcast:255.255.255.255:60000")
	flagset.StringVar(&cmd.pipe, "pipe", "", "TCP pipe connection e.g. tcp/server:0.0.0.0:54321 or tcp/client:101.102.103.104:54321")
	flagset.Var(&cmd.usergroup, "user", "user:group for uhppoted-tunnel service")

	return flagset
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Daemonizes %s as a service/daemon", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return "daemonize [--user <user:group>] --portal <UDP connection> --pipe <TCP connection> [--label <label>]"
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize [--user <user:group>] --portal <UDP connection> --pipe <TCP connection> [--label <label>]\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Registers %s as a systemd service/daemon that runs on startup.\n", SERVICE)
	fmt.Println("      Defaults to the user:group uhppoted:uhppoted unless otherwise specified")
	fmt.Println("      with the --user option")
	fmt.Println()
	fmt.Println("   --portal <connection>  Specifies the UDP connection to use for the daemon/service")
	fmt.Println("   --pipe   <connection>  Specifies the TCP connection to use for the daemon/service")
	fmt.Println("   --label  <label>       Suffix to append to the service to uniquely identify different tunnels")
	fmt.Println("                          running on the same machine")
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
	}
	// ... check UDP packet handler
	switch {
	case cmd.portal == "":
		return fmt.Errorf("--portal argument is required")

	case strings.HasPrefix(cmd.portal, "udp/listen:"):

	case strings.HasPrefix(cmd.portal, "udp/broadcast:"):

	default:
		return fmt.Errorf("Invalid --portal argument (%v)", cmd.portal)
	}

	// ... check TCP/IP pipe
	switch {
	case cmd.pipe == "":
		return fmt.Errorf("--pipe argument is required")

	case strings.HasPrefix(cmd.pipe, "tcp/client:"):

	case strings.HasPrefix(cmd.pipe, "tcp/server:"):

	default:
		return fmt.Errorf("Invalid --pipe argument (%v)", cmd.pipe)
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

	service := fmt.Sprintf("%v", SERVICE)
	if cmd.label != "" {
		service = fmt.Sprintf("%v-%v", SERVICE, cmd.label)
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

	pid := fmt.Sprintf("/var/uhppoted/%v.pid", SERVICE)
	if cmd.label != "" {
		pid = fmt.Sprintf("/var/uhppoted/%v-%v.pid", SERVICE, cmd.label)
	}

	fmt.Println()
	fmt.Println("   ... daemonizing")

	i := info{
		Description:   "UHPPOTE UTO311-L0x access card controllers UDP tunnel service/daemon ",
		Documentation: "https://github.com/uhppoted/uhppoted-tunnel",
		Executable:    executable,
		Portal:        cmd.portal,
		Pipe:          cmd.pipe,
		PID:           pid,
		User:          "uhppoted",
		Group:         "uhppoted",
		Uid:           uid,
		Gid:           gid,
		LogFiles: []string{
			fmt.Sprintf("/var/log/uhppoted/%s.log", service),
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
	var tcp *net.TCPAddr

	switch {
	case strings.HasPrefix(cmd.portal, "udp/listen:"):
		udp, _ = net.ResolveUDPAddr("udp", cmd.portal[11:])

	case strings.HasPrefix(cmd.portal, "udp/broadcast:"):
		udp, _ = net.ResolveUDPAddr("udp", cmd.portal[14:])
	}

	switch {
	case strings.HasPrefix(cmd.pipe, "tcp/client:"):
		tcp, _ = net.ResolveTCPAddr("tcp", cmd.pipe[11:])

	case strings.HasPrefix(cmd.pipe, "tcp/server:"):
		tcp, _ = net.ResolveTCPAddr("tcp", cmd.pipe[11:])
	}

	fmt.Printf("   ... %s registered as a systemd service\n", service)
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Printf("     > sudo systemctl start  %s\n", service)
	fmt.Printf("     > sudo systemctl status %s\n", service)
	fmt.Println()

	if udp != nil {
		fmt.Println("   The firewall may need additional rules to allow UDP broadcast e.g. for UFW:")
		fmt.Println()
		fmt.Printf("     > sudo ufw allow from %s to any port 60000 proto udp\n", udp)
		fmt.Println()
	}

	if tcp != nil {
		fmt.Println("   The firewall may need additional rules to allow external access to the tunnel TCP pipe e.g. for UFW:")
		fmt.Println()
		fmt.Printf("     > sudo ufw allow from %s to any port 8080 proto tcp\n", tcp)
		fmt.Println()
	}

	fmt.Printf("   The installation can be verified by running the %v service in 'console' mode:\n", SERVICE)
	fmt.Println()
	fmt.Printf("     > sudo su %v\n", username)
	fmt.Printf("     > ./%v --debug --console --portal %v --pipe %v\n", SERVICE, cmd.portal, cmd.pipe)
	fmt.Println()
	fmt.Println()

	return nil
}

func (cmd *Daemonize) systemd(i *info) error {
	service := fmt.Sprintf("%v.service", SERVICE)
	if cmd.label != "" {
		service = fmt.Sprintf("%v-%v.service", SERVICE, cmd.label)
	}

	path := filepath.Join("/etc/systemd/system", service)
	t := template.Must(template.New(service).Parse(serviceTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
	service := fmt.Sprintf("%v.service", SERVICE)
	if cmd.label != "" {
		service = fmt.Sprintf("%v-%v.service", SERVICE, cmd.label)
	}

	path := filepath.Join("/etc/logrotate.d", service)
	t := template.Must(template.New(fmt.Sprintf("%s.logrotate", service)).Parse(logRotateTemplate))

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
