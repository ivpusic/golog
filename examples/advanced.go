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

	// you can disable some logger completely
	// you have to provide logger name in order to disable it
	golog.Disable("github.com/someuser/somelib")

	// all loggers are enabled by default
	// if you have case that at some point you disable it,
	// and later you want to enable it again, you can use this method
	golog.Enable("github.com/someuser/somelib")
}
