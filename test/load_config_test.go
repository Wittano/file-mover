package test

import (
	"github.com/wittano/file-mover/src/config"
	"log"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := config.LoadConfig("./test_config.toml")
	if err != nil {
		log.Print(err)
		t.FailNow()
	}

	l := len(c)
	if l != 2 {
		t.Errorf("Invalid number of directories. Expected 2, acually %d", l)
	}

	expPictureDir := config.Directory{
		Src:       "/tmp/test",
		Dest:      "/tmp/test2",
		Recursive: true,
	}

	expMusicDir := config.Directory{
		Src:       "/tmp/test2/*.mp4",
		Dest:      "/tmp",
		Recursive: false,
	}

	testDirectory(t, c, "Picture", expPictureDir)
	testDirectory(t, c, "Music", expMusicDir)
}

func testDirectory(t *testing.T, config config.Config, dir string, exp config.Directory) {
	d := config[dir]
	if d != exp {
		t.Errorf("Invalid %s directory config. Expected %v, acually %v", dir, exp, d)
	}
}
