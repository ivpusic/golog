package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

func main() {
	logger := golog.Default
	logger.Info("some text")

	appLogger := golog.GetLogger("application")

	appLogger.Level = golog.WARN

	fileAppender := appenders.File(golog.Conf{
		"path": "./log.txt",
	})

	appLogger.Enable(fileAppender)

	appLogger.Enable(appenders.Mongo(golog.Conf{
		"port":       "27017",
		"db":         "somedb",
		"collection": "logs",
	}))

	appLogger.Error("log from application logger")
	appLogger.Warn("log from application logger")

	appLogger.Disable(fileAppender)
}
