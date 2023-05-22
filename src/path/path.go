package path

import (
	"os"
	p "path"
	"regexp"
	"strings"
)

func GetPathsFromPattern(path string) ([]string, error) {
	dir, pattern := p.Split(path)

	if strings.Contains(pattern, "*") {
		pattern = "\\w" + pattern
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() && reg.Match([]byte(f.Name())) {
			paths = append(paths, dir+f.Name())
		}
	}

	return paths, err
}
