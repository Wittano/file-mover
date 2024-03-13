package setting

import (
	"bytes"
	"github.com/wittano/filebot/internal/test"
	"io"
	"os"
	"testing"
	"time"
)

func TestLogInfo(t *testing.T) {
	path := test.CreateTempFile(t)
	Flags.LogFilePath = path

	logger := Logger()

	const input = "Test"
	logger.Info(input)

	time.Sleep(10 * time.Millisecond)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	res, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(res, []byte(input)) {
		t.Fatalf("Log didn't save")
	}
}

func TestFileStdWriter(t *testing.T) {
	path := test.CreateTempFile(t)
	writer := fileStdWriter{path}

	exp := []byte("test")

	_, err := writer.Write(exp)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	res, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(exp, res) {
		t.Fatalf("Invalid data saved in log file. Expected %v, Actually: %v", exp, res)
	}
}
