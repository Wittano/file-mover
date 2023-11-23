package cron

import (
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"time"
)

// TODO Improve Trash path for other block devices e.g. for NTFS devices
var TrashPath = filepath.Join(os.Getenv("HOME"), ".locale", "share", "Trash", "files")

func moveToTrashTask() {
	c := setting.Flags.Config()

	for _, dir := range c.Dirs {
		if dir.MoveToTrash {
			moveFileToTrash(dir)
		}
	}
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
			go file.MoveToDestination(TrashPath, p)
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
