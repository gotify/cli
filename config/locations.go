package config

import (
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/gotify/cli/utils"
)

func userDir() string {
	usr, err := user.Current()
	if err != nil {
		utils.Exit1With(err)
	}
	return usr.HomeDir
}

func GetLocations() []string {
	inUserDir := filepath.Join(userDir(), ".gotify", "cli.json")
	etcPath := "/etc/gotify/cli.json"
	relativePath := "./cli.json"
	if runtime.GOOS == "windows" {
		return []string{relativePath, inUserDir}
	} else {
		return []string{relativePath, inUserDir, etcPath}
	}
}
