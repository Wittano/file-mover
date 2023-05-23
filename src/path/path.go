package path

import (
	"os"
	p "path"
	"regexp"
)

func GetPathsFromPattern(path string) ([]string, error) {
	dir, pattern := p.Split(path)

	regex, err := regexp.Compile("^\\*")
	if err != nil {
		return nil, err
	}

	pattern = "^" + string(regex.ReplaceAll([]byte(pattern), []byte("\\w*"))) + "$"

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, len(files))

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	size := uint(0)
	for _, f := range files {
		if !f.IsDir() && reg.Match([]byte(f.Name())) {
			paths[size] = dir + f.Name()
			size++
		}
	}

	if size == 0 {
		return nil, nil
	}

	return paths[0:size], nil
}
