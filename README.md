golog
=====
[![Build Status](https://travis-ci.org/ivpusic/golog.svg?branch=master)](https://travis-ci.org/ivpusic/golog)

Simple but powerful logging library

![alt text](http://s2.postimg.org/gocipacjt/Screenshot_from_2014_12_11_13_11_11.png "")

### Example
```Go
package main

import "github.com/ivpusic/golog"

func main() {
	// get default logger
	logger := golog.Default
	
	// default level for all loggers is DEBUG
	// you can easily change it it you want
	logger.Level = golog.WARN

	// log something
	logger.Debug("some message")
}
```

### Features
- Multiple loggers
- Stdout appender
- File appender
- Mongo appender
- Simple API for writing custom appenders
- Enabling/disabling appenders

### Installation
```Shell
go get github.com/ivpusic/golog
```

### Levels
Currently supported levels are
- DEBUG
- INFO
- WARN
- ERROR
- PANIC (in this case program will panic)

### Multiple loggers
You can ask ``golog`` for logger instance. Logger instances are singletones.
```Go
package main

import "github.com/ivpusic/golog"

func main() {
	// default logger
	logger := golog.Default
	logger.Debug("some message")
	
	// another logger
	// in this case we are issuing logger with name application
	// but it could be something else also
	application := golog.GetLogger("application")
	application.Debug("some message")
}
```

### Appenders
Golog provided set of default appenders and provides ultra simple API for adding new ones.

#### Enabling appenders
Stdout appender is enabled by default, and you can enable additional appenders using ``Enable`` method of appender.

##### File
```Go
package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

func main() {
	logger := golog.Default

	// make instance of file appender and enable it
	logger.Enable(appenders.File(golog.Conf{
		"path": "/path/to/log.txt",
	}))

	logger.Debug("some message")
}
```

##### Mongo
```Go
package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

func main() {
	logger := golog.Default

	// make instance of mongo appender and enable it
	logger.Enable(appenders.Mongo(golog.Conf{
		"host":       "127.0.0.1:27017",
		"db":         "somedb",
		"collection": "logs",
		"username":   "myusername",
		"password":   "mypassword",
	}))

	logger.Debug("some message")
}
```

#### Disabling appenders
You can disable appender by calling ``Disable`` method on appender.

```Go
package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

func main() {
	logger := golog.Default

	appender := appenders.File(golog.Conf{
		"path": "/path/to/log.txt",
	})
	
	// let we say that we first enabled appender
	logger.Enable(appender)

	// and at the some point we want to disable it
	logger.Disable(appender)
	
	// you can also disable appender by passing appender id
	// in this case we are disabling file appender, so we will pass it's id
	logger.Disable("github.com/ivpusic/golog/appender/file")

  // this log won't go to disable appender
	logger.Debug("some message")
}
```

#### Custom appenders
Writing your own appender to really simple. You have to implement ``Id`` and ``Append`` methods. In ``Append`` method you will receive log instance, and you can do with it wathever you want.

```Go
package main

import "github.com/ivpusic/golog"

type CustomAppender struct {
}

func (s *CustomAppender) Append(log golog.Log) {
	// do something with log
	// for example you can save it to database
	// send it to some service, etc.
}

func (s *CustomAppender) Id() string {
	return "id/of/custom/appender"
}

func main() {
	logger := golog.Default

	// make custom appender instance
	appender := &CustomAppender{}

	// enable custom appender
	logger.Enable(appender)

	// log something
	logger.Debug("this will go to custom appender also")
}
```

### GoDoc
For additional documentation and detailed info about package structures please visit  [this](https://godoc.org/github.com/ivpusic/golog) link.

# License
MIT
