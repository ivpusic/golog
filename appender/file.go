package appender

import (
	"encoding/json"
	"fmt"
	"github.com/ivpusic/golog"
	"os"
)

type Conf map[string]string

type FileAppender struct {
	path string
}

func (fa *FileAppender) Id() string {
	return "github.com/ivpusic/golog/appender/file"
}

func (fa *FileAppender) Append(log golog.Log) error {
	f, err := os.OpenFile(fa.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer f.Close()

	line, err := json.Marshal(log)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	line = append(line, byte('\n'))

	f.Write(line)
	f.Sync()

	return nil
}

func GetFileAppender(cnf Conf) *FileAppender {
	return &FileAppender{
		path: cnf["path"],
	}
}
