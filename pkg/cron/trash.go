package cron

import (
	"fmt"
	"github.com/wittano/filebot/cmd/filebot/cmd"
	"github.com/wittano/filebot/pkg/config"
	"github.com/wittano/filebot/pkg/file"
	"log"
	"os"
	"path/filepath"
	"time"
)

var TrashPath = filepath.Join(os.Getenv("HOME"), ".locale", "share", "Trash", "files")

func moveToTrashTask() {
	c, err := config.Get(cmd.Flags.ConfigPath)
	if err != nil {
		log.Println("Failed to get config. " + err.Error())
		return
	}

	for _, dir := range c.Dirs {
		if dir.MoveToTrash {
			moveFileToTrash(dir)
		}
	}
}

func moveFileToTrash(dir config.Directory) {
	for _, path := range dir.Src {
		if isAfterDateOfRemovingFile(path, dir.After) {
			go file.MoveFilesToDestination(TrashPath, path)
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
