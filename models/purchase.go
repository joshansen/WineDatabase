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
	Wine             Wine
	WineID           bson.ObjectId
	Store            Store
	StoreID          bson.ObjectId
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

//Slice of purchases
type Purchases []Purchase

//Find all purchases with the matching wineID.
func (ps *Purchases) FindByWineID(wineID bson.ObjectId, db *mgo.Database) error {
	return ps.coll(db).Find(bson.M{"wineid": wineID}).All(ps)
}

//Find all stores with the matching storeID.
func (ps *Purchases) FindByStoreID(storeID bson.ObjectId, db *mgo.Database) error {
	return ps.coll(db).Find(bson.M{"storeid": storeID}).All(ps)
}

//Find all purchases.
func (ps *Purchases) FindAll(db *mgo.Database) error {
	return ps.coll(db).Find(nil).All(ps)
}

//Return the collection of purchases.
func (*Purchases) coll(db *mgo.Database) *mgo.Collection {
	return db.C("purchase")
}
