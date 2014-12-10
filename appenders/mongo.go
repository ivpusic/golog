package appenders

import (
	"github.com/ivpusic/golog"
)

type MongoAppender struct {
	host       string
	db         string
	collection string
}

func (ma *MongoAppender) Id() string {
	return "github.com/ivpusic/golog/appenders/mongo"
}

func (ma *MongoAppender) Append(log golog.Log) error {
	return nil
}

func Mongo(cnf golog.Conf) *MongoAppender {
	return &MongoAppender{
		host:       cnf["host"],
		db:         cnf["db"],
		collection: cnf["collection"],
	}
}
