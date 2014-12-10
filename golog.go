package golog

import "os"

type Conf map[string]string

var (
	Default *Logger
	loggers map[string]*Logger
)

func init() {
	out = os.Stdout

	Default = &Logger{
		Name:  "default",
		Level: DEBUG,
	}

	Default.Enable(StdoutAppender())
}

func GetLogger(name string) *Logger {
	logger, ok := loggers["foo"]
	if !ok {
		logger = &Logger{
			Name:  name,
			Level: DEBUG,
		}

		logger.Enable(StdoutAppender())
	}

	return logger
}
