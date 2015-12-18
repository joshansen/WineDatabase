package controllers

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/models"
	"github.com/joshansen/WineDatabase/utils"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"strings"
	"time"
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
	router.HandleFunc("/store/{id}", sc.single)
	router.HandleFunc("/store", sc.form).Methods("GET")
	router.HandleFunc("/store", sc.create).Methods("POST")
}

func (sc *StoreControllerImpl) single(w http.ResponseWriter, r *http.Request) {
	//written below
	data, err := sc.get(w, r)

	if err != nil {
		//TODO Fix this so it doesn't respond with only text
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		t, _ := template.ParseFiles("views/layout.html", "views/store.html")
		t.Execute(w, data)
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

func (sc *StoreControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/layout.html", "views/create_store.html")
	t.Execute(w, nil)
}

func (sc *StoreControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	website := strings.TrimPrefix(strings.TrimPrefix(r.FormValue("Website"), "http://"), "https://")
	so := models.Store{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		Name:         r.FormValue("Name"),
		Address:      r.FormValue("Address"),
		City:         r.FormValue("City"),
		State:        r.FormValue("State"),
		Zip:          r.FormValue("Zip"),
		Website:      website,
	}

	so.Geocode()
	so.Save(sc.database.Get(r))
	utils.Redirect(w, r, "/store/"+so.Id.Hex())
}
