package udp

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"sync/atomic"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

var PACKETID uint32 = 0

func nextID() uint32 {
	return atomic.AddUint32(&PACKETID, 1)
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}

func dumpf(tag string, message []byte, format string, args ...any) {
	hex := dump(message, "                                  ")
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
