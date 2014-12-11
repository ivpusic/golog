package appenders

import (
	"github.com/ivpusic/golog"
	"gopkg.in/mgo.v2"
)

type MongoAppender struct {
	session    *mgo.Session
	db         string
	collection string
}

// github.com/ivpusic/golog/appenders/mongo
func (ma *MongoAppender) Id() string {
	return "github.com/ivpusic/golog/appenders/mongo"
}

func (ma *MongoAppender) Append(log golog.Log) {
	ma.session.Copy()
	defer ma.session.Close()

	c := ma.session.DB(ma.db).C(ma.collection)
	c.Insert(log)
}

func Mongo(cnf golog.Conf) *MongoAppender {
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Database: cnf["db"],
		Username: cnf["username"],
		Password: cnf["password"],
		Addrs:    []string{cnf["host"]},
	})

	if err != nil {
		panic(err)
	}

	return &MongoAppender{
		db:         cnf["db"],
		collection: cnf["collection"],
		session:    sess,
	}
}
