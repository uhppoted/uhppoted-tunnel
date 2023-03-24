package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrLabel = errors.New("invalid label")

func (cmd *Daemonize) validate() (string, error) {
	in := ""
	out := ""
	label := ""

	if configuration, err := configure(cmd.conf); err != nil {
		return label, err
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
		return label, fmt.Errorf("a valid IN connector is required")

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
		return label, fmt.Errorf("invalid IN connector (%v)", in)
	}

	// ... verify OUT connector
	switch {
	case out == "":
		return label, fmt.Errorf("a valid OUT connector is required")

	case
		strings.HasPrefix(out, "udp/broadcast:"),
		strings.HasPrefix(out, "tcp/client:"),
		strings.HasPrefix(out, "tcp/server:"),
		strings.HasPrefix(out, "tls/client:"),
		strings.HasPrefix(out, "tls/server:"):
	// OK

	default:
		return label, fmt.Errorf("invalid OUT connector (%v)", out)
	}

	// ... check label
	if label == "" {
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
			return label, ErrLabel
		}
	}

	return label, nil
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
