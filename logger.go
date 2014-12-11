package golog

import (
	"os"
	"strings"
	"time"
)

var (
	DEBUG = Level{
		value: 10,
		color: "blue",
		icon:  "★",
		Name:  "DEBUG",
	}

	INFO = Level{
		value: 20,
		color: "green",
		icon:  "♥",
		Name:  "INFO",
	}

	WARN = Level{
		value: 30,
		color: "yellow",
		icon:  "\u26A0",
		Name:  "WARN",
	}

	ERROR = Level{
		value: 40,
		color: "red",
		Name:  "ERROR",
		icon:  "✖",
	}

	PANIC = Level{
		value: 50,
		color: "black",
		icon:  "☹",
		Name:  "PANIC",
	}

	// limit when logger name will be normalized
	// normalized names are shown in console using stdout appender
	namelen = 18

	// supported name separators
	separators []byte = []byte{'/', '.', '-'}
)

// Representing log level
type Level struct {
	// level priority value
	// bigger number has bigger priority
	value int

	// color which will used by stdout appender
	// github.com/ivpusic/go-clicolor
	color string

	// ascii icon of level
	icon string

	// level name
	Name string
}

// Representing one Log instance
type Log struct {
	// date and time of log
	Time time.Time `json:"time"`

	// logged message
	Message string `json:"message"`

	// log level
	Level Level `json:"level"`

	// additional data sent to log
	// this part should be handled by appenders
	// appender can decide to ignore data or to store it on specific way
	Data []interface{} `json:"data"`

	// id of process which made log
	Pid int `json:"pid"`

	// logger instance
	Logger *Logger `json:"logger"`
}

// Representing one logger instance
// Logger can have multiple appenders, it can enable it,
// or disable it. Also you can define level which will be specific to this logger.
type Logger struct {
	// list of appenders
	appenders []Appender

	// is logged disabled
	disabled bool

	// name of logger
	// logger name will be shown in stdout appender output
	// also it can be used to enable/disable logger
	Name string `json:"name"`

	// minimum level of log to be shown
	Level Level `json:"-"`

	// if this flag is set to true, in case any errors in appender
	// appender should panic. This also depends on appender implementation,
	// so appender can decide to ignore or to accept information in this flag
	DoPanic bool `json:"-"`
}

// Making and sending log entry to appenders if log level is appropriate.
func (l *Logger) makeLog(msg string, lvl Level, data []interface{}) {
	if l.disabled {
		return
	}

	if lvl.value >= l.Level.value {
		log := Log{
			Time:    time.Now(),
			Message: msg,
			Level:   lvl,
			Data:    data,
			Logger:  l,
			Pid:     os.Getpid(),
		}

		for _, appender := range l.appenders {
			appender.Append(log)
		}
	}
}

// method will normalize names if they are too big or too short
// normal name length if defined by namelen variable
// if
func (l *Logger) normalizeName() {
	length := len(l.Name)

	// name is ok as it is
	if length == namelen {
		return
	}

	// name is too short, add some spaces
	if length < namelen {
		l.normalizeNameLen()
	}

	// name is too long
	// do best to normalize it

	var (
		normalized string
		parts      []string
		separator  byte
	)

	// try split long name using different separators
	// this first one which can split name into smaller parts will be used
	for _, sep := range separators {
		parts = strings.Split(l.Name, string(sep))
		if len(parts) > 1 {
			separator = sep
			break
		}
	}

	// if we sucesufully splitted string into multiple parts
	if len(parts) > 1 {
		for i, str := range parts {
			// if part length is bigger than zero
			if len(str) > 0 {
				normalized += str[:1]
				if i != (len(parts) - 1) {
					normalized += string(separator)
				}
			}
		}

		// if still to long
		if len(normalized) > namelen {
			normalized = normalized[:namelen]
		}
	} else {
		length := len(l.Name)
		if length > namelen {
			normalized = l.Name[:namelen]
		} else {
			normalized = l.Name[0 : length-1]
		}
	}

	l.Name = normalized
	if len(l.Name) < namelen {
		l.normalizeNameLen()
	}
}

// if name is still to short we will add spaces
func (l *Logger) normalizeNameLen() {
	length := len(l.Name)
	missing := namelen - length
	for i := 0; i < missing; i++ {
		l.Name += " "
	}
}

// Making log with DEBUG level.
func (l *Logger) Debug(msg string, data ...interface{}) {
	l.makeLog(msg, DEBUG, data)
}

// Making log with INFO level.
func (l *Logger) Info(msg string, data ...interface{}) {
	l.makeLog(msg, INFO, data)
}

// Making log with WARN level.
func (l *Logger) Warn(msg string, data ...interface{}) {
	l.makeLog(msg, WARN, data)
}

// Making log with ERROR level.
func (l *Logger) Error(msg string, data ...interface{}) {
	l.makeLog(msg, ERROR, data)
}

// Making log with PANIC level.
func (l *Logger) Panic(msg string, data ...interface{}) {
	l.makeLog(msg, PANIC, data)
	panic(msg)
}

// When you want to send logs to another appender,
// you should create instance of appender and call this method.
// Method is expecting appender instance to be passed
// to this method. At the end passed appender will receive logs
func (l *Logger) Enable(appender Appender) {
	l.appenders = append(l.appenders, appender)
}

// If you want to disable logs from some appender you can use this method.
// You have to call method either with appender instance,
// or you can pass appender Id as argument.
// If appender is found, it will be removed from list of appenders of this logger,
// and all other further logs won't be received by this appender.
func (l *Logger) Disable(target interface{}) {
	var id string
	var appender Appender

	switch object := target.(type) {
	case string:
		id = object
	case Appender:
		appender = object
	default:
		l.Warn("Error while disabling logger. Cannot cast to target type.")
		return
	}

	for i, app := range l.appenders {
		// if we can find the same appender reference
		// or we can extract and match id from appender
		// or we can match received id string argument with one of appender's id
		if (appender != nil && (app == appender || appender.Id() == app.Id())) || id == app.Id() {
			var toAppend []Appender

			if len(l.appenders) >= i+1 {
				toAppend = l.appenders[i+1:]
			}

			l.appenders = append(l.appenders[:i], toAppend...)
			return
		}
	}
}
