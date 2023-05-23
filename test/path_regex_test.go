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

func TestGetPathsFromRegexButFunctionReturnNil(t *testing.T) {
	p := createTempDirWithFile(t)
	exp := strings.Replace(p, "test", "tset", 1)

	paths, err := path.GetPathsFromPattern(exp)
	if err != nil || paths != nil || len(paths) != 0 {
		t.Fatalf("Failed got paths. Expected 1, acually %d", len(paths))
	}
}

func BenchmarkGetPathsFromRegex(b *testing.B) {
	dir := b.TempDir()
	pattern := "test"

	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		b.Fatalf("Failed created temp file. %s", err)
	}

	defer file.Close()

	exp := dir + "/" + pattern

	for i := 0; i < b.N; i++ {
		paths, err := path.GetPathsFromPattern(exp)
		if err == nil && len(paths) != 1 {
			b.Fatalf("Failed got paths. Expected 1, acually %d", len(paths))
		}
	}
}

func createTempDirWithFile(t *testing.T) string {
	dir := t.TempDir()
	pattern := "test"

	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatalf("Failed created temp file. %s", err)
	}
	defer file.Close()

	return dir + "/" + pattern
}
