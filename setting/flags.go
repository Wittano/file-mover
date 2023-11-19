package setting

import (
	"github.com/mitchellh/go-homedir"
	"log"
	"path/filepath"
	"time"
)

var Flags = NewFlag()

type Flag struct {
	ConfigPath     string
	UpdateInterval time.Duration
}

func NewFlag() Flag {
	return Flag{
		GetDefaultConfigPath(),
		GetDefaultUpdateInterval(),
	}
}

func (f Flag) GetConfig() *Config {
	c, err := GetConfig(f.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to get config file: %s", f.ConfigPath)
	}

	return c
}

func GetDefaultConfigPath() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	// TODO Change config file from TOML to YAML
	return filepath.Join(homeDir, ".setting", "filebot", "setting.toml")
}

func GetDefaultUpdateInterval() time.Duration {
	duration, _ := time.ParseDuration("10m")
	return duration
}
