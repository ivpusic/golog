package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appender"

func main() {
	logger := golog.Default
	logger.Info("some text")

	appLogger := golog.GetLogger("application")

	appLogger.Level = golog.WARN

	fileAppender := appender.GetFileAppender(appender.Conf{
		"path": "./log.txt",
	})

	appLogger.Enable(fileAppender)

	appLogger.Error("log from application logger")
	appLogger.Warn("log from application logger")

	appLogger.Disable(fileAppender)
}
