package path

import (
	"os"
	"path/filepath"
)

func GetPathsFromPattern(src string) ([]string, error) {
	reg, err := GetPathRegex(src)
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Dir(src)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files))

	for _, f := range files {
		if !f.IsDir() && reg.MatchString(f.Name()) {
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

	var (
		paths = make([]string, len(files))
		size  = uint(0)
	)

	for _, f := range files {
		var path []string

		if f.IsDir() {
			path, err = GetPathFromPatternRecursive(dir + f.Name())
		} else {
			path, err = GetPathsFromPattern(filepath.Join(dir, f.Name()))
		}

		if err != nil {
			return nil, err
		}

		paths = append(paths[0:size], path...)
		size = uint(len(paths))
	}

	if size == 0 {
		return nil, nil
	}

	return paths[0:size], nil
}