package golog

// Convinient type for representing appender configuration
type Conf map[string]string

var (
	// instance of default logger
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

// Function for getting logger instance.
// Method returns singleton logger instance.
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
