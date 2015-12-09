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
	Brand        string
	Information  string
	//ImageUrl string
	//OriginalImageUrl string
	//Types []string
	//Year int
	Bottles []Bottle
	Stores  []Stores
}

func NewWine(name, brand, information string) *Wine {
	return &Wine{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		Name:         name,
		Brand:        brand,
		Information:  information,
		//OriginalImageUrl: originalImageUrl,
	}
}

func (w *Wine) Update(name, brand, information string) {
	//need to check if variables present?
	w.ModifiedDate = time.Now()
	w.Name = name
	w.Brand = brand
	w.Information = information
	w.Save(db)
}

func (w *Wine) AddBottleStore(bottleId, storeId bson.ObjectId) {

	if len(w.Bottles) == 0 {
		w.Bottles = append([]Bottle, bottleId)
	} else {
		w.Bottles = append(w.Bottles, bottleId)
	}

	if len(w.Stores) == 0 {
		w.Stores = append([]Store, storeId)
	} else {
		w.Stores = append(w.Stores, storeId)
	}

	w.Save(db)
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

//need to implement find by name, brand, year, etc.
//could do serverside
