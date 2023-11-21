//go:build windows

package setting

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfigPath(t *testing.T) {
	dir := DefaultConfigPath()

	if dir != filepath.Join(os.Getenv("USERPROFILE"), ".setting\\filebot\\setting.toml") {
		t.Fatalf("Invalid default setting path. Expected %s\\.setting\\filebot\\setting.toml, Acually: %s", os.Getenv("USERPROFILE"), dir)
	}
}
