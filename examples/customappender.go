package main

import "github.com/ivpusic/golog"

type CustomAppender struct {
}

func (s *CustomAppender) Append(log golog.Log) {
	// do something with log
	// for example you can save it to database
	// send it to some service, etc.
}

func (s *CustomAppender) Id() string {
	return "id/of/custom/appender"
}

func customappender() {
	logger := golog.Default

	// make custom appender instance
	appender := &CustomAppender{}

	// enable custom appender
	logger.Enable(appender)

	// log something
	logger.Debug("this will go to custom appender also")
}
