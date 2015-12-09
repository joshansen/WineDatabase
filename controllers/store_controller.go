package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/models"
	"github.com/joshansen/WineDatabase/utils"
	"html/template"
	"net/http"
	"strings"
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

	resultString, _ := json.Marshal(data)
	t, _ := template.ParseFiles("views/layout.html", "views/store.html")
	t.Execute(w, string(resultString))
}

func (sc *StoreControllerImpl) get(w http.ResponseWriter, r *http.Request) (*models.Store, error) {
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		return storeResult{}, errors.New("Not a valid ID.")
	}
	store := new(models.Store)
	db := sc.database.Get(r)
	if err := store.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); !service.ID.Valid() || err != nil {
		return storeResult{}, errors.New("No such store.")
	}

	return store, nil
}
