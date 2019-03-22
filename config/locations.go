package config

import (
	"os/user"
	"path/filepath"
	"runtime"
)

func userDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func GetLocations() (res []string) {
	res = append(res, "./cli.json")

	if usrDir, err := userDir(); err != nil {
		res = append(res, filepath.Join(usrDir, ".gotify", "cli.json"))
	}

	if runtime.GOOS != "windows" {
		res = append(res, "/etc/gotify/cli.json")
	}
	return
}
