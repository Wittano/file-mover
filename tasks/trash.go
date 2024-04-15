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

func MoveToTrashTask(ctx context.Context) (err error) {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	config, err := setting.Flags.Config()
	if err != nil {
		return err
	}

	for _, dir := range config.Dirs {
		if dir.MoveToTrash {
			if err = moveFileToTrash(dir); err != nil {
				return
			}
		}
	}

	setting.Logger().Debug("Complete 'moveToTrash' taskDetails")

	return
}

func moveFileToTrash(dir setting.Directory) error {
	paths, err := dir.RealPaths()
	if err != nil {
		return err
	}

	for _, p := range paths {
		if slices.Contains(dir.Exceptions, p) {
			continue
		}

		if isAfterDateOfRemovingFile(p, dir.After) {
			trashPath, err := dir.TrashDir()
			if err != nil {
				return err
			}

			go func(dest string, src string) {
				if err = file.MoveToDestination(dest, src); err != nil {
					setting.Logger().Error(fmt.Sprintf("One of soruce file wasn't moved to destination directory"), err)
					return
				}
			}(trashPath, p)
		}
	}

	return nil
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
