package setting

import "github.com/wittano/filebot/logger"

var myLogger logger.Logger

func Logger() logger.Logger {
	if myLogger == nil {
		myLogger = logger.NewLogger(Flags.LogFilePath, Flags.LogLevel())
	}

	return myLogger
}
