package cli

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets []HttpTarget `yaml:"targets"`
}

type HttpTarget struct {
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
}

func NewConfig(opts ...ConfigOpt) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

type ConfigOpt func(*Config)

func FromFile(name string) ConfigOpt {
	return func(cfg *Config) {
		f, err := os.Open(name)
		if err != nil {
			fmt.Printf("could not open config: %s\n", err)
			os.Exit(1)
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
			fmt.Printf("could not yaml decode config: %s\n", err)
			os.Exit(1)
		}
	}
}
