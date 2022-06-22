package conn

import (
	"encoding/hex"
	"fmt"
	"regexp"
)

func Dump(m []byte, prefix string) string {
	p := regexp.MustCompile(`\s*\|.*?\|`).ReplaceAllString(hex.Dump(m), "")
	q := regexp.MustCompile("(?m)^(.*)").ReplaceAllString(p, prefix+"$1")

	return fmt.Sprintf("%s", q)
}
