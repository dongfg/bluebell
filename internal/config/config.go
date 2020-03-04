package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Config application config
type Config struct {
	Port   int
	Series struct {
		Domain string
	}
}

// New config from config.yml file
func New() *Config {
	configPath := "config.yml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
