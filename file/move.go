package file

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func MoveToDestination(dest string, paths ...string) {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		log.Printf("Destination directory %s doesn't exist", dest)
		return
	}

	for _, src := range paths {
		newPath := filepath.Join(dest, filepath.Base(src))

		if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
			err := os.Rename(src, newPath)
			if err != nil {
				log.Printf("Failed to move file from %s to %s. %s", src, newPath, err)
				return
			}

			log.Printf("Moved file from %s to %s", src, dest)
		}
	}
}
