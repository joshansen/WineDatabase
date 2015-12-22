package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Purchase struct {
	Id               bson.ObjectId `bson:"_id"`
	CreatedDate      time.Time
	ModifiedDate     time.Time
	Wine             bson.ObjectId
	Store            bson.ObjectId
	Rating           int
	BuyAgain         bool
	Price            float64
	DatePurchased    time.Time
	DateDrank        time.Time
	MemoryCue        string
	Year             int
	Notes            string
	OnSale           bool
	ImageOriginalURL string
	ImageResizedURL  string
}

func (p *Purchase) Save(db *mgo.Database) error {
	_, err := p.coll(db).UpsertId(p.Id, p)
	return err
}

func (p *Purchase) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return p.coll(db).FindId(id).One(p)
}

func (*Purchase) coll(db *mgo.Database) *mgo.Collection {
	return db.C("purchase")
}
