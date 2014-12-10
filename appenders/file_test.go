package appenders

import (
	"encoding/json"
	"github.com/ivpusic/golog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
}

func TestFileId(t *testing.T) {
	appender := File(golog.Conf{})
	assert.Equal(t, "github.com/ivpusic/golog/appender/file", appender.Id())
}

func TestFileAppend(t *testing.T) {
	logfile := "./log.txt"
	logtext := "some message"
	os.Remove(logfile)
	defer os.Remove(logfile)

	_, err := os.Stat(logfile)
	assert.NotNil(t, err)

	appender := File(golog.Conf{
		"path": "log.txt",
	})

	log := golog.Log{
		Message: logtext,
		// todo: add other
	}

	appender.Append(log)
	_, err = os.Stat(logfile)
	assert.Nil(t, err)

	content, err := ioutil.ReadFile(logfile)
	if err != nil {
		panic(err)
	}

	logInstance := &golog.Log{}
	err = json.Unmarshal(content, &logInstance)
	assert.Equal(t, logtext, logInstance.Message)
}
