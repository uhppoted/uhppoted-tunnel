package commands

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/pelletier/go-toml/v2"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

const (
	SERVICE = `uhppoted-tunnel`
)

func configure(configuration string) (map[string]any, error) {
	config := map[string]any{}

	if configuration == "" {
		return config, nil
	}

	file := configuration
	section := ""
	if match := regexp.MustCompile("(.*?)(?:::|#)(.*)").FindStringSubmatch(configuration); match != nil {
		file = match[1]
		section = match[2]
	}

	if bytes, err := os.ReadFile(file); err != nil {
		return nil, err
	} else {
		c := map[string]any{}
		if err := toml.Unmarshal(bytes, &c); err != nil {
			return nil, err
		}

		if m, ok := c["defaults"]; ok {
			if defaults, ok := m.(map[string]any); ok {
				for k, v := range defaults {
					config[k] = v
				}
			}
		}

		if m, ok := c[section]; ok {
			if tunnel, ok := m.(map[string]any); ok {
				for k, v := range tunnel {
					config[k] = v
				}
			}
		}
	}

	return config, nil
}

func helpOptions(flagset *flag.FlagSet) {
	flags := 0
	count := 0

	flag.VisitAll(func(f *flag.Flag) {
		count++
	})

	flagset.VisitAll(func(f *flag.Flag) {
		flags++
		fmt.Printf("    --%-13s %s\n", f.Name, f.Usage)
	})

	if count > 0 {
		fmt.Println()
		fmt.Println("  Options:")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("    --%-13s %s\n", f.Name, f.Usage)
		})
	}

	if flags > 0 {
		fmt.Println()
	}
}

func infof(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Infof(f, args...)
}

func warnf(tag, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Warnf(f, args...)
}

func errorf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Errorf(f, args...)
}

func fatalf(format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", "", format)

	log.Fatalf(f, args...)
}
