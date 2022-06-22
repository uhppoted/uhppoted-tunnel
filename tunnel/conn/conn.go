package conn

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

type Conn struct {
	Tag string
}

func (c Conn) Dumpf(message []byte, format string, args ...any) {
	Dumpf(c.Tag, message, format, args...)
}

func (c Conn) Debugf(format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", c.Tag, format)

	log.Debugf(f, args...)
}

func (c Conn) Infof(format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", c.Tag, format)

	log.Infof(f, args...)
}

func (c Conn) Warnf(format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", c.Tag, format)

	log.Warnf(f, args...)
}

func (c Conn) Errorf(format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", c.Tag, format)

	log.Errorf(f, args...)
}

func (c Conn) Fatalf(format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", c.Tag, format)

	log.Fatalf(f, args...)
}

func Dump(m []byte, prefix string) string {
	p := regexp.MustCompile(`\s*\|.*?\|`).ReplaceAllString(hex.Dump(m), "")
	q := regexp.MustCompile("(?m)^(.*)").ReplaceAllString(p, prefix+"$1")

	return fmt.Sprintf("%s", q)
}

func Dumpf(tag string, message []byte, format string, args ...any) {
	hex := Dump(message, "                                  ")
	preamble := fmt.Sprintf(format, args...)

	debugf(tag, "%v\n%s", preamble, hex)
}
