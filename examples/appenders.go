package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

func mongo() {
	logger := golog.Default

	// make instance of mongo appender and enable it
	logger.Enable(appenders.Mongo(golog.Conf{
		"host":       "127.0.0.1:27017",
		"db":         "somedb",
		"collection": "logs",
		"username":   "myusername",
		"password":   "mypassword",
	}))

	logger.Debug("some message")
}

func file() {
	logger := golog.Default

	// make instance of file appender and enable it
	logger.Enable(appenders.File(golog.Conf{
		"path": "/path/to/log.txt",
	}))

	logger.Debug("some message")
}
