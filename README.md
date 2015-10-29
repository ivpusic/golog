golog
=====
[![Build Status](https://travis-ci.org/ivpusic/golog.svg?branch=master)](https://travis-ci.org/ivpusic/golog)

Simple but powerful go logging library

![alt text](http://s2.postimg.org/gocipacjt/Screenshot_from_2014_12_11_13_11_11.png "")

### Example
```Go
package main

import "github.com/ivpusic/golog"

func main() {
	// get default logger
	logger := golog.Default
	
	// default level for all loggers is DEBUG
	// you can easily change it if you want
	logger.Level = golog.WARN

	// log something
	logger.Debug("some message")
}
```

### Features
- Multiple loggers
- Logger Context
- Copying Logger
- Appenders
	- Stdout appender
	- File appender
	- Mongo appender
- Simple API for writing custom appenders
- Enabling/disabling appenders
- Enabling/disabling loggers
- Attaching log data
- Formatting logs

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

### Formatting
Normally you call one of ``Debug``, ``Info``, etc.. methods of logger when you want to log some string. But sometimes you want to format your log, so you want to pass format and parameters related to format. Let's see example:
```Go
package main

import "github.com/ivpusic/golog"

func main() {
	logger := golog.Default

	// will output `some cool number 4`
	// the same you can do for other levels
	logger.Debugf("some %s number %d", "cool", 4)
}
```

### Logger Context
Every logger can have it's own context. Later every appender will be able to extract context which is shared between all logs sent by that logger.

```Go
package main

import "github.com/ivpusic/golog"

func main() {
  logger := golog.Default

  logger.SetContext(golog.Ctx{
    "field1": 123,
    "field2": "value"
  })

  // Calling copy will make new instance of logger, 
  // but with all values copied from original logger.
  // After copy loggers are independent instances.
  logger = logger.Copy().SetContext(golog.Ctx{
    "field1": 123,
    "field2": "value"
  })
  
  logger.Debug("message")
}
```

### Attaching data
You can attach data to log. Be aware that your appender have to support this. Appender will be able to access passed data using ``Data`` member of ``golog.Log`` type. First argument of logging method is string which is actual log message, and other parameters are data which can be optionally attached to log.
```Go
package main

import "github.com/ivpusic/golog"

type SomeType struct {
	Something 		string
	SomethingElse 	int
}

func main() {
	logger := golog.Default

	data1 := SomeType{"blabla", 10}
	data2 := SomeType{"blabla", 10}
	// first parameter is string, and second is ...inteface{}
	// the same you can do for other level
	// once again, appender which you are using should support (save) passed data
	logger.Debug("some log message", data1, data2)
}
```

### Multiple loggers
You can ask ``golog`` for logger instance. Logger instances are singletons.
```Go
package main

import "github.com/ivpusic/golog"

func main() {	
	// get logger with name github.com/someuser/somelib
	// if logged doesn't exists, it will be created
	logger := golog.GetLogger("github.com/someuser/somelib")

	// this can be very useful if library which you are using uses
	// golog too. Then you can control logger level or logger appenders
	logger.Level = golog.DEBUG

	// or you can just log something using that logger
	logger.Debug("some message")
}
```

#### Enabling/Disabling loggers
If library which you are using uses ``golog`` you can explicitly enable or disable logger. This gives you control over logs which you want to see (in your console for example), and the ones which you don't.
```Go
package main

import "github.com/ivpusic/golog"

func main() {
	// you have to provide logger name in order to disable it
	golog.Disable("github.com/someuser/somelib")

	// all loggers are enabled by default
	// if you have case that at some point you disable it,
	// and later you want to enable it again, you can use this method
	golog.Enable("github.com/someuser/somelib")
}
```

### Appenders
Golog provides set of default appenders and ultra simple API for adding new ones.

##### Stdout
This appender is enabled by default, so by default you should see messages in your terminal with following structure:
```
{logger_name} {date} {level} {message}
```

#### Enabling appenders
As you know stdout appender is enabled by default. You can enable additional appenders using ``Enable`` method of logger.

##### File
```Go
package main

import "github.com/ivpusic/golog"
import "github.com/ivpusic/golog/appenders"

func main() {
	logger := golog.Default

	// make instance of file appender and enable it
	logger.Enable(appenders.File(golog.Conf{
		// file in which logs will be saved
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
		// host and database port of target mongo database
		// where logs will be saved
		"host":       "127.0.0.1:27017",
		// target database in which logs will be saved
		"db":         "somedb",
		// target collection in which logs will be saved
		"collection": "logs",
		// database username (if exists)
		"username":   "myusername",
		// database password (if exists)
		"password":   "mypassword",
	}))

	logger.Debug("some message")
}
```

#### Disabling appenders
You can disable appender by calling ``Disable`` method of logger.

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

  	// this log won't go to disabled appender
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

### Conventions
We should name propperly our loggers and appenders if we want that others don't have troubles when they want to use them.

#### Logger names
Logger name should be ``go get`` path of your library. So for example if you have library hosted on github, and users are getting library with ``go get github.com/someuser/somelibrary``, then you shold use logger:
```Go
logger := golog.GetLogger("github.com/someuser/somelibrary")
```
On this way users will easily get logger of your library, change its level, enable it or disable it.
#### Appender names
Let we say that you hosted your appender on github, and whole repo is reserved only for that appender, so users will get your appender with ``go get github.com/someuser/someappender``. In this case ID of appender should be ``github.com/someuser/someappender``.

In case that you have one repo and in that repo you have multiple appenders available, then you append appender name to ``go get`` path. So if your repository is ``github.com/someuser/appenders``, and you have appender A and appender B available in that repo, then ID of appender A should be ``github.com/someuser/appenders/A`` and ID of appender B should be ``github.com/someuser/appenders/B``.

### GoDoc
For additional documentation and detailed info about package structures please visit  [this](https://godoc.org/github.com/ivpusic/golog) link.

# License
MIT
