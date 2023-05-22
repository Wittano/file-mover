package main

import (
	"flag"
	"github.com/wittano/file-mover/src/path"
)

func parseFlags(configPath *string) {
	flag.StringVar(configPath, "c", path.GetDefaultConfigPath(), "Set configuration file path")

	flag.Parse()
}
