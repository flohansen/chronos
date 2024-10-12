package scrape

import "time"

var DefaultConfig = Config{
	Targets: []HttpTarget{
		{
			Scheme:  "http",
			Host:    "127.0.0.1:3000",
			Path:    "/metrics",
			Timeout: 10 * time.Second,
		},
	},
}

type Config struct {
	Targets []HttpTarget `yaml:"targets"`
}

type HttpTarget struct {
	Scheme  string        `yaml:"scheme"`
	Host    string        `yaml:"host"`
	Path    string        `yaml:"path"`
	Timeout time.Duration `yaml:"timeout"`
}
