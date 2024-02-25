package setting

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml/v2"
	"github.com/wittano/filebot/internal/filesystem"
	"github.com/wittano/filebot/path"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// TODO Migrate from TOML to YAML with new features
type Config struct {
	Dirs []Directory `validate:"required"`
}

var config *Config

type Directory struct {
	Src         []string `validate:"required"`
	Dest        string   `validate:"required_if=MoveToTrash false"`
	Recursive   bool
	MoveToTrash bool `validate:"required_without=Dest"`
	After       uint
	Exceptions  []string
}

func (d Directory) RealPaths() (paths []string, err error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	err = v.Struct(d)
	if err != nil {
		return
	}

	for _, exp := range d.Src {
		if d.Recursive {
			paths, err = path.PathsFromPatternRecursive(exp)
		} else {
			paths, err = path.PathsFromPattern(exp)
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

func (d Directory) TrashDir() (trashDir string, err error) {
	if !d.MoveToTrash {
		return "", nil
	}

	dir := filepath.Dir(d.Src[0])
	_, err = os.Stat(dir)
	if err != nil {
		return
	}

	fs, err := filesystem.MountedList()
	if err != nil {
		return
	}

	trashName := ".Trash-" + strconv.Itoa(os.Getuid())

	for _, device := range fs {
		if strings.Contains(dir, device.MountedPoint) && device.MountedPoint != "/" {
			trashDir = filepath.Join(device.MountedPoint, trashName, "files")
			break
		} else if device.MountedPoint == "/" && isUserRoot() {
			trashDir = "/root/.Trash-0/files"
		}
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		return
	}

	if trashDir == "" {
		trashDir = filepath.Join(homeDir, ".local", "share", trashName, "files")
	}

	err = os.MkdirAll(trashDir, 0700)
	if err != nil {
		return
	}

	return
}

func isUserRoot() bool {
	return os.Getuid() == 0
}

func load(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var unmarshal map[string]Directory
	if err := toml.Unmarshal(bytes, &unmarshal); err != nil {
		return nil, err
	}

	if len(unmarshal) == 0 {
		return nil, errors.New("config file is empty")
	}

	config = new(Config)
	config.Dirs = maps.Values(unmarshal)

	v := validator.New(validator.WithRequiredStructEnabled())

	for _, d := range config.Dirs {
		if err = v.Struct(d); err != nil {
			return nil, err
		}
	}

	return config, nil
}
