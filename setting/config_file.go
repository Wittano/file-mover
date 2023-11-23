package setting

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/wittano/filebot/path"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"regexp"
)

// TODO Add validation
type Config struct {
	Dirs []Directory
}

type Directory struct {
	Src         []string
	Dest        string
	Recursive   bool
	MoveToTrash bool
	After       uint
	Exceptions  []string
}

func (d Directory) RealPaths() (paths []string, err error) {
	for _, exp := range d.Src {
		if d.Recursive {
			paths, err = path.GetPathFromPatternRecursive(exp)
		} else {
			paths, err = path.GetPathsFromPattern(exp)
		}

		if err != nil {
			Logger().Errorf("Failed get files from pattern '%s'", err, exp)
			return
		}

		paths = append(paths, paths...)
	}

	if d.Exceptions != nil {
		return d.filterRealPaths(paths), nil
	}

	return
}

func (d Directory) filterRealPaths(paths []string) (res []string) {
	for _, p := range paths {
		f, err := os.Stat(p)
		if err != nil {
			continue
		}

		if !f.IsDir() && slices.Contains(d.Exceptions, p) {
			res = append(res, p)
			continue

		}

		for _, exp := range d.Exceptions {
			reg, err := regexp.Compile(exp)
			if err != nil {
				Logger().Warnf("Failed to compile regex: '%s'", nil, exp)
				continue
			}

			filename := filepath.Base(p)
			if exp != filename && !reg.MatchString(filename) {
				res = append(res, p)
			}
		}
	}

	return
}

var config *Config

func load(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var unmarshal map[string]Directory
	if err := toml.Unmarshal(bytes, &unmarshal); err != nil {
		return nil, err
	}

	// TODO Add validator
	config = new(Config)
	config.Dirs = maps.Values(unmarshal)

	return config, nil
}
