package golog

type Conf map[string]string

var (
	Default *Logger
	loggers map[string]*Logger
)

func init() {
	Default = &Logger{
		Name:  "default",
		Level: DEBUG,
	}

	Default.Enable(StdoutAppender())
}

// getting logger instance
// method returns singleton logger instance
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
