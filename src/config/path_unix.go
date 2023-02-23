package config

import "os"

func GetDefaultConfigPath() string {
	home := os.Getenv("HOME")

	return home + "/.config/file-mover/config.ini"
}
