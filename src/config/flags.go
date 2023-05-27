package config

import (
	"flag"
	"github.com/wittano/file-mover/src/path"
	"time"
)

type FlagConfig struct {
	ConfigPath     string
	UpdateInterval time.Duration
}

func ParseFlags() FlagConfig {
	var config FlagConfig

	flag.StringVar(&config.ConfigPath, "c", path.GetDefaultConfigPath(), "Configuration file path")

	defaultUpdateInterval, _ := time.ParseDuration("10m")
	flag.DurationVar(
		&config.UpdateInterval,
		"u",
		defaultUpdateInterval,
		"Time after observable file list will be update",
	)

	flag.Parse()

	return config
}
