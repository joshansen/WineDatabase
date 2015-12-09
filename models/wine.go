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
	Information  string
	//ImageUrl string
	//OriginalImageUrl string
	Name  string
	Brand string
	//Types []string
	//Year int
	Bottles []Bottle
	Stores  []Stores
}

func NewWine(name, information, brand string) *Wine {
	return &Wine{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		Name:         name,
		Information:  information,
		//OriginalImageUrl: originalImageUrl,
		Brand: brand,
	}
}

func (w *Wine) Update(name, information, brand string) {
	//need to check if variables present?
	w.ModifiedDate = time.Now()
	w.Name = name
	w.Information = information
	w.Brand = brand
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
