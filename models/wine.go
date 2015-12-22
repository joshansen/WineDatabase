package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Need to implement find by name, geographic distance, etc.
// We could also do this client side

type Wine struct {
	Id           bson.ObjectId `bson:"_id"`
	CreatedDate  time.Time
	ModifiedDate time.Time
	Name         string
	Winery       string
	Information  string
	Variety      bson.ObjectId
	Style        string
	Region       string
	Purchases    []bson.ObjectId
	Stores       []bson.ObjectId
}

func (w *Wine) AddPurchaseStore(purchaseId, storeId bson.ObjectId, db *mgo.Database) error {
	w.ModifiedDate = time.Now()

	w.Purchases = append(w.Purchases, purchaseId)
	//Appends only if store not already present
	w.Stores = appendIfMissing(w.Stores, storeId)

	return w.Save(db)
}

func appendIfMissing(slice []bson.ObjectId, i bson.ObjectId) []bson.ObjectId {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func (w *Wine) Save(db *mgo.Database) error {
	_, err := w.coll(db).UpsertId(w.Id, w)
	return err
}

func (w *Wine) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return w.coll(db).FindId(id).One(w)
}

func (*Wine) coll(db *mgo.Database) *mgo.Collection {
	return db.C("wine")
}

type Wines []Wine

func (w *Wines) FindAll(db *mgo.Database) error {
	return w.coll(db).Find(nil).All(w)
}

func (*Wines) coll(db *mgo.Database) *mgo.Collection {
	return db.C("wine")
}
