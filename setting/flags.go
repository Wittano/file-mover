package setting

import (
	"github.com/mitchellh/go-homedir"
	"log"
	"path/filepath"
	"time"
)

type Flag struct {
	ConfigPath     string        `validation:"required,filepath"`
	UpdateInterval time.Duration `validation:"required,min=1m"`
	LogFilePath    string        `validation:"filepath"`
	LogLevelName   string        `validation:"required_with=LogFilePath"`
}

var Flags = Flag{
	DefaultConfigPath(),
	DefaultUpdateInterval(),
	"",
	"",
}

func (f Flag) Config() (Config, error) {
	if config.Dirs != nil {
		return config, nil
	}

	return load(f.ConfigPath)
}

func (f Flag) LogLevel() LogLevel {
	var level LogLevel

	switch f.LogLevelName {
	case "ALL":
		level = ALL
	case "DEBUG":
		level = DEBUG
	case "WARN":
		level = WARN
	default:
		level = INFO
	}

	return level
}

func DefaultConfigPath() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Failed to find home directory. %s", err)
	}

	return filepath.Join(homeDir, ".config", "filebot", "setting.toml")
}

func DefaultUpdateInterval() time.Duration {
	duration, _ := time.ParseDuration("10m")
	return duration
}
