package main

import "github.com/ivpusic/golog"

func simple() {
	// get default logger
	logger := golog.Default

	// log something
	logger.Debug("some message")
}
