package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//This type holds all purchase records. These records record the information associated with a specific purchase or bottle of wine.
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
	ImageResizedURL  string //Not currently in use
}

//Save the purchase record to the database.
func (p *Purchase) Save(db *mgo.Database) error {
	_, err := p.coll(db).UpsertId(p.Id, p)
	return err
}

//Return a purchase record given its ID.
func (p *Purchase) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return p.coll(db).FindId(id).One(p)
}

//Return the purchase collection.
func (*Purchase) coll(db *mgo.Database) *mgo.Collection {
	return db.C("purchase")
}
