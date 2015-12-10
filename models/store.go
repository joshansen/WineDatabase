package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"strings"
	"net/url"
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
	Website      url.URL
	Lattitude    float64
	Longitutde   float64
	Bottles      []bson.ObjectId
}

// func NewStore(name, address, city, state, zip, website string) *Store {
// 	//website = strings.TrimPrefix(strings.TrimPrefix(website, "http://"), "https://")
// 	//should this be abstracted out??
// 	lat, lng, err := geocode(address + "," + city + "," + state + " " + zip)
// 	if err != nil {
// 		fmt.Printf("Error in geocoding store address: %v", err)
// 	}
// 	return &Store{
// 		Id:           bson.NewObjectId(),
// 		CreatedDate:  time.Now(),
// 		ModifiedDate: time.Now(),
// 		Name:         name,
// 		Address:      address,
// 		City:         city,
// 		State:        state,
// 		Zip:          zip,
// 		Website:      website,
// 		Lattitude:    lat,
// 		Longitutde:   lng,
// 	}
// }

func (s *Store) Update(name, address, city, state, zip string, website url.URL, db *mgo.Database) {
	//need to check if variables present?

	//website = strings.TrimPrefix(strings.TrimPrefix(website, "http://"), "https://")
	//should this be abstracted out??
	lat, lng, err := geocode(address + "," + city + "," + state + " " + zip)
	if err != nil {
		fmt.Printf("Error in geocoding store address: %v", err)
	}

	s.ModifiedDate = time.Now()
	s.Name = name
	s.Address = address
	s.City = city
	s.State = state
	s.Zip = zip
	s.Website = website
	s.Lattitude = lat
	s.Longitutde = lng
	s.Save(db)
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

func (s *Store) AddBottle(bottleId bson.ObjectId, db *mgo.Database) {
	if len(s.Bottles) == 0 {
		s.Bottles = append(make([]bson.ObjectId, 3), bottleId)
	} else {
		s.Bottles = append(s.Bottles, bottleId)
	}

	s.Save(db)
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
