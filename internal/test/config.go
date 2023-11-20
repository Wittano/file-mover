package test

import (
	"github.com/wittano/filebot/setting"
	"os"
	"testing"
)

func CreateTestConfiguration(t *testing.T) *setting.Config {
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

func CreateTestConfigurationWithRecursive(t *testing.T) *setting.Config {
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
