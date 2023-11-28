package setting

import (
	"fmt"
	"github.com/wittano/filebot/internal/test"
	"os"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := load("testdata/config.toml")
	if err != nil {
		t.Fatalf("Failed load conf causes %s", err)
	}

	for _, dir := range conf.Dirs {
		t.Run(dir.Dest, func(t *testing.T) {
			if len(dir.Src) == 1 && dir.Src[0] == "" {
				t.Fatalf("Invalid source paths. Expacted [ '/tmp/test' ], acually %v", dir.Src)
			}

			if dir.Dest == "" && !dir.MoveToTrash {
				t.Fatalf("Invalid destination path paths. Expacted '/tmp/test', acually %s", dir.Dest)
			}
		})
	}
}

func TestFailedLoadingConfig(t *testing.T) {
	_, err := load("/invalid/path")
	if err == nil {
		t.Fatal("Loaded setting file from invalid path")
	}
}

func TestGetTrashDir(t *testing.T) {
	destDir := t.TempDir()
	tempFile := test.CreateTempFile(t)

	d := Directory{
		Src:         []string{tempFile},
		Dest:        destDir,
		MoveToTrash: true,
	}

	res, err := d.TrashDir()
	if err != nil {
		t.Fatal(err)
	}

	if res == "" {
		t.Fatal("MoveToTrash field is false")
	}

	var exp string
	if strings.Contains(destDir, "/tmp") {
		exp = fmt.Sprintf("/tmp/.Trash-%d/files", os.Getuid())
	} else {
		exp = fmt.Sprintf("%s/.local/share/.Trash-%d/files", os.Getenv("HOME"), os.Getuid())
	}

	if exp != res {
		t.Fatalf("Trash dir is diffrent. Expected: %s, Actually: %s", exp, res)
	}
}
