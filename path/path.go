package path

import (
	"os"
	"path/filepath"
	"regexp"
)

func GetPathsFromPattern(src string) ([]string, error) {
	reg, err := getPathRegex(src)
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Dir(src)

	if f, _ := os.Stat(dirPath); f != nil && !f.IsDir() {
		if reg.Match([]byte(f.Name())) {
			return []string{f.Name()}, err
		}
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() && reg.Match([]byte(f.Name())) {
			paths = append(paths, filepath.Join(dirPath, f.Name()))
		}
	}

	if len(paths) == 0 {
		return nil, nil
	}

	return paths, nil
}

func GetPathFromPatternRecursive(path string) ([]string, error) {
	dir, pattern := filepath.Split(path)
	if !isFilePathIsRegex(pattern) {
		dir = path
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return GetPathsFromPattern(dir)
	}

	paths := make([]string, len(files))

	size := uint(0)
	for _, f := range files {
		if f.IsDir() {
			recPath, err := GetPathFromPatternRecursive(dir + f.Name())
			if err != nil {
				return nil, err
			}

			paths = append(paths[0:size], recPath...)
			size = uint(len(paths))
		} else {
			path, err := GetPathsFromPattern(filepath.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}

			if len(path) > 0 {
				paths = append(paths[0:size], path...)
				size++
			}
		}
	}

	if size == 0 {
		return nil, nil
	}

	return paths[0:size], nil
}

func getPathRegex(src string) (*regexp.Regexp, error) {
	pattern := filepath.Base(src)

	reg, err := regexp.Compile("^\\*")
	if err != nil {
		return nil, err
	}

	pattern = "^" + string(reg.ReplaceAll([]byte(pattern), []byte("\\w*"))) + "$"

	return regexp.Compile(pattern)
}

func isFilePathIsRegex(reg string) bool {
	specialChars := "*+?|[]{}()"

	for _, specChar := range specialChars {
		for _, char := range reg {
			if char == specChar {
				return true
			}
		}
	}

	return false
}
