//go:build linux

package watcher

import (
	"errors"
	"github.com/wittano/filebot/setting"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAddFileToObservable(t *testing.T) {
	conf := createTestConfiguration(t)

	w := NewWatcher()
	w.AddFilesToObservable(*conf)

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

	dest := filepath.Join(conf.Dirs[0].Dest, filepath.Base(src))
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("File %s didn't move to destination location", dest)
	}
}

func TestAddFileToObservableButDestinationPathHasEnvVariable(t *testing.T) {
	conf := createTestConfiguration(t)

	os.Setenv("TEST", filepath.Dir(t.TempDir()))

	conf.Dirs[0].Dest = strings.ReplaceAll(conf.Dirs[0].Dest, filepath.Dir(t.TempDir()), "$TEST")

	w := NewWatcher()
	w.AddFilesToObservable(*conf)

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

	conf.Dirs[0].Dest = strings.ReplaceAll(conf.Dirs[0].Dest, "$TEST", filepath.Dir(t.TempDir()))

	dest := filepath.Join(conf.Dirs[0].Dest, filepath.Base(src))
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("File %s didn't move to destination location", dest)
	}
}

func TestAddFileToObservableButFilesAreInExceptions(t *testing.T) {
	conf := createTestConfiguration(t)

	w := NewWatcher()
	w.AddFilesToObservable(*conf)

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
	conf.Dirs[0].Exceptions = []string{filepath.Base(conf.Dirs[0].Src[0])}

	w := NewWatcher()
	w.AddFilesToObservable(*conf)

	if len(w.WatchList()) != 0 {
		t.Fatalf("Invalid number of watched files. Expected 0, actually %d", len(w.Watcher.WatchList()))
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

	conf := &setting.Config{Dirs: []setting.Directory{
		{
			Src:       []string{tempFile.Name()},
			Dest:      secondTempDir,
			Recursive: false,
		},
	}}

	w := NewWatcher()

	for i := 0; i < b.N; i++ {
		w.AddFilesToObservable(*conf)
	}

	b.ReportAllocs()
}

func createTestConfiguration(t *testing.T) *setting.Config {
	tempDir := t.TempDir()
	secondTempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test.mp4")
	if err != nil {
		t.Fatalf("Failed creating temp file: %s", err)
	}
	defer tempFile.Close()

	return &setting.Config{Dirs: []setting.Directory{
		{
			Src:       []string{tempFile.Name()},
			Dest:      secondTempDir,
			Recursive: false,
		},
	}}
}

func createTestConfigurationWithRecursive(t *testing.T) *setting.Config {
	tempDir := t.TempDir()
	secondTempDir := tempDir + "/test"

	os.Mkdir(secondTempDir, 0777)

	tempFile, err := os.CreateTemp(secondTempDir, "test.mp4")
	if err != nil {
		t.Fatalf("Failed creating temp file: %s", err)
	}
	defer tempFile.Close()

	return &setting.Config{Dirs: []setting.Directory{
		{
			Src:       []string{tempFile.Name()},
			Dest:      tempDir,
			Recursive: true,
		},
	}}
}
