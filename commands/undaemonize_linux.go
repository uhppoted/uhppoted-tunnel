package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

var UNDAEMONIZE = Undaemonize{
	workdir: "/var/uhppoted/tunnel",
	logdir:  "/var/log/uhppoted",
	config:  "/etc/uhppoted/uhppoted.conf",
	etc:     "/usr/etc/uhppoted/tunnel",
}

type Undaemonize struct {
	workdir string
	logdir  string
	config  string
	etc     string
	label   string
}

func (cmd *Undaemonize) Name() string {
	return "undaemonize"
}

func (cmd *Undaemonize) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("undaemonize", flag.ExitOnError)

	flagset.StringVar(&cmd.label, "label", "", "Identifying label for the service (to distinguish multiple tunnels running on the same machine)")

	return flagset
}

func (cmd *Undaemonize) Description() string {
	return fmt.Sprintf("Deregisters %s as a service/daemon", SERVICE)
}

func (cmd *Undaemonize) Usage() string {
	return ""
}

func (cmd *Undaemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s undaemonize [--label <label>]\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Deregisters %s from launchd as a service/daemon", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Undaemonize) Execute(args ...interface{}) error {
	fmt.Println("   ... undaemonizing")

	service := SERVICE
	if cmd.label != "" {
		service = fmt.Sprintf("%v-%v", SERVICE, cmd.label)
	}

	if err := cmd.systemd(); err != nil {
		return err
	}

	if err := cmd.logrotate(); err != nil {
		return err
	}

	if err := cmd.clean(); err != nil {
		return err
	}

	fmt.Printf("   ... %s unregistered as a systemd service\n", service)
	fmt.Printf(`
       NOTE: Configuration files in %s,
             working files in %s,
             log files in %s
             were not removed and should be deleted manually
`, filepath.Dir(cmd.config), cmd.workdir, cmd.logdir)
	fmt.Println()

	return nil
}

func (cmd *Undaemonize) systemd() error {
	service := SERVICE
	if cmd.label != "" {
		service = fmt.Sprintf("%v-%v", SERVICE, cmd.label)
	}

	path := filepath.Join("/etc/systemd/system", fmt.Sprintf("%v.service", service))
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... nothing to do for 'systemd'   (%s does not exist)\n", path)
		return nil
	}

	fmt.Printf("   ... stopping %s service\n", service)
	command := exec.Command("systemctl", "stop", service)
	out, err := command.CombinedOutput()
	if strings.TrimSpace(string(out)) != "" {
		fmt.Printf("   > %s\n", out)
	}
	if err != nil {
		return fmt.Errorf("ERROR: Failed to stop '%s' (%v)\n", service, err)
	}

	fmt.Printf("   ... removing '%s'\n", path)
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Undaemonize) logrotate() error {
	service := SERVICE
	if cmd.label != "" {
		service = fmt.Sprintf("%v-%v", SERVICE, cmd.label)
	}

	path := filepath.Join("/etc/logrotate.d", service)

	fmt.Printf("   ... removing '%s'\n", path)

	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Undaemonize) clean() error {
	pid := fmt.Sprintf("/var/uhppoted/%v.pid", SERVICE)
	if cmd.label != "" {
		pid = fmt.Sprintf("/var/uhppoted/%v-%v.pid", SERVICE, cmd.label)
	}

	files := []string{
		filepath.Join(cmd.workdir, fmt.Sprintf("%s.pid", pid)),
	}

	directories := []string{
		cmd.logdir,
		cmd.workdir,
	}

	for _, f := range files {
		fmt.Printf("   ... removing '%s'\n", f)
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	warnings := []string{}
	for _, dir := range directories {
		fmt.Printf("   ... removing '%s'\n", dir)
		if err := os.Remove(dir); err != nil && !os.IsNotExist(err) {
			patherr, ok := err.(*os.PathError)
			if !ok {
				return err
			}

			syserr, ok := patherr.Err.(syscall.Errno)
			if !ok {
				return err
			}

			if syserr != syscall.ENOTEMPTY {
				return err
			}

			warnings = append(warnings, fmt.Sprintf("could not remove directory '%s' (%v)", dir, syserr))
		}
	}

	if len(warnings) > 0 {
		fmt.Println()
		for _, w := range warnings {
			fmt.Printf("   ... WARNING: %v\n", w)
		}
	}

	return nil
}
