package tailscale

import (
	"fmt"
	"regexp"
	"strconv"
)

func resolveTailscaleAddr(spec string) (string, uint16, error) {
	if match := regexp.MustCompile("(.*?):([0-9]+)").FindStringSubmatch(spec); len(match) < 3 {
		return "", 0, fmt.Errorf("invalid tailscale address (%v)", spec)
	} else if port, err := strconv.ParseUint(match[2], 10, 16); err != nil {
		return "", 0, err
	} else if port > 65535 {
		return "", 0, fmt.Errorf("invalid tailscale port (%v)", spec)
	} else {
		return match[1], uint16(port), nil
	}
}
