package golog

import (
	"fmt"
	color "github.com/ivpusic/go-clicolor/clicolor"
	"io"
	"os"
	"text/tabwriter"
)

// Interface for implementing custom appenders.
type Appender interface {
	// method for injecting log to some source
	// when appender receives Log instance through this method,
	// it should decide what to do with log
	Append(log Log)

	// method will return appender ID
	// it will be used for disabling appenders
	Id() string
}

// Representing stdout appender.
type Stdout struct {
	writer     *tabwriter.Writer
	dateformat string
}

var (
	instance *Stdout
	out      io.Writer
)

// Appending logs to stdout.
func (s *Stdout) Append(log Log) {
	msg := fmt.Sprintf(" {cyan}%s {default}%s {%s}%s[%s] ▶ %s",
		log.Logger.Name,
		log.Time.Format(s.dateformat),
		log.Level.color,
		log.Level.icon,
		log.Level.Name[:4],
		log.Message)

	if out != nil {
		color.Out = out
	} else {
		color.Out = s.writer
	}

	color.Print(msg).InFormat()
	s.writer.Flush()
}

// Getting Id of stdout appender
// Id of default stdout appender is "github.com/ivpusic/golog/stdout"
func (s *Stdout) Id() string {
	return "github.com/ivpusic/golog/stdout"
}

// Function for creating and returning new stdout appender instance.
func StdoutAppender() *Stdout {
	if instance == nil {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		instance = &Stdout{
			writer:     w,
			dateformat: "15:04:05",
		}
	}

	return instance
}
