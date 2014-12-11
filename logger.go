package golog

import (
	"os"
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
)

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

type Logger struct {
	// list of appenders
	appenders []Appender
	// name of logger
	// logger name will be shown in stdout appender output
	Name string `json:"name"`
	// minimum level of log to be shown
	Level Level `json:"-"`
	// if this flag is set to true, in case any errors in appender
	// appender should panic. This also depends on appender implementation,
	// so appender can decide to ignore or to accept information in this flag
	DoPanic bool `json:"-"`
}

// making and sending log entry to appenders if log level is appropriate
func (l *Logger) makeLog(msg string, lvl Level, data []interface{}) {
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

// making log with DEBUG level
func (l *Logger) Debug(msg string, data ...interface{}) {
	l.makeLog(msg, DEBUG, data)
}

// making log with INFO level
func (l *Logger) Info(msg string, data ...interface{}) {
	l.makeLog(msg, INFO, data)
}

// making log with WARN level
func (l *Logger) Warn(msg string, data ...interface{}) {
	l.makeLog(msg, WARN, data)
}

// making log with ERROR level
func (l *Logger) Error(msg string, data ...interface{}) {
	l.makeLog(msg, ERROR, data)
}

// making log with PANIC level
func (l *Logger) Panic(msg string, data ...interface{}) {
	l.makeLog(msg, PANIC, data)
	panic(msg)
}

// when want to send logs to another appender,
// you should create instance of appender and call this method
// method is expecting appender instance to be passed
// after this method, passed appender will receive logs
func (l *Logger) Enable(appender Appender) {
	l.appenders = append(l.appenders, appender)
}

// if you want to disable logs from some appender you can use this method
// you have to call method either with appender instance
// or you can pass appender Id as argument
// if appender is found, it will be removed from list of appenders of this logger,
// and all other further logs won't be received by this appender
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
