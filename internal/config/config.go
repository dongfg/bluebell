package config

import (
	"github.com/dongfg/bluebell/internal/consul"
	"gopkg.in/yaml.v2"
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
func Load(client *consul.Consul) error {
	rawConfig := client.Fetch("config/bluebell/yaml")
	config := Config{}
	err := yaml.Unmarshal([]byte(rawConfig), &config)
	if err != nil {
		return err
	}
	Basic = config
	return nil
}
