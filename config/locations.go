package config

import (
	"os/user"
	"path/filepath"
	"runtime"
)

func GetLocations() (res []string) {
	res = append(res, "./cli.json")

	if usr, err := user.Current(); err == nil {
		res = append(res, filepath.Join(usr.HomeDir, ".gotify", "cli.json"))
	}

	if runtime.GOOS != "windows" {
		res = append(res, "/etc/gotify/cli.json")
	}
	return
}
