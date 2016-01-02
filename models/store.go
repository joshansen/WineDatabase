package models

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
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
	Website      string
	Lattitude    float64
	Longitutde   float64
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

func (s *Store) Geocode() {
	lat, lng, err := geocode(s.Address + "," + s.City + "," + s.State + " " + s.Zip)
	if err != nil {
		fmt.Printf("Error in geocoding store address: %v", err)
	}

	s.ModifiedDate = time.Now()
	s.Lattitude = lat
	s.Longitutde = lng
}

func geocode(address string) (lat float64, lng float64, err error) {
	const (
		geocodeURL = "https://maps.googleapis.com/maps/api/geocode/json?address="
	)

	type geocodingResults struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	//can add api key with the following  + "&key=" + apiKey
	resp, err := http.Get(geocodeURL + url.QueryEscape(address))

	if err != nil {
		return 0, 0, fmt.Errorf("Error geocoding address: <%v>", err)
	}

	defer resp.Body.Close()

	var result geocodingResults

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)

	if err != nil {
		return 0, 0, fmt.Errorf("Error decoding geocoding result: <%v>", err)
	}

	if len(result.Results) > 0 {
		lat = result.Results[0].Geometry.Location.Lat
		lng = result.Results[0].Geometry.Location.Lng
	}

	return lat, lng, nil
}
