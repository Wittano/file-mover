package cron

import (
	"fmt"
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"log"
	"os"
	"path/filepath"
	"time"
)

// TODO Improve Trash path for other block devices e.g. for NTFS devices
var TrashPath = filepath.Join(os.Getenv("HOME"), ".locale", "share", "Trash", "files")

func moveToTrashTask() {
	c := setting.Flags.GetConfig()

	for _, dir := range c.Dirs {
		if dir.MoveToTrash {
			moveFileToTrash(dir)
		}
	}
}

func moveFileToTrash(dir setting.Directory) {
	for _, p := range dir.Src {
		if isAfterDateOfRemovingFile(p, dir.After) {
			go file.MoveToDestination(TrashPath, p)
		}
	}
}

func isAfterDateOfRemovingFile(path string, after uint) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to load file info from %s: %s", path, err))
		return false
	}

	afterTime := time.Hour * 24 * time.Duration(after)
	if afterTime == 0 {
		return true
	}

	return stat.ModTime().Add(afterTime).After(time.Now())
}
