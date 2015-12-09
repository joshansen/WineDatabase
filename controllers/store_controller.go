package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/models"
	"github.com/joshansen/WineDatabase/utils"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

type StoreControllerImpl struct {
	database utils.DatabaseAccessor
}

func NewStoreController(database utils.DatabaseAccessor) *StoreControllerImpl {
	return &StoreControllerImpl{
		database: database,
	}
}

func (sc *StoreControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/service/{id}", sc.single)
}

func (sc *StoreControllerImpl) single(w http.ResponseWriter, r *http.Request) {
	//written below
	data, err := sc.get(w, r)

	if err != nil {
		//TODO Fix this so it doesn't respond with only text
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else{
		resultString, _ := json.Marshal(data)
		t, _ := template.ParseFiles("views/layout.html", "views/store.html")
		t.Execute(w, string(resultString))
	}
}

func (sc *StoreControllerImpl) get(w http.ResponseWriter, r *http.Request) (*models.Store, error) {
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		return new(models.Store), errors.New("Not a valid ID.")
	}
	store := new(models.Store)
	db := sc.database.Get(r)
	if err := store.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		return new(models.Store), errors.New("No such store.")
	}

	return store, nil
}
