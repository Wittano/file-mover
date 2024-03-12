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
	MoveToDestination(resPath, src.Name())

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
	MoveToDestination(resPath, src.Name())

	if _, err = os.Stat(resPath); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Failed move file from %s to %s", src.Name(), resPath)
	}
}

func TestMoveFiletToDestinationButDestHasTildaChar(t *testing.T) {
	dir := t.TempDir()
	src, err := os.CreateTemp(dir, "tests")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(src.Name())

	const homeTestDir = "~/testDir"

	defer os.Remove(homeTestDir)
	MoveToDestination(homeTestDir, src.Name())

	if _, err = os.Stat(filepath.Join(homeTestDir, filepath.Base(src.Name()))); err != nil {
		t.Fatal(err)
	}
}
