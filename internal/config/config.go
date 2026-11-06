package config

import (
	"errors"
	"os"
)

type Config struct {
	GladosURL  string
	ErdtreeURL string
}

func LoadConfig() (Config, error) {
	var conf Config
	conf.GladosURL = os.Getenv("GLADOS_URL")
	conf.ErdtreeURL = os.Getenv("ERDTREE_URL")

	if conf.GladosURL == "" {
		return Config{}, errors.New("GLADOS_URL is empty")
	}

	if conf.ErdtreeURL == "" {
		return Config{}, errors.New("ERDTREE_URL is empty")
	}

	return conf, nil
}
