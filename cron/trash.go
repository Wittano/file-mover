package cron

import (
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"golang.org/x/exp/slices"
	"os"
	"time"
)

func moveToTrashTask() {
	c := setting.Flags.Config()

	for _, dir := range c.Dirs {
		if dir.MoveToTrash {
			moveFileToTrash(dir)
		}
	}

	setting.Logger().Debug("Complete 'moveToTrash' task", nil)
}

func moveFileToTrash(dir setting.Directory) {
	paths, err := dir.RealPaths()
	if err != nil {
		setting.Logger().Error("Failed to get file paths", err)
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
