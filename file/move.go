package file

import (
	"errors"
	"github.com/wittano/filebot/path"
	"github.com/wittano/filebot/setting"
	"os"
	"path/filepath"
)

func MoveToDestination(dest string, paths ...string) {
	dest = path.ReplaceEnvVariablesInPath(dest)

	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		setting.Logger().Errorf("Destination directory %s doesn't exist", err, dest)
		return
	}

	for _, src := range paths {
		newPath := filepath.Join(dest, filepath.Base(src))

		if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
			err := os.Rename(src, newPath)
			if err != nil {
				setting.Logger().Errorf("Failed to move file from %s to %s", err, src, newPath)
				return
			}

			setting.Logger().Info("Moved file from %s to %s", src, dest)
		}
	}
}
