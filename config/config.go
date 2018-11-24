package config

type Config struct {
	Token        string `json:"token"`
	URL          string `json:"url"`
	FromLocation string `json:"-"`
}
