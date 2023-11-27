package linux

import (
	"bufio"
	"os"
	"strings"
)

type FileSystem struct {
	Device       string
	MountedPoint string
	Type         string
}

func MountedList() (fss []FileSystem, err error) {
	f, err := os.Open("/etc/mtab")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		buf := strings.Split(scanner.Text(), " ")

		fs := FileSystem{
			Device:       buf[0],
			MountedPoint: buf[1],
			Type:         buf[2],
		}

		fss = append(fss, fs)
	}

	return
}
