package tasks

import (
	"context"
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"golang.org/x/exp/slices"
	"os"
	"time"
)

func MoveToTrashTask(cancel context.CancelFunc) {
	c := setting.Flags.Config()

	for _, dir := range c.Dirs {
		if dir.MoveToTrash {
			moveFileToTrash(cancel, dir)
		}
	}

	setting.Logger().Debug("Complete 'moveToTrash' task", nil)
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

			go file.MoveToDestination(trashPath, p)
		}
	}
}

func isAfterDateOfRemovingFile(path string, after uint) bool {
	stat, err := os.Stat(path)
	if err != nil {
		setting.Logger().Warn("Failed to load file info from "+path, err)
		return false
	}

	afterTime := time.Hour * 24 * time.Duration(after)
	if afterTime == 0 {
		return true
	}

	return stat.ModTime().Add(afterTime).After(time.Now())
}
