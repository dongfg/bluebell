package consul

import (
	"github.com/hashicorp/consul/api"
)

// Consul client wrap
type Consul struct {
	client *api.Client
}

// Service consul agent service wrap
type Service struct {
	Name          string
	Address       string
	Port          int
	CheckUrl      string
	CheckInterval string
}

// New consul config with address(ip:port) and acl token
func New(address string, token string) *Consul {
	client, err := api.NewClient(&api.Config{
		Address: address,
		Token:   token,
	})
	if err != nil {
		panic(err)
	}
	return &Consul{
		client: client,
	}
}

// Fetch k/v config with given key
func (consul *Consul) Fetch(key string) string {
	pair, _, err := consul.client.KV().Get(key, nil)
	if err != nil {
		panic(err)
	}
	return string(pair.Value)
}

// Register service
func (consul *Consul) Register(service Service) {
	err := consul.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      service.Name,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Check: &api.AgentServiceCheck{
			HTTP:     service.CheckUrl,
			Interval: service.CheckInterval,
		},
	})
	if err != nil {
		panic(err)
	}
}

// Deregister service
func (consul *Consul) Deregister(serviceName string) {
	err := consul.client.Agent().ServiceDeregister(serviceName)
	if err != nil {
	}
}
