package config

import (
	"golang.org/x/exp/maps"
	"os"
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
	Src         []string
	Dest        string
	Recursive   bool
	MoveToTrash bool
	After       uint
}

var config *Config

func Get(path string) (*Config, error) {
	if config != nil {
		return config, nil
	}

	return load(path)
}

func load(path string) (*Config, error) {
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
