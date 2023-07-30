package main

import (
	"flag"
	"github.com/wittano/file-mover/pkg/config"
	"time"
)

func parseFlags() config.Flags {
	var flags config.Flags

	flag.StringVar(&flags.ConfigPath, "c", config.GetDefaultConfigPath(), "Configuration file path")

	defaultUpdateInterval, _ := time.ParseDuration("10m")
	flag.DurationVar(
		&flags.UpdateInterval,
		"u",
		defaultUpdateInterval,
		"Time after observable file list will be update",
	)

	flag.Parse()

	return flags
}
