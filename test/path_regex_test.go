package test

import (
	"github.com/wittano/file-mover/src/path"
	"os"
	"strings"
	"testing"
)

func TestGetPathsFromRegex(t *testing.T) {
	exp := createTempDirWithFile(t)

	paths, err := path.GetPathsFromPattern(exp)
	if err == nil && len(paths) != 1 {
		t.Fatalf("Failed got paths. Expected 1, acually %d", len(paths))
	}
}

func TestGetPathsFromRegexButRegexStartFromStar(t *testing.T) {
	p := createTempDirWithFile(t)
	exp := strings.Replace(p, "test", "*est", 1)

	paths, err := path.GetPathsFromPattern(exp)
	if err == nil && len(paths) != 1 {
		t.Fatalf("Failed got paths. Expected 1, acually %d", len(paths))
	}
}

func createTempDirWithFile(t *testing.T) string {
	dir := t.TempDir()
	pattern := "test"

	_, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatalf("Failed created temp file. %s", err)
	}

	return dir + "/" + pattern
}
