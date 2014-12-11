package appenders

import (
	"encoding/json"
	"fmt"
	"github.com/ivpusic/golog"
	"os"
)

type FileAppender struct {
	path string
}

// github.com/ivpusic/golog/appender/file
func (fa *FileAppender) Id() string {
	return "github.com/ivpusic/golog/appender/file"
}

func (fa *FileAppender) Append(log golog.Log) {
	f, err := os.OpenFile(fa.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
		if log.Logger.DoPanic {
			panic(err)
		}
	}

	defer f.Close()

	line, err := json.Marshal(log)
	if err != nil {
		fmt.Println(err.Error())
		if log.Logger.DoPanic {
			panic(err)
		}
	}

	line = append(line, byte('\n'))

	f.Write(line)
	f.Sync()
}

func File(cnf golog.Conf) *FileAppender {
	return &FileAppender{
		path: cnf["path"],
	}
}
