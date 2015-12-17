package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/models"
	"github.com/joshansen/WineDatabase/utils"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"time"
)

type WineControllerImpl struct {
	database utils.DatabaseAccessor
}

func NewWineController(database utils.DatabaseAccessor) *WineControllerImpl {
	return &WineControllerImpl{
		database: database,
	}
}

func (wc *WineControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/wine/{id}", wc.single)
	router.HandleFunc("/wine", wc.form).Methods("GET")
	router.HandleFunc("/wine", wc.create).Methods("POST")
}

func (wc *WineControllerImpl) single(w http.ResponseWriter, r *http.Request) {
	//written below
	data, err := wc.get(w, r)

	if err != nil {
		//TODO Fix this so it doesn't respond with only text
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		resultString, _ := json.Marshal(data)
		t, _ := template.ParseFiles("views/layout.html", "views/wine.html")
		t.Execute(w, string(resultString))
	}
}

func (wc *WineControllerImpl) get(w http.ResponseWriter, r *http.Request) (*models.Wine, error) {
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		return new(models.Wine), errors.New("Not a valid ID.")
	}
	wine := new(models.Wine)
	db := wc.database.Get(r)
	if err := wine.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		return new(models.Wine), errors.New("No such wine.")
	}

	return wine, nil
}

func (wc *WineControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/layout.html", "views/create_wine.html")

	db := wc.database.Get(r)
	varieties := new(models.Varieties)
	if err := varieties.FindAll(db); err != nil {
		//return new(models.Wines), errors.New("Could not retrieve all wines.")
	}

	if err := t.Execute(w, varieties); err != nil {
		fmt.Printf("\nThe following error occured when compiling the template: %v\n", err)
	}
}

func (wc *WineControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	wo := models.Wine{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		Name:         r.FormValue("Name"),
		Winery:       r.FormValue("Winery"),
		Information:  r.FormValue("Information"),
		Style:        r.FormValue("Style"),
		Variety:      bson.ObjectIdHex(r.FormValue("Variety")),
		Region:       r.FormValue("Region"),
	}

	db := wc.database.Get(r)
	if err := wo.Save(db); err != nil {
		fmt.Printf("Failed to save wine with error: %v\n", err)
		return
	}

	variety := new(models.Variety)
	if err := variety.FindByID(wo.Variety, db); err != nil {
		fmt.Printf("Could not find variety: %v", err)
	}
	variety.AddWine(wo.Id, db)

	utils.Redirect(w, r, "/wine/"+wo.Id.Hex())
}
