package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Need to implement find by name, geographic distance, etc.
// We could also do this client side

type Variety struct {
	Id           bson.ObjectId `bson:"_id"`
	CreatedDate  time.Time
	ModifiedDate time.Time
	Name         string
}

func (v *Variety) Save(db *mgo.Database) error {
	_, err := v.coll(db).UpsertId(v.Id, v)
	return err
}

func (v *Variety) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return v.coll(db).FindId(id).One(v)
}

func (*Variety) coll(db *mgo.Database) *mgo.Collection {
	return db.C("variety")
}

type Varieties []Variety

func (v *Varieties) FindAll(db *mgo.Database) error {
	return v.coll(db).Find(nil).All(v)
}

func (*Varieties) coll(db *mgo.Database) *mgo.Collection {
	return db.C("variety")
}
