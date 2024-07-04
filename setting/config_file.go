package setting

import (
	"errors"
	"fmt"
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

type Config struct {
	Dirs []Directory `validate:"required"`
}

var config Config

type Directory struct {
	Src         []string `validate:"required"`
	Dest        string   `validate:"required_if=MoveToTrash false"`
	Recursive   bool
	MoveToTrash bool `validate:"required_without=Dest"`
	After       uint
	Exceptions  []string
	UID         uint32
	GID         uint32
	IsRoot      bool
}

func (d Directory) RealPaths() (paths []string, err error) {
	for _, exp := range d.Src {
		if d.Recursive {
			paths, err = path.PathsFromPatternRecursive(exp)
		} else {
			paths, err = path.PathsFromPattern(exp)
		}

		if err != nil {
			Logger().Error(fmt.Sprintf("Failed get files from pattern '%s'", exp), err)
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
				Logger().Warn(fmt.Sprintf("Failed to compile regex: '%s'", exp))
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

	if trashDir, err = d.defaultHomeTrash(); err == nil {
		return
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

	name := ".Trash-" + strconv.Itoa(os.Getuid())

	for _, device := range fs {
		if strings.Contains(dir, device.MountedPoint) && device.MountedPoint != "/" {
			trashDir = filepath.Join(device.MountedPoint, name, "files")
			break
		} else if device.MountedPoint == "/" && isUserRoot() {
			trashDir = "/root/.Trash-0/files"
		}
	}

	if trashDir == "" {
		homeDir, err := homedir.Dir()
		if err != nil {
			return "", err
		}

		trashDir = filepath.Join(homeDir, ".local", "share", name, "files")
	}

	if err = os.MkdirAll(trashDir, 0700); err != nil {
		return "", err
	}

	return
}

func (d Directory) defaultHomeTrash() (trashDir string, err error) {
	var isHomeDir bool

	for _, src := range d.Src {
		if strings.HasPrefix(path.ReplaceEnvVariablesInPath(src), "/home") {
			isHomeDir = true
			break
		}
	}

	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	if isHomeDir {
		trashDir = filepath.Join(dir, ".local", "share", "Trash", "files")
	} else {
		err = errors.New("settings: sources doesn't refer to home directory")
	}

	return
}

func isUserRoot() bool {
	return os.Getuid() == 0
}

func load(path string) (Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var unmarshal map[string]Directory
	if err := toml.Unmarshal(bytes, &unmarshal); err != nil {
		return Config{}, err
	}

	if len(unmarshal) == 0 {
		return Config{}, errors.New("config file is empty")
	}

	config.Dirs = maps.Values(unmarshal)

	v := validator.New(validator.WithRequiredStructEnabled())

	for _, d := range config.Dirs {
		if err = v.Struct(d); err != nil {
			return Config{}, err
		}
	}

	return config, nil
}
