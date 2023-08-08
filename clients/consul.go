package clients

import (
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type Consul struct {
	client *api.Client
}

func NewConsulAPI() (*Consul, error) {
	cfg := api.DefaultConfig()
	consul, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Consul{client: consul}, nil
}

func (c *Consul) Agent() *api.Agent {
	return c.client.Agent()
}

type RegisterCfg struct {
	ID          string
	Host        string
	Port        string
	Environment string
}

func (c *Consul) Register(cfg *RegisterCfg) error {
	var address string
	if cfg.Environment == "dev" {
		address = "127.0.0.1"
	} else {
		address = cfg.Host
	}

	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s/", net.JoinHostPort(cfg.Host, cfg.Port)),
		Interval:                       "10s",
		Timeout:                        "30s",
		CheckID:                        fmt.Sprintf("service:%s:http", cfg.ID),
		DeregisterCriticalServiceAfter: "1m",
		TLSSkipVerify:                  func() bool { return cfg.Environment == "dev" }(),
	}

	port, _ := strconv.Atoi(cfg.Port)

	serviceDefinition := &api.AgentServiceRegistration{
		ID:      cfg.ID,
		Name:    cfg.ID + "_master",
		Port:    port,
		Address: address,
		Tags:    []string{cfg.Environment, cfg.ID},
		Check:   check,
	}

	if err := c.client.Agent().ServiceRegister(serviceDefinition); err != nil {
		return err
	}

	return nil
}
