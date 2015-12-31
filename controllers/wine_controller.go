package controllers

import (
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
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		//return new(models.Wine), errors.New("Not a valid ID.")
	}

  wine := new(models.Wine)
	// variety := new(models.Variety)
	// purchases := new(models.Purchases)

	db := wc.database.Get(r)
	if err := wine.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		//return new(models.Wine), errors.New("No such wine.")
	}
	// if err := variety.FindByID(wine.Variety, db); err != nil {
	// 	//return new(models.Wine), errors.New("No such wine.")
	// }
	// if err := purchases.FindByWineID(wine.Id, db); err != nil {
	// 	//return new(models.Wine), errors.New("No such wine.")
	// }

	// type PurchaseWithStore struct{
	// 	*models.Purchase
	// 	*models.Store
	// }

	// purchasesWithStores := make([]PurchaseWithStore, len(purchases))

	// store := new(models.Store)
	// for _, purchase := range purchases {
	// 	store.FindById(purchase.Store, db)
	// 	purchasesWithStores = append(purchasesWithStores, PurchaseWithStore{purchase, store})

	// }

	// data := struct {
	// 	Wine      *models.Wine
	// 	Variety   *models.Variety
	// 	Purchases *models.Purchases
	// 	PurchasesWithStores PurchaseWithStore
	// }{
	// 	wine,
	// 	variety,
	// 	purchases,
	// 	purchasesWithStores,
	// }

	data := wine

	t, err := template.ParseFiles("views/layout.html", "views/wine.html")
	if err != nil {
		fmt.Printf("The following error occured when compiling wine.html template: %v", err)
	}

	if err := t.Execute(w, data); err != nil {
		fmt.Printf("The following error occured when compiling the create_wine template: %v", err)
	}
}

func (wc *WineControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/layout.html", "views/create_wine.html")
	if err != nil {
		fmt.Printf("The following error occured when compiling create_wine.html template: %v", err)
	}

	db := wc.database.Get(r)
	varieties := new(models.Varieties)
	if err := varieties.FindAll(db); err != nil {
		fmt.Printf("The following error occured when getting all varieties: %v", err)
		return
	}

	if err := t.Execute(w, varieties); err != nil {
		fmt.Printf("The following error occured when executing the create_wine template: %v", err)
		return
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
		Region:       r.FormValue("Region"),
	}

	db := wc.database.Get(r)
	variety := new(models.Variety)
	if err := variety.FindByID(bson.ObjectIdHex(r.FormValue("Variety")), db); err != nil {
		fmt.Printf("Could not find variety: %v", err)
	}

	wo.Variety = *variety

	if err := wo.Save(db); err != nil {
		fmt.Printf("Failed to save wine with error: %v\n", err)
		return
	}

	variety.AddWine(wo.Variety.Id, db)

	utils.Redirect(w, r, "/wine/"+wo.Id.Hex())
}
