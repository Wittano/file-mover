package test

import (
	"path/filepath"
)

func LoadTestData(name string) string {
	return filepath.Join("..", "..", "test", "testdata", name)
}
