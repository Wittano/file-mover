package file

import (
	"errors"
	"fmt"
	"github.com/wittano/filebot/path"
	"github.com/wittano/filebot/setting"
	"os"
	"path/filepath"
)

func MoveToDestination(dest string, paths ...string) {
	dest = path.ReplaceEnvVariablesInPath(dest)

	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(dest, 0644)
	}

	for _, src := range paths {
		newPath := filepath.Join(dest, filepath.Base(src))

		if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
			err := os.Rename(src, newPath)
			if err != nil {
				setting.Logger().Error(fmt.Sprintf("Failed to move file from %s to %s", src, newPath), err)
				return
			}

			setting.Logger().Info(fmt.Sprintf("Moved file from %s to %s", src, dest))
		}
	}
}
