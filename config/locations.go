package config

import (
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
)

func GetLocations() (res []string) {
	res = append(res, "./cli.json")

	xdgPath, err := xdg.ConfigFile(filepath.Join("gotify", "cli.json"))
	if err == nil {
		res = append(res, xdgPath)
	}

	if usr, err := user.Current(); err == nil {
		res = append(res, filepath.Join(usr.HomeDir, ".gotify", "cli.json"))
	}

	if runtime.GOOS != "windows" {
		if usr, err := user.Current(); err == nil {
			fallbackXdgPath := filepath.Join(usr.HomeDir, ".config", "gotify", "cli.json")
			if xdgPath != fallbackXdgPath {
				res = append(res, fallbackXdgPath)
			}
		}

		res = append(res, "/etc/gotify/cli.json")
	}
	return
}
