package config

type Config struct {
	Token           string `json:"token"`
	URL             string `json:"url"`
	DefaultPriority int    `json:"defaultPriority"`
	FromLocation    string `json:"-"`
}
