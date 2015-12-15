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
	"strconv"
	"time"
)

type BottleControllerImpl struct {
	database utils.DatabaseAccessor
}

func NewBottleController(database utils.DatabaseAccessor) *BottleControllerImpl {
	return &BottleControllerImpl{
		database: database,
	}
}

func (bc *BottleControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/bottle/{id}", bc.single)
	router.HandleFunc("/bottle/", bc.form).Methods("GET")
	router.HandleFunc("/bottle/", bc.create).Methods("POST")
}

func (bc *BottleControllerImpl) single(w http.ResponseWriter, r *http.Request) {
	//written below
	data, err := bc.get(w, r)

	if err != nil {
		//TODO Fix this so it doesn't respond with only text
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		resultString, _ := json.Marshal(data)
		t, _ := template.ParseFiles("views/layout.html", "views/bottle.html")
		t.Execute(w, string(resultString))
	}
}

func (bc *BottleControllerImpl) get(w http.ResponseWriter, r *http.Request) (*models.Bottle, error) {
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		return new(models.Bottle), errors.New("Not a valid ID.")
	}
	bottle := new(models.Bottle)
	db := bc.database.Get(r)
	if err := bottle.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		return new(models.Bottle), errors.New("No such bottle.")
	}

	return bottle, nil
}

func (bc *BottleControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/layout.html", "views/create_bottle.html")

	db := bc.database.Get(r)
	wines := new(models.Wines)
	if err := wines.FindAll(db); err != nil {
		//return new(models.Wines), errors.New("Could not retrieve all wines.")
	}
	stores := new(models.Stores)
	if err := stores.FindAll(db); err != nil {
		//return new(models.Stores), errors.New("Could not retrieve all stores.")
	}

	type WinesStores struct {
		Wines  *models.Wines
		Stores *models.Stores
	}

	//resultString, _ := json.Marshal(WinesStores{Wines: wines, Stores: stores})

	if err := t.Execute(w, WinesStores{Wines: wines, Stores: stores}); err != nil {
		fmt.Printf("\nThe following error occured when compiling the template: %v\n", err)
	}
}

func (bc *BottleControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	bo := models.Bottle{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		Wine:         bson.ObjectIdHex(r.FormValue("Wine")),
		Store:        bson.ObjectIdHex(r.FormValue("Store")),
		MemoryCue:    r.FormValue("MemoryCue"),
		Notes:        r.FormValue("Notes"),
	}

	switch r.FormValue("BuyAgain") {
	case "true", "on", "1":
		bo.BuyAgain = true
	case "false", "off", "0":
		bo.BuyAgain = false
	}

	switch r.FormValue("OnSale") {
	case "true", "on", "1":
		bo.OnSale = true
	case "false", "off", "0":
		bo.OnSale = false
	}

	bo.Price, _ = strconv.ParseFloat(r.FormValue("Price"), 64)
	bo.DatePurchased, _ = time.Parse("2006-01-02", r.FormValue("DatePurchased"))
	bo.DateDrank, _ = time.Parse("2006-01-02", r.FormValue("DateDrank"))
	bo.Year, _ = strconv.Atoi(r.FormValue("Year"))
	bo.Rating, _ = strconv.Atoi(r.FormValue("Rating"))

	db := bc.database.Get(r)

	if err := bo.Save(db); err != nil {
		fmt.Printf("Failed to save bottle with error: %v\n", err)
		return
	}

	wine := new(models.Wine)
	if err := wine.FindByID(bo.Wine, db); err != nil {
		fmt.Printf("Could not find wine: %v", err)
	}
	wine.AddBottleStore(bo.Id, bo.Store, db)

	store := new(models.Store)
	if err := store.FindByID(bo.Store, db); err != nil {
		fmt.Printf("Could not find store: %v", err)
	}
	store.AddBottle(bo.Id, db)

	http.Redirect(w, r, "/bottle/"+bo.Id.Hex(), http.StatusSeeOther)
}
