package test

import (
	"os"
	"testing"
)

func CreateTempFile(t *testing.T) string {
	file, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatalf("Failed created temp file. %s", err)
	}
	defer file.Close()

	return file.Name()
}

func CreateNestedTempDirWithFiles(t *testing.T, filename string) (string, string) {
	dir := t.TempDir()
	nestedDir := dir + "test"

	os.Mkdir(nestedDir, 0777)

	file, err := os.CreateTemp(nestedDir, filename)
	if err != nil {
		t.Fatalf("Failed created temp file. %s", err)
	}
	defer file.Close()

	return nestedDir, file.Name()
}
