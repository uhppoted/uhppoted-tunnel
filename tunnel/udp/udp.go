package udp

import (
	"fmt"

	"github.com/uhppoted/uhppoted-tunnel/log"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

func dumpf(tag string, message []byte, format string, args ...any) {
	hex := conn.Dump(message, "                                  ")
	preamble := fmt.Sprintf(format, args...)

	debugf(tag, "%v\n%s", preamble, hex)
}

func debugf(tag, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Debugf(f, args...)
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
	log.Fatalf(format, args...)
}
