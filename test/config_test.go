package test

import (
	"github.com/wittano/fmanager/pkg/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := config.Load("./testdata/config.toml")
	if err != nil {
		t.Fatalf("Failed load conf causes %s", err)
	}

	if len(conf.Dirs) != 2 {
		t.Fatalf("Number of loaded configuration directories is invalid. Expected 2, acually %d", len(conf.Dirs))
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
	_, err := config.Load("/invalid/path")
	if err == nil {
		t.Fatal("Loaded config file from invalid path")
	}
}
