package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

type TestStruct struct {
	field string
}

func advanced() {
	logger := golog.Default
	logger.Debug("some message")

	// you can provide additional data to log
	additionalData := &TestStruct{}
	logger.Debug("some message", additionalData)

	// you can require logger instance
	application := golog.GetLogger("application")

	// set log level
	application.Level = golog.WARN

	application.Info("log from application logger")

	appender := appenders.File(golog.Conf{
		"path": "/path/to/log.txt",
	})

	// you can enable appender
	application.Enable(appender)

	// you can disable appender by passing reference
	application.Disable(appender)

	// you can disable appender by passing appender id
	// id is returned with appender.Id() method
	application.Disable("github.com/ivpusic/golog/appender/file")
}
