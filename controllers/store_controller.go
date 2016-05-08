package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/models"
	"github.com/joshansen/WineDatabase/utils"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"
)

//StoreControllerImpl is a struct that holds a reference to the database to be used by the store controller.
type StoreControllerImpl struct {
	database utils.DatabaseAccessor
}

//NewStoreController returns a reference to a new StoreControllerImpl.
func NewStoreController(database utils.DatabaseAccessor) *StoreControllerImpl {
	return &StoreControllerImpl{
		database: database,
	}
}

//Register registers the routes controlled by the store controller:
// /store/{id}
// /store (Get)
// /store (Post)
func (sc *StoreControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/store/{id}", sc.single)
	router.HandleFunc("/store", sc.form).Methods("GET")
	router.HandleFunc("/store", sc.create).Methods("POST")
}

//Create three types that will be used to create the list of purchases grouped by wine.
type wineWithStats struct {
	models.Wine
	LastPurchased time.Time
	NumPurchased  int
	MinPrice      float64
}
type purchasesByWine struct {
	Wine      wineWithStats
	Purchases []models.Purchase
}
type purchasesByWines []purchasesByWine

//purchasesByWine Len method returns the length of the purchases slice.
func (p purchasesByWine) Len() int {
	return len(p.Purchases)
}

//purchasesByWine Less method compares purchase record dates.
func (p purchasesByWine) Less(i, j int) bool {
	return p.Purchases[i].DatePurchased.After(p.Purchases[j].DatePurchased)
}

//purchasesByWine Swap method swaps purchase records.
func (p purchasesByWine) Swap(i, j int) {
	p.Purchases[i], p.Purchases[j] = p.Purchases[j], p.Purchases[i]
}

//purchasesByWines Len method returns the length of the slice.
func (ps purchasesByWines) Len() int {
	return len(ps)
}

//purchasesByWines Swap method compares purchasesByWine records by wine name.
func (ps purchasesByWines) Less(i, j int) bool {
	return ps[i].Wine.Name < ps[j].Wine.Name
}

//purchasesByWines Swap method swaps purchasesByWine records.
func (ps purchasesByWines) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

//single serves the store page associated with the store id in the /store/{id} URL.
//This function also query's the the associated Wine and Purchase records to be displayed on the store page.
func (sc *StoreControllerImpl) single(w http.ResponseWriter, r *http.Request) {

	//Check if the provided ID is in the form of a valid ID.
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		http.Error(w, "Not a valid store ID.", http.StatusBadRequest)
		return
	}

	//Create empty store and purchases records.
	store := new(models.Store)
	purchases := new(models.Purchases)

	//Get a database connection.
	db := sc.database.Get(r)

	//Populate the created store and purchases records by querying the database.
	if err := store.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		http.Error(w, "No such store.", http.StatusBadRequest)
		return
	}
	if err := purchases.FindByStoreID(store.Id, db); err != nil {
		http.Error(w, "Error finding associated purchase records.", http.StatusBadRequest)
		return
	}

	//Initilize variables to be used when grouping purchases by wine.
	var wines purchasesByWines
	var match bool

	//Loop through all purchases to group them by wine.
	for _, purchase := range *purchases {
		//Initialize match to false on first pass through loop.
		match = false

		//Loop through the existing array of wines.
		for wineIndex, wine := range wines {
			//If Wine Id's Match, add purchase
			if wine.Wine.Id == purchase.Wine.Id {
				wines[wineIndex].Purchases = append(wine.Purchases, purchase)
				match = true
				break
			}
		}
		//If no match was found, append the wine and purchase.
		if !match {
			wines = append(wines, purchasesByWine{wineWithStats{purchase.Wine, time.Time{}, 0, 0.0}, []models.Purchase{purchase}})
		}
	}

	//Sort wines in alphabetical order.
	sort.Sort(wines)

	//Loop over wine purchases to calculate wine specific statistics.
	for wineIndex, wine := range wines {
		//Set NumPurchased to the length of the purchases slice.
		wines[wineIndex].Wine.NumPurchased = len(wine.Purchases)
		//Sort purchases in order of most recently purchased to least recently purchased.
		sort.Sort(wines[wineIndex])

		//Loop over all purchases for wine.
		for _, purchase := range wine.Purchases {
			//Update LastPurchased if DatePurchased is before LastPurchased.
			if wines[wineIndex].Wine.LastPurchased.Before(purchase.DatePurchased) {
				wines[wineIndex].Wine.LastPurchased = purchase.DatePurchased
			}
			//Set MinPrice to a purchased price to initialize.
			if wines[wineIndex].Wine.MinPrice == 0 {
				wines[wineIndex].Wine.MinPrice = purchase.Price
			}
			//Update MinPrice if purchase price is less than MinPrice.
			if wines[wineIndex].Wine.MinPrice > purchase.Price {
				wines[wineIndex].Wine.MinPrice = purchase.Price
			}
		}
	}

	//Create an anonymous struct that will be used to pass the populated store, purchases, and wines records to the page.
	data := struct {
		Store     *models.Store
		Purchases *models.Purchases
		Wines     purchasesByWines
	}{
		store,
		purchases,
		wines,
	}

	//Parse and execute the views/store.html template.
	t, err := template.ParseFiles("views/layout.html", "views/store.html")
	if err != nil {
		fmt.Printf("The following error occurred when parsing store.html: %v", err)
	}
	if err = t.Execute(w, data); err != nil {
		fmt.Printf("The following error occurred when executing store.html: %v", err)
	}
}

//form serves the form used to create new store records. This page is found at the URL /store/
func (sc *StoreControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	//Parse and execute the views/create_store.html template.
	t, err := template.ParseFiles("views/layout.html", "views/create_store.html")
	if err != nil {
		fmt.Printf("The following error occured when compiling create_store.html template: %v", err)
	}
	if err := t.Execute(w, nil); err != nil {
		fmt.Printf("The following error occured when executing the create_store template: %v", err)
		return
	}
}

//create receives and parses the form data submited to create a new wine record, and creates the record.
func (sc *StoreControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	//Parse the form.
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Could not parse store form: %v", err)
	}

	//Create a new record from the parsed data.
	so := models.Store{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		Name:         r.FormValue("Name"),
		Address:      r.FormValue("Address"),
		City:         r.FormValue("City"),
		State:        r.FormValue("State"),
		Zip:          r.FormValue("Zip"),
		Website:      strings.TrimPrefix(strings.TrimPrefix(r.FormValue("Website"), "http://"), "https://"),
	}

	//Geocode the record address.
	so.Geocode()

	//Save the record to the database.
	so.Save(sc.database.Get(r))

	//Redirect to a new URL.
	utils.Redirect(w, r, "/store/"+so.Id.Hex())
}
