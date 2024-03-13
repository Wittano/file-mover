package file

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFileToDestination(t *testing.T) {
	dir := t.TempDir()
	src, err := os.CreateTemp(dir, "tests")
	if err != nil {
		t.Fatal(err)
	}

	resPath := filepath.Join(dir, "test2")
	if err = MoveToDestination(resPath, src.Name()); err != nil {
		t.Fatal(err)
	}

	if _, err = os.Stat(resPath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Failed move file from %s to %s", src.Name(), resPath)
	}
}

func TestMoveFileToDestinationButDestDirNotExist(t *testing.T) {
	dir := t.TempDir()
	src, err := os.CreateTemp(dir, "tests")
	if err != nil {
		t.Fatal(err)
	}

	resPath := filepath.Join(dir, "path", "to", "my", "testFile")
	if err = MoveToDestination(resPath, src.Name()); err != nil {
		t.Fatal(err)
	}

	if _, err = os.Stat(resPath); err != nil {
		t.Fatalf("Failed move file from %s to %s", src.Name(), resPath)
	}
}
