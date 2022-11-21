package commands

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

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
		return label, fmt.Errorf("A valid IN connector is required")

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
		return label, fmt.Errorf("Invalid IN connector (%v)", in)
	}

	// ... verify OUT connector
	switch {
	case out == "":
		return label, fmt.Errorf("A valid OUT connector is required")

	case
		strings.HasPrefix(out, "udp/broadcast:"),
		strings.HasPrefix(out, "tcp/client:"),
		strings.HasPrefix(out, "tcp/server:"),
		strings.HasPrefix(out, "tls/client:"),
		strings.HasPrefix(out, "tls/server:"):
	// OK

	default:
		return label, fmt.Errorf("Invalid OUT connector (%v)", out)
	}

	return label, nil
}
