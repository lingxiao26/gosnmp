package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Hosts     []*Host    `yaml:"hosts"`
	Threshold *Threshold `yaml:"threshold"`
	Interval  *Interval  `yaml:"interval"`
}

type Interval struct {
	Disk   time.Duration `yaml:"disk"`
	Memory time.Duration `yaml:"memory"`
}

type Threshold struct {
	Disk   int `yaml:"disk"`
	Memory int `yaml:"memory"`
}

type Host struct {
	At      []string `yaml:"at"`
	Addr    string   `yaml:"addr"`
	Webhook string   `yaml:"webhook"`
}

func Load(cfgFile string) (*Configuration, error) {
	cfgBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(%s): %v", cfgFile, err)
	}

	cfg := &Configuration{}
	if err := yaml.Unmarshal(cfgBytes, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal configfile data: %v", err)
	}

	return cfg, nil
}

func GetHostByAddr(hs []*Host, addr string) *Host {
	for _, h := range hs {
		if h.Addr == addr {
			return h
		}
	}
	return nil
}
