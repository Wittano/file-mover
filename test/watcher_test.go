package test

import (
	"errors"
	"github.com/wittano/file-mover/src/config"
	"github.com/wittano/file-mover/src/watcher"
	"os"
	"testing"
	"time"
)

func TestAddFileToObservable(t *testing.T) {
	conf := createTestConfiguration(t)

	w := watcher.NewWatcher()
	w.AddFilesToObservable(conf)

	if len(w.WatchList()) != 1 {
		t.Fatalf("Invalid number of watched files. Expected 1, actually %d", len(w.Watcher.WatchList()))
	}

	src := conf.Dirs[0].Src[0]
	if w.WatchList()[0] != src {
		t.Fatalf(
			"Invalid path was added to observation list. Expected %s, actually %s",
			src,
			w.WatchList()[0],
		)
	}

	duration, _ := time.ParseDuration("0.5s")
	time.Sleep(duration)
	if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("File didn't move to destination location")
	}

	dest := conf.Dirs[0].Dest
	if _, err := os.Stat(dest); err != nil {
		t.Fatal("File didn't move to destination location")
	}
}

func createTestConfiguration(t *testing.T) config.Config {
	tempDir := t.TempDir()
	secondTempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test.mp4")
	if err != nil {
		t.Fatalf("Failed creating temp file: %s", err)
	}
	defer tempFile.Close()

	return config.Config{Dirs: []config.Directory{
		{
			Src:       []string{tempFile.Name()},
			Dest:      secondTempDir,
			Recursive: false,
		},
	}}
}
