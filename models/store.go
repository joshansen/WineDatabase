package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Store struct {
	Id           bson.ObjectId `bson:"_id"`
	CreatedDate  time.Time
	ModifiedDate time.Time
	Name         string
	Address      string
	City         string
	State        string
	Zip          string
	Website      string
	Lattitude    float64
	Longitutde   float64
	Purchases    []bson.ObjectId
}

func (s *Store) Geocode() {
	lat, lng, err := geocode(s.Address + "," + s.City + "," + s.State + " " + s.Zip)
	if err != nil {
		fmt.Printf("Error in geocoding store address: %v", err)
	}

	s.ModifiedDate = time.Now()
	s.Lattitude = lat
	s.Longitutde = lng
}

func (s *Store) AddPurchase(purchaseId bson.ObjectId, db *mgo.Database) error {
	s.ModifiedDate = time.Now()

	s.Purchases = append(s.Purchases, purchaseId)

	return s.Save(db)
}

func (s *Store) Save(db *mgo.Database) error {
	_, err := s.coll(db).UpsertId(s.Id, s)
	return err
}

func (s *Store) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return s.coll(db).FindId(id).One(s)
}

func (*Store) coll(db *mgo.Database) *mgo.Collection {
	return db.C("store")
}

type Stores []Store

func (s *Stores) FindAll(db *mgo.Database) error {
	return s.coll(db).Find(nil).All(s)
}

func (*Stores) coll(db *mgo.Database) *mgo.Collection {
	return db.C("store")
}
