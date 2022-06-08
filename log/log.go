package log

import (
	"fmt"
	syslog "log"
)

var Debug = false

func Debugf(format string, args ...any) {
	if Debug {
		syslog.Printf("%-5v  %v", "DEBUG", fmt.Sprintf(format, args...))
	}
}

func Infof(format string, args ...any) {
	syslog.Printf("%-5v  %v", "INFO", fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
	syslog.Printf("%-5v  %v", "WARN", fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
	syslog.Printf("%-5v  %v", "ERROR", fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...any) {
	syslog.Fatalf("%-5v  %v", "FATAL", fmt.Sprintf(format, args...))
}
