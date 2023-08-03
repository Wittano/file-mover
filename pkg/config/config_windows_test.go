//go:build windows

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfigPath(t *testing.T) {
	dir := GetDefaultConfigPath()

	if dir != filepath.Join(os.Getenv("USERPROFILE"), ".config\\filebot\\config.toml") {
		t.Fatalf("Invalid default config path. Expected %s\\.config\\filebot\\config.toml, Acually: %s", os.Getenv("USERPROFILE"), dir)
	}
}
