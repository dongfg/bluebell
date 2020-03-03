package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Basic global config
var Basic Config

// Config application config
type Config struct {
	Port    int
	Service struct {
		Name          string
		Address       string
		Port          int
		CheckURL      string `yaml:"check-url"`
		CheckInterval string `yaml:"check-interval"`
	}
	Series struct {
		Domain string
	}
}

// Load config from consul server
func Load() error {
	configPath := "config.yml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	Basic = config
	return nil
}
