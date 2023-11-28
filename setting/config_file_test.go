package setting

import (
	"fmt"
	"github.com/wittano/filebot/internal/test"
	"golang.org/x/exp/slices"
	"os"
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
	tempFile := test.CreateTempFile(t)

	d := Directory{
		Src:         []string{tempFile},
		Dest:        t.TempDir(),
		MoveToTrash: true,
	}

	res, err := d.TrashDir()
	if err != nil {
		t.Fatal(err)
	}

	if res == "" {
		t.Fatal("MoveToTrash field is false")
	}

	exp := []string{
		fmt.Sprintf("/tmp/.Trash-%d/files", os.Getuid()),
		fmt.Sprintf("%s/.local/share/.Trash-%d/files", os.Getenv("HOME"), os.Getuid()),
	}

	if !slices.Contains(exp, res) {
		t.Fatalf("Trash dir is diffrent. Expected one of them: %v, Actually: %s", exp, res)
	}
}
