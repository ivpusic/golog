package golog

// Convinient type for representing appender configuration
type Conf map[string]string

var (
	// instance of default logger
	Default *Logger
	loggers map[string]*Logger
)

func init() {
	loggers = map[string]*Logger{}
	Default = GetLogger("default")
}

// Function for getting logger instance.
// Method returns singleton logger instance.
func GetLogger(name string) *Logger {
	logger, ok := loggers[name]
	if !ok {
		logger = &Logger{
			Name:  name,
			Level: DEBUG,
		}

		logger.Enable(StdoutAppender())
		logger.normalizeName()

		// recalculate names
		curnamelen = len(logger.Name)
		for _, _logger := range loggers {
			_logger.normalizeName()
		}

		loggers[name] = logger
	}

	return logger
}

// Will disable all logs comming from logger with provided name
func Disable(name string) {
	logger := loggers[name]
	if logger == nil {
		Default.Warn("cannot find logger " + name)
		return
	}

	logger.disabled = true
}

// Will enable all logs comming to logger with provided name
func Enable(name string) {
	logger := loggers[name]
	if logger == nil {
		Default.Warn("cannot find logger " + name)
		return
	}

	logger.disabled = false
}
