package appenders

import (
	"github.com/ivpusic/golog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestMongoId(t *testing.T) {
	appender := Mongo(golog.Conf{})
	assert.Equal(t, "github.com/ivpusic/golog/appenders/mongo", appender.Id())
}

func TestMongoAppend(t *testing.T) {
	db := "test"
	coll := "logs"

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Database: db,
		Addrs:    []string{"127.0.0.1:27017"},
	})

	if err != nil {
		panic(err)
	}

	logtext := "some message"
	defer session.Close()

	c := session.DB(db).C(coll)
	// remove old logs (if any)
	c.RemoveAll(bson.M{
		"message": logtext,
	})

	// check if all old logs are removed
	log := golog.Log{}
	count, err := c.Find(bson.M{
		"message": logtext,
	}).Count()
	assert.Nil(t, err)
	assert.Exactly(t, 0, count)

	log = golog.Log{
		Message: logtext,
		// todo: add other
	}

	appender := Mongo(golog.Conf{
		"host":       "127.0.0.1:27017",
		"db":         db,
		"collection": coll,
	})

	appender.Append(log)

	// check if new log is sucesufully added
	log = golog.Log{}
	count, err = c.Find(bson.M{
		"message": logtext,
	}).Count()
	assert.Nil(t, err)
	assert.Exactly(t, 1, count)
}
