package utils

import (
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type DatabaseAccessor struct {
	*mgo.Session
	url  string
	name string
	key  int
}

func NewDatabaseAccessor(url, name string, key int) *DatabaseAccessor {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Panicf("The following error occured when accessing the database: %v", err)
	}

	return &DatabaseAccessor{session, url, name, key}
}

func (d *DatabaseAccessor) Set(r *http.Request, db *mgo.Session) {
	context.Set(r, d.key, db.DB(d.name))
}

func (d *DatabaseAccessor) Get(r *http.Request) *mgo.Database {
	if rv := context.Get(r, d.key); rv != nil {
		return rv.(*mgo.Database)
	}
	return nil
}
