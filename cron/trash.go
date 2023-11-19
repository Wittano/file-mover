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
	for _, paths := range dir.Src {
		if isAfterDateOfRemovingFile(paths, dir.After) {
			go file.MoveToDestination(TrashPath, paths)
		}
	}
}

func isAfterDateOfRemovingFile(path string, after uint) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to load file info from %s: %s", path, err))
		return false
	}

	return stat.ModTime().Add(time.Hour * 24 * time.Duration(after)).After(time.Now())
}
