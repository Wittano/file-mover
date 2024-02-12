package setting

import (
	"github.com/mitchellh/go-homedir"
	"github.com/wittano/filebot/logger"
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

func (f Flag) Config() *Config {
	if config != nil {
		return config
	}

	c, err := load(f.ConfigPath)
	if err != nil {
		Logger().Fatal("Failed to load config file", err)
	}

	return c
}

func (f Flag) LogLevel() logger.LogLevel {
	var level logger.LogLevel

	switch f.LogLevelName {
	case "ALL":
		level = logger.ALL
	case "DEBUG":
		level = logger.DEBUG
	case "WARN":
		level = logger.WARN
	default:
		level = logger.INFO
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
