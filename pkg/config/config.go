package config

import (
	"github.com/mitchellh/go-homedir"
	"golang.org/x/exp/maps"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Flags struct {
	ConfigPath     string
	UpdateInterval time.Duration
}

type Config struct {
	Dirs []Directory
}

type Directory struct {
	Src       []string
	Dest      string
	Recursive bool
}

var config *Config

func GetDefaultConfigPath() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".config", "filebot", "config.toml")
}

func GetDefaultUpdateInterval() time.Duration {
	duration, _ := time.ParseDuration("10m")
	return duration
}

func Get(path string) (*Config, error) {
	if config != nil {
		return config, nil
	}

	return Load(path)
}

func Load(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var unmarshal map[string]Directory
	if err := toml.Unmarshal(bytes, &unmarshal); err != nil {
		return nil, err
	}

	config = new(Config)
	config.Dirs = maps.Values(unmarshal)

	return config, nil
}
