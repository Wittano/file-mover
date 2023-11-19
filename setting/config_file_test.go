package setting

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := load("testdata/config.toml")
	if err != nil {
		t.Fatalf("Failed load conf causes %s", err)
	}

	dir := conf.Dirs[0]
	if len(dir.Src) == 1 && dir.Src[0] != "/tmp/test" {
		t.Fatalf("Invalid source paths. Expacted [ '/tmp/test' ], acually %v", dir.Src)
	}

	if dir.Dest != "/tmp/test2" {
		t.Fatalf("Invalid destination path paths. Expacted '/tmp/test', acually %s", dir.Dest)
	}
}

func TestFailedLoadingConfig(t *testing.T) {
	_, err := load("/invalid/path")
	if err == nil {
		t.Fatal("Loaded setting file from invalid path")
	}
}
