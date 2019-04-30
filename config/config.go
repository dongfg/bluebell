package config

import (
	"github.com/dongfg/bluebell/consul"
	"gopkg.in/yaml.v2"
)

var Basic Config

type Config struct {
	Port    int
	Service struct {
		Name          string
		Address       string
		Port          int
		CheckUrl      string `yaml:"check-url"`
		CheckInterval string `yaml:"check-interval"`
	}
}

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
