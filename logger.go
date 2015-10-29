package golog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	DEBUG = Level{
		Value: 10,
		color: "blue",
		icon:  "★",
		Name:  "DEBUG",
	}

	INFO = Level{
		Value: 20,
		color: "green",
		icon:  "♥",
		Name:  "INFO",
	}

	WARN = Level{
		Value: 30,
		color: "yellow",
		icon:  "\u26A0",
		Name:  "WARN",
	}

	ERROR = Level{
		Value: 40,
		color: "red",
		Name:  "ERROR",
		icon:  "✖",
	}

	PANIC = Level{
		Value: 50,
		color: "black",
		icon:  "☹",
		Name:  "PANIC",
	}

	// limit when logger name will be normalized
	// normalized names are shown in console using stdout appender
	maxnamelen = 20
	curnamelen = 7

	// supported name separators
	separators []byte = []byte{'/', '.', '-'}
)

type Ctx map[string]interface{}

// Representing log level
type Level struct {
	// level priority value
	// bigger number has bigger priority
	Value int

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

	// represents data bound to contextual logger
	Ctx Ctx

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

	// represents data bound to contextual logger
	ctx Ctx
}

func (l *Logger) shouldAppend(lvl Level) bool {
	if l.disabled || lvl.Value < l.Level.Value {
		return false
	}

	return true
}

// Making and sending log entry to appenders if log level is appropriate.
func (l *Logger) makeLog(msg interface{}, lvl Level, data []interface{}) {
	log := Log{
		Time:    time.Now(),
		Message: l.toString(msg),
		Level:   lvl,
		Data:    data,
		Logger:  l,
		Pid:     os.Getpid(),
		Ctx:     l.ctx,
	}

	for _, appender := range l.appenders {
		appender.Append(log)
	}
}

func (l *Logger) toString(object interface{}) string {
	return fmt.Sprintf("%v", object)
}

// method will normalize names if they are too big or too short
// normal name length if defined by namelen variable
func (l *Logger) normalizeName() {
	length := len(l.Name)

	// name is ok as it is
	if length == maxnamelen || length == curnamelen {
		return
	}

	// name is too short, add some spaces
	if length < curnamelen {
		l.normalizeNameLen()
		return
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
		appendSeparator := true

		for i, str := range parts {
			// if part length is bigger than zero
			switch len(str) {
			case 0:
				appendSeparator = false
				break
			case 1:
				normalized += str[:1]
				break
			case 2:
				normalized += str[:2]
				break
			default:
				normalized += str[:3]
				break
			}

			if appendSeparator && (i != (len(parts) - 1)) {
				normalized += string(separator)
			}
		}

		// if still to long
		if len(normalized) > maxnamelen {
			normalized = normalized[:maxnamelen]
		}
	} else {
		length := len(l.Name)
		if length > maxnamelen {
			normalized = l.Name[:maxnamelen]
		} else {
			normalized = l.Name[0:length]
		}
	}

	l.Name = normalized
	if len(normalized) >= curnamelen {
		curnamelen = len(normalized)
	} else {
		l.normalizeNameLen()
	}
}

// if name is still to short we will add spaces
func (l *Logger) normalizeNameLen() {
	length := len(l.Name)
	missing := curnamelen - length
	for i := 0; i < missing; i++ {
		l.Name += " "
	}
}

// Making log with DEBUG level.
func (l *Logger) Debug(msg interface{}, data ...interface{}) {
	if l.shouldAppend(DEBUG) {
		l.makeLog(msg, DEBUG, data)
	}
}

// Making log with INFO level.
func (l *Logger) Info(msg interface{}, data ...interface{}) {
	if l.shouldAppend(INFO) {
		l.makeLog(msg, INFO, data)
	}
}

// Making log with WARN level.
func (l *Logger) Warn(msg interface{}, data ...interface{}) {
	if l.shouldAppend(WARN) {
		l.makeLog(msg, WARN, data)
	}
}

// Making log with ERROR level.
func (l *Logger) Error(msg interface{}, data ...interface{}) {
	if l.shouldAppend(ERROR) {
		l.makeLog(msg, ERROR, data)
	}
}

// Making log with PANIC level.
func (l *Logger) Panic(msg interface{}, data ...interface{}) {
	if l.shouldAppend(PANIC) {
		l.makeLog(msg, PANIC, data)
		panic(msg)
	}
}

// Making formatted log with DEBUG level.
func (l *Logger) Debugf(msg string, params ...interface{}) {
	if l.shouldAppend(DEBUG) {
		l.makeLog(fmt.Sprintf(msg, params...), DEBUG, nil)
	}
}

// Making formatted log with INFO level.
func (l *Logger) Infof(msg string, params ...interface{}) {
	if l.shouldAppend(INFO) {
		l.makeLog(fmt.Sprintf(msg, params...), INFO, nil)
	}
}

// Making formatted log with WARN level.
func (l *Logger) Warnf(msg string, params ...interface{}) {
	if l.shouldAppend(WARN) {
		l.makeLog(fmt.Sprintf(msg, params...), WARN, nil)
	}
}

// Making formatted log with ERROR level.
func (l *Logger) Errorf(msg string, params ...interface{}) {
	if l.shouldAppend(ERROR) {
		l.makeLog(fmt.Sprintf(msg, params...), ERROR, nil)
	}
}

// Making formatted log with PANIC level.
func (l *Logger) Panicf(msg string, params ...interface{}) {
	if l.shouldAppend(PANIC) {
		l.makeLog(fmt.Sprintf(msg, params...), PANIC, nil)
		panic(msg)
	}
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

// Will set context to current logger.
// Later appenders will be able to extract context from Log instance.
func (l *Logger) SetContext(ctx Ctx) *Logger {
	l.ctx = ctx
	return l
}

// Will copy current logger and return instance of new one.
func (l *Logger) Copy() *Logger {
	ctxLogger := &Logger{}
	*ctxLogger = *l
	return ctxLogger
}
