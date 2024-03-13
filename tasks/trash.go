package tasks

import (
	"context"
	"fmt"
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"golang.org/x/exp/slices"
	"os"
	"time"
)

func MoveToTrashTask(cancel context.CancelFunc) {
	for _, dir := range setting.Flags.Config().Dirs {
		if dir.MoveToTrash {
			moveFileToTrash(cancel, dir)
		}
	}

	setting.Logger().Debug("Complete 'moveToTrash' task")
}

func moveFileToTrash(cancel context.CancelFunc, dir setting.Directory) {
	paths, err := dir.RealPaths()
	if err != nil {
		setting.Logger().Error("Failed to get file paths", err)
		cancel()
		return
	}

	for _, p := range paths {
		if slices.Contains(dir.Exceptions, p) {
			continue
		}

		if isAfterDateOfRemovingFile(p, dir.After) {
			trashPath, err := dir.TrashDir()
			if err != nil {
				setting.Logger().Error("Failed to find trash directory", err)
				cancel()
				return
			}

			go func(dest string, src string) {
				if err = file.MoveToDestination(dest, src); err != nil {
					setting.Logger().Error(fmt.Sprintf("One of soruce file wasn't moved to destination directory"), err)
					return
				}
			}(trashPath, p)
		}
	}
}

func isAfterDateOfRemovingFile(path string, after uint) bool {
	stat, err := os.Stat(path)
	if err != nil {
		setting.Logger().Warn(fmt.Sprintf("Failed to load file info from %s: %s", path, err))
		return false
	}

	afterTime := time.Hour * 24 * time.Duration(after)
	if afterTime == 0 {
		return true
	}

	return stat.ModTime().Add(afterTime).After(time.Now())
}
