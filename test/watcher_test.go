package test

import (
	"github.com/wittano/file-mover/src/watcher"
	"testing"
)

func TestAddDirectoryToObservable(t *testing.T) {
	dir := t.TempDir()

	w := watcher.NewWatcher()

	w.AddFileToObservable(dir)

	actually := len(w.Watcher.WatchList())
	if actually != 1 {
		t.Fatalf("Failed to add directory to observable list. Expected %d, acutally %d", 1, actually)
	}
}

func TestFailedAddDirectoryToObservableList(t *testing.T) {
	dir := t.TempDir() + "/failed"

	w := watcher.NewWatcher()

	w.AddFileToObservable(dir)

	actually := len(w.Watcher.WatchList())
	if actually != 0 {
		t.Fatalf("Failed to add directory to observable list. Expected %d, acutally %d", 0, actually)
	}
}
