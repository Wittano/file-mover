package file

import (
	"errors"
	"os"
)

func createDestDir(path string) (err error) {
	_, err = os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		return
	}

	const defaultDirPermission = 0755
	return os.MkdirAll(path, defaultDirPermission)
}
