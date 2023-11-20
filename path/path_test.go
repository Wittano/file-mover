package path

import (
	"github.com/wittano/filebot/internal/test"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetPathsFromRegex(t *testing.T) {
	exp := test.CreateTempFile(t)

	paths, err := GetPathsFromPattern(exp)
	if err == nil && len(paths) != 1 && paths[0] == exp {
		t.Fatalf("Failed got paths. Expected 1, acually %d", len(paths))
	}
}

func TestGetPathsFromRegexButRegexStartFromStar(t *testing.T) {
	p := test.CreateTempFile(t)
	exp := strings.Replace(p, "test", "*est", 1)

	paths, err := GetPathsFromPattern(exp)
	if err == nil && len(paths) != 1 && paths[0] == exp {
		t.Fatalf("Failed got paths. Expected 1, acually %d", len(paths))
	}
}

func TestGetPathsFromRegexButFunctionReturnNil(t *testing.T) {
	p := test.CreateTempFile(t)
	exp := strings.Replace(p, "test", "tset", 1)

	paths, err := GetPathsFromPattern(exp)
	if err != nil || paths != nil || len(paths) != 0 {
		t.Fatalf("Failed got paths. Expected 1, acually %d", len(paths))
	}
}

func TestGetPathsFromRegexRecursive(t *testing.T) {
	_, expFilename := test.CreateNestedTempDirWithFiles(t, "test.mp4")
	dir := filepath.Dir(expFilename)

	paths, err := GetPathFromPatternRecursive(dir + "*.mp4*")
	if err != nil || len(paths) != 1 {
		t.Fatalf("Failed got paths. Expected 1, acually %d. Error: %s", len(paths), err)
	}

	if expFilename == filepath.Base(paths[0]) {
		t.Fatalf("Expected file not found")
	}
}

func TestGetPathsFromRegexRecursiveButFunctionReturnNil(t *testing.T) {
	expDir, _ := test.CreateNestedTempDirWithFiles(t, "test")
	dir := strings.Replace(expDir, "test", "tset", 1)

	paths, err := GetPathFromPatternRecursive(dir)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(paths) != 0 {
		t.Fatalf("Failed got paths. Expected 0, acually %d.", len(paths))
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

	exp := file.Name()

	for i := 0; i < b.N; i++ {
		GetPathsFromPattern(exp)
	}

	b.ReportAllocs()
}
