package setting

import (
	"github.com/mitchellh/go-homedir"
	"log"
	"path/filepath"
	"time"
)

type Flag struct {
	ConfigPath     string
	UpdateInterval time.Duration
}

var Flags = Flag{
	DefaultConfigPath(),
	DefaultUpdateInterval(),
}

func (f Flag) Config() *Config {
	if config != nil {
		return config
	}

	c, err := load(f.ConfigPath)
	if err != nil {
		// TODO Add debug logs
		log.Fatalf("Failed to load config file: %s", f.ConfigPath)
	}

	return c
}

func DefaultConfigPath() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Failed to find home directory. %s", err)
	}

	// TODO Change config file from TOML to YAML
	return filepath.Join(homeDir, ".setting", "filebot", "setting.toml")
}

func DefaultUpdateInterval() time.Duration {
	duration, _ := time.ParseDuration("10m")
	return duration
}
