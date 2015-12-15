package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Bottle struct {
	Id           bson.ObjectId `bson:"_id"`
	CreatedDate  time.Time
	ModifiedDate time.Time
	Wine         bson.ObjectId
	Store        bson.ObjectId
	Rating       int
	BuyAgain     bool
	//Is this the proper type
	Price         float64
	DatePurchased time.Time
	DateDrank     time.Time
	MemoryCue     string
	Year          int
	Notes         string
	OnSale        bool
}

// func (b *Bottle) Update(wine bson.ObjectId, store bson.ObjectId, notes string, memoryCue string, buyAgain bool, doWeLike bool, price float64, datePurchased time.Time, dateDrank time.Time, year int, db *mgo.Database) {
// 	//need to check if variables present?
// 	b.ModifiedDate = time.Now()
// 	b.Notes = notes
// 	b.MemoryCue = memoryCue
// 	b.Wine = wine
// 	b.Store = store
// 	b.BuyAgain = buyAgain
// 	b.DoWeLike = doWeLike
// 	b.Price = price
// 	b.DatePurchased = datePurchased
// 	b.DateDrank = dateDrank
// 	b.Year = year
// 	b.Save(db)
// }

func (b *Bottle) Save(db *mgo.Database) error {
	_, err := b.coll(db).UpsertId(b.Id, b)
	return err
}

func (b *Bottle) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return b.coll(db).FindId(id).One(b)
}

func (*Bottle) coll(db *mgo.Database) *mgo.Collection {
	return db.C("bottle")
}
