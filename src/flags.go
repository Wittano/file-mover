package main

import (
	"flag"

	"github.com/wittano/file-mover/src/config"
)

func parseFlags(configPath *string) {
	flag.StringVar(configPath, "c", config.GetDefaultConfigPath(), "Set configuration file path")

	flag.Parse()
}
