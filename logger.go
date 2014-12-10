package golog

import (
	"os"
	"time"
)

var (
	DEBUG = Level{
		value: 10,
		color: "blue",
		Name:  "DEBUG",
	}

	INFO = Level{
		value: 20,
		color: "green",
		Name:  "INFO",
	}

	WARN = Level{
		value: 30,
		color: "yellow",
		Name:  "WARN",
	}

	ERROR = Level{
		value: 40,
		color: "red",
		Name:  "ERROR",
	}

	PANIC = Level{
		value: 50,
		color: "black",
		Name:  "PANIC",
	}
)

type Level struct {
	value int
	Name  string
	color string
}

type Log struct {
	Time    time.Time
	Message string
	Level   Level
	Data    []interface{}
	Pid     int
	Logger  *Logger
}

type Logger struct {
	appenders []Appender
	Name      string
	Level     Level `json:"-"`
}

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

func (l *Logger) Debug(msg string, data ...interface{}) {
	l.makeLog(msg, DEBUG, data)
}

func (l *Logger) Info(msg string, data ...interface{}) {
	l.makeLog(msg, INFO, data)
}

func (l *Logger) Warn(msg string, data ...interface{}) {
	l.makeLog(msg, WARN, data)
}

func (l *Logger) Error(msg string, data ...interface{}) {
	l.makeLog(msg, ERROR, data)
}

func (l *Logger) Panic(msg string, data ...interface{}) {
	l.makeLog(msg, PANIC, data)
	panic(msg)
}

func (l *Logger) Enable(appender Appender) {
	l.appenders = append(l.appenders, appender)
}

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
