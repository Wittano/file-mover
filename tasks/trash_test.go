package tasks

import (
	"github.com/wittano/filebot/internal/test"
	"github.com/wittano/filebot/setting"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsAfterDateOfRemovingFile(t *testing.T) {
	f := test.CreateTempFile(t)

	res := isAfterDateOfRemovingFile(f, 1)

	if res == false {
		t.Fatalf("File %s is older then should be", f)
	}
}

func TestIsAfterDateOfRemovingFileButFileDoesNotExist(t *testing.T) {
	res := isAfterDateOfRemovingFile("/path/to/non/existing/file", 1)

	if res {
		t.Fatalf("Non-exisiting file was found")
	}
}

func TestIsAfterDateOfRemovingFileButAfterTimeIsEqualZero(t *testing.T) {
	f := test.CreateTempFile(t)

	res := isAfterDateOfRemovingFile(f, 0)

	if res == false {
		t.Fatalf("File %s shouldn't move. Function isAfterDateOfRemovingFile returned true", f)
	}
}

func TestMoveFileToTrash(t *testing.T) {
	f := test.CreateTempFile(t)
	dir := setting.Directory{
		Src:         []string{f},
		MoveToTrash: true,
		After:       0,
	}

	moveFileToTrash(dir)

	time.Sleep(10 * time.Millisecond)

	if _, err := os.Stat(f); err == nil {
		t.Fatalf("File %s didn't move from original source", f)
	}

	trashDir, err := dir.TrashDir()
	if err != nil {
		t.Fatal(err)
	}

	newFilePath := filepath.Join(trashDir, filepath.Base(f))

	if _, err := os.Stat(newFilePath); err != nil {
		t.Fatalf("File %s didn't move to destination directory. %s", f, err)
	}
}
