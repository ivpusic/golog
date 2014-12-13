package main

import "github.com/ivpusic/golog"

func formatting() {
	logger := golog.Default

	// will output `some cool number 4`
	// the same you can do for other levels
	logger.Debugf("some %s number %d", "cool", 4)
}
