package logger

import (
	"fmt"
	"log"
	"os"
)

type LogLevel uint

const (
	ALL   LogLevel = 0
	DEBUG LogLevel = 1
	WARN  LogLevel = 2
	INFO  LogLevel = 2
)

const errorPrefix = "[ERROR] "

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, err error)
	Warnf(msg string, err error, args ...any)
	Error(msg string, err error)
	Errorf(msg string, err error, args ...any)
	Debug(msg string, err error)
	Fatal(msg string, err error)
}

type fileStdWriter struct {
	filePath string
}

func (f fileStdWriter) Write(p []byte) (n int, err error) {
	if f.filePath != "" {
		go writeToLogFile(f.filePath, p)
	}

	return os.Stdout.Write(p)
}

func writeToLogFile(path string, p []byte) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open log file")
	}
	defer f.Close()

	_, err = f.Write(p)
	if err != nil {
		log.Fatalf("Failed to log into file. %s", err)
	}
}

type myLogger struct {
	log   *log.Logger
	level LogLevel
}

func (m myLogger) Fatal(msg string, err error) {
	m.log.Fatalf("[ERROR]: %s. %s", msg, err)
}

func (m myLogger) Info(msg string, args ...any) {
	m.log.Println(fmt.Sprintf("[INFO] %s", fmt.Sprintf(msg, args...)))
}

func (m myLogger) Warn(msg string, err error) {
	if m.level == WARN || m.level == ALL {
		m.log.Println(appendError(fmt.Sprintf("[WARN] %s", msg), err))
	}
}

func (m myLogger) Warnf(msg string, err error, args ...any) {
	if m.level == WARN || m.level == ALL {
		m.log.Println(appendError(fmt.Sprintf("[WARN] %s", fmt.Sprintf(msg, args...)), err))
	}
}

func (m myLogger) Error(msg string, err error) {
	m.log.Println(appendError(errorPrefix+msg, err))
}

func (m myLogger) Errorf(msg string, err error, args ...any) {
	m.log.Println(appendError(fmt.Sprintf(errorPrefix+fmt.Sprintf(msg, args)), err))
}

func (m myLogger) Debug(msg string, err error) {
	if m.level == DEBUG || m.level == ALL {
		m.log.Println(appendError(errorPrefix+msg, err))
	}
}

func appendError(msg string, err error) string {
	if err != nil {
		return msg + ", " + err.Error()
	}

	return msg
}

func NewLogger(filePath string, level LogLevel) Logger {
	logger := log.New(fileStdWriter{filePath}, "[FileBot]", log.LstdFlags|log.Lmsgprefix)

	return myLogger{logger, level}
}
