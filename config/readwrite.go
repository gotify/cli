package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrNoneSet = errors.New("no config set, run 'gotify init'")
)

func ExistingConfig(locations []string) (string, error) {
	for _, location := range locations {
		_, err := os.Stat(location)

		if err == nil {
			return location, nil
		}
	}

	for _, location := range locations {
		_, err := os.Stat(location)

		if os.IsPermission(err) {
			return "", errors.New("Permission denied " + location)
		}
	}

	return "", ErrNoneSet
}

func ReadConfig(locations []string) (*Config, error) {
	location, err := ExistingConfig(locations)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(location)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	conf := Config{}
	conf.FromLocation = location
	err = json.NewDecoder(file).Decode(&conf)

	if err != nil {
		return nil, errors.New("Error while decoding " + err.Error())
	}
	return &conf, nil
}

func WriteConfig(location string, conf *Config) error {
	bytes, err := json.MarshalIndent(conf,"", "  ")
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(location), 0755)
	return ioutil.WriteFile(location, bytes, 0644)
}
