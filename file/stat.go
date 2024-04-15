package file

import (
	"github.com/wittano/filebot/setting"
	"os"
	"path/filepath"
	"syscall"
)

type uID int
type gID int

const (
	rootID  uID = 0
	rootGID gID = 0
)

func uid(path string, ogStat os.FileInfo, config setting.Config) uID {
	baseDir := filepath.Dir(path)

	for _, dir := range config.Dirs {
		if dir.Dest == baseDir {
			uid := uID(os.Getuid())

			if dir.UID > 0 {
				uid = uID(dir.UID)
			} else if dir.IsRoot {
				uid = rootID
			}

			return uid
		}
	}

	return uID(ogStat.Sys().(*syscall.Stat_t).Uid)
}

func gid(path string, ogStat os.FileInfo, config setting.Config) gID {
	baseDir := filepath.Dir(path)

	for _, dir := range config.Dirs {
		if dir.Dest == baseDir {
			uid := gID(os.Getgid())
			if dir.UID > 0 {
				uid = gID(dir.GID)
			} else if dir.IsRoot {
				uid = rootGID
			}

			return uid
		}
	}

	return gID(ogStat.Sys().(*syscall.Stat_t).Gid)
}
