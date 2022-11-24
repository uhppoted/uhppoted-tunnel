package commands

import (
	"path/filepath"
)

var DefaultConfig = filepath.Join(workdir(), "uhppoted-tunnel.toml")
var DefaultLockfile = filepath.Join(workdir(), "uhppoted-tunnel.pid")
