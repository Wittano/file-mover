package path

import (
	"os"
	"path"
)

func GetDefaultConfigPath() string {
	return path.Join(os.Getenv("HOME"), ".config", "file-mover", "config.toml")
}
