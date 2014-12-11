package main

import "github.com/ivpusic/golog"

func simple() {
	// get default logger
	logger := golog.Default

	// default level for all loggers is DEBUG
	// you can easily change it it you want
	logger.Level = golog.DEBUG

	// log something
	logger.Debug("some message")
}

func main() {
	simple()
}
