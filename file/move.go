package file

import (
	"errors"
	"fmt"
	"github.com/wittano/filebot/path"
	"github.com/wittano/filebot/setting"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	unixOwnerWritePermIndex = 2
	unixGroupWritePermIndex = 5
	unixWritePermission     = "w"
)

func MoveToDestination(dest string, paths ...string) (err error) {
	fixedDest := path.ReplaceEnvVariablesInPath(dest)

	if err = createDestDir(fixedDest); err != nil {
		return
	}

	for _, src := range paths {
		srcInfo, err := os.Stat(src)
		if err != nil {
			setting.Logger().Warn(fmt.Sprintf("Failed read stats from %s", src))
			continue
		}

		if err = checkFilePermissions(srcInfo); err != nil {
			setting.Logger().Error("User hasn't permission to move file", err)
			continue
		}

		newPath := filepath.Join(fixedDest, filepath.Base(src))

		if _, err = os.Stat(src); !errors.Is(err, os.ErrNotExist) {
			err = os.Rename(src, newPath)
			if err != nil {
				setting.Logger().Error(fmt.Sprintf("Failed to move file from %s to %s", src, newPath), err)
				continue
			}

			setting.Logger().Info(fmt.Sprintf("Moved file from %s to %s", src, dest))
		}
	}

	return nil
}

func checkFilePermissions(stat os.FileInfo) error {
	writePermIndex := strings.Index(stat.Mode().String(), unixWritePermission)
	if writePermIndex == -1 {
		return errors.New("no one has write permission")
	}

	switch writePermIndex {
	case unixOwnerWritePermIndex:
		ownerUID := int(stat.Sys().(*syscall.Stat_t).Uid)

		if os.Getuid() != ownerUID {
			return errors.New("user isn't owner")
		}
	case unixGroupWritePermIndex:
		groupID := int(stat.Sys().(*syscall.Stat_t).Gid)

		if os.Getgid() != groupID {
			return errors.New("user doesn't belong to group")
		}
	}

	return nil
}
