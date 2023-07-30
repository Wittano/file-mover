package config

import (
	"golang.org/x/exp/maps"
	"log"
	"os"
	"path"
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
	return path.Join(os.Getenv("HOME"), ".config", "file_mover", "config.toml")
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
		log.Printf("Failed loaded configuration from path %s causes: %s", path, err)
		return nil, err
	}

	var unmarshal map[string]Directory
	if err := toml.Unmarshal(bytes, &unmarshal); err != nil {
		log.Printf("Failed unmarshal configuration causes: %s", err)
		return nil, err
	}

	config = new(Config)
	config.Dirs = maps.Values(unmarshal)

	return config, nil
}
