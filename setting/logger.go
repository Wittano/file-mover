package setting

import (
	"log"
	"os"
)

type LogLevel uint8

const (
	ALL   LogLevel = 0
	DEBUG LogLevel = 1
	WARN  LogLevel = 2
	INFO  LogLevel = 2
)

var logger FileStdLogger

type fileStdWriter struct {
	filePath string
}

func (f fileStdWriter) Write(data []byte) (n int, err error) {
	if f.filePath != "" {
		go writeToLogFile(f.filePath, data)
	}

	return os.Stdout.Write(data)
}

func writeToLogFile(path string, data []byte) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		log.Fatalf("Failed to log into %s: %s", f.Name(), err)
	}
}

type FileStdLogger struct {
	log   *log.Logger
	level LogLevel
}

func (m FileStdLogger) Fatal(msg string, err error) {
	m.log.Fatalf("[ERROR]: %s. %s\n", msg, err)
}

func (m FileStdLogger) Info(msg string) {
	m.log.Printf("[INFO] %s\n", msg)
}

func (m FileStdLogger) Warn(msg string) {
	if m.level == WARN || m.level == ALL {
		m.log.Printf("[WARN] %s\n", msg)
	}
}

func (m FileStdLogger) Error(msg string, err error) {
	m.log.Printf("[ERROR] %s: %s\n", msg, err)
}

func (m FileStdLogger) Debug(msg string) {
	if m.level == DEBUG || m.level == ALL {
		m.log.Printf("[DEBUG] %s", msg)
	}
}

func Logger() FileStdLogger {
	if logger != (FileStdLogger{}) {
		return logger
	}

	stdLog := log.New(fileStdWriter{Flags.LogFilePath}, "", log.LstdFlags)
	logger = FileStdLogger{stdLog, Flags.LogLevel()}

	return logger
}
