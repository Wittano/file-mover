package config

import (
	"golang.org/x/exp/maps"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Dirs []Directory
}

type Directory struct {
	Src       []string
	Dest      string
	Recursive bool
}

func LoadConfig(path string) (Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed loaded configuration from path %s causes: %s", path, err)
		return Config{}, err
	}

	var unmarshal map[string]Directory
	if err := toml.Unmarshal(bytes, &unmarshal); err != nil {
		log.Printf("Failed unmarshal configuration causes: %s", err)
		return Config{}, err
	}

	var config Config
	config.Dirs = maps.Values(unmarshal)

	return config, nil
}
