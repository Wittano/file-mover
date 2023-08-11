//go:build linux

package watcher

import (
	"errors"
	"github.com/wittano/filebot/pkg/config"
	"os"
	"testing"
	"time"
)

func TestAddFileToObservable(t *testing.T) {
	conf := createTestConfiguration(t)

	w := NewWatcher()
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
		t.Fatal("File didn't move from started location")
	}

	dest := conf.Dirs[0].Dest
	if _, err := os.Stat(dest); err != nil {
		t.Fatal("File didn't move to destination location")
	}
}

func TestAddFileToObservableRecursive(t *testing.T) {
	conf := createTestConfigurationWithRecursive(t)

	w := NewWatcher()
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
		t.Fatal("File didn't move from started location")
	}

	dest := conf.Dirs[0].Dest
	if _, err := os.Stat(dest); err != nil {
		t.Fatal("File didn't move to destination location")
	}
}

func BenchmarkAddFilesToObservable(b *testing.B) {
	tempDir := b.TempDir()
	secondTempDir := b.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test.mp4")
	if err != nil {
		b.Fatalf("Failed creating temp file: %s", err)
	}
	defer tempFile.Close()

	conf := &config.Config{Dirs: []config.Directory{
		{
			Src:       []string{tempFile.Name()},
			Dest:      secondTempDir,
			Recursive: false,
		},
	}}

	w := NewWatcher()

	for i := 0; i < b.N; i++ {
		w.AddFilesToObservable(conf)
	}

	b.ReportAllocs()
}

func createTestConfiguration(t *testing.T) *config.Config {
	tempDir := t.TempDir()
	secondTempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test.mp4")
	if err != nil {
		t.Fatalf("Failed creating temp file: %s", err)
	}
	defer tempFile.Close()

	return &config.Config{Dirs: []config.Directory{
		{
			Src:       []string{tempFile.Name()},
			Dest:      secondTempDir,
			Recursive: false,
		},
	}}
}

func createTestConfigurationWithRecursive(t *testing.T) *config.Config {
	tempDir := t.TempDir()
	secondTempDir := tempDir + "/test"

	os.Mkdir(secondTempDir, 0777)

	tempFile, err := os.CreateTemp(secondTempDir, "test.mp4")
	if err != nil {
		t.Fatalf("Failed creating temp file: %s", err)
	}
	defer tempFile.Close()

	return &config.Config{Dirs: []config.Directory{
		{
			Src:       []string{tempFile.Name()},
			Dest:      tempDir,
			Recursive: true,
		},
	}}
}
