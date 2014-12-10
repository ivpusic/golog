package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

func main() {
	logger := golog.Default
	logger.Info("some text")

	appLogger := golog.GetLogger("application")

	appLogger.Level = golog.DEBUG

	fileAppender := appenders.File(golog.Conf{
		"path": "./log.txt",
	})

	appLogger.Enable(fileAppender)

	appLogger.Enable(appenders.Mongo(golog.Conf{
		"host":       "127.0.0.1:27017",
		"db":         "somedb",
		"collection": "logs",
	}))

	appLogger.Error("log from application logger")
	appLogger.Warn("log from application logger")

	appLogger.Disable(fileAppender)
}
