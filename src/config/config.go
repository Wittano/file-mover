package config

import (
	"github.com/pelletier/go-toml"
	"os"
)

type Config map[string]Directory

type Directory struct {
	Src       string
	Dest      string
	Recursive bool
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string]Directory

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
