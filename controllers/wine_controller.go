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
	"time"
)

//WineControllerImpl is a struct that holds a reference to the database to be used by the wine controller.
type WineControllerImpl struct {
	database utils.DatabaseAccessor
}

//NewWineController returns a reference to a new WineControllerImpl.
func NewWineController(database utils.DatabaseAccessor) *WineControllerImpl {
	return &WineControllerImpl{
		database: database,
	}
}

//Register registers the routes controlled by the wine controller:
// /wine/{id}
// /wine (Get)
// /wine (Post)
func (wc *WineControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/wine/{id}", wc.single)
	router.HandleFunc("/wine", wc.form).Methods("GET")
	router.HandleFunc("/wine", wc.create).Methods("POST")
}

//Define the bestYears type that will be used in single.
type bestYears []int

//Define the String method on bestYears that will be used to print a list of best years with commas.
func (ys bestYears) String() string {
	stringOfYears := ""

	for i, y := range ys {
		if i == 0 {
			stringOfYears = fmt.Sprint(y)
			continue
		}
		stringOfYears = stringOfYears + ", " + fmt.Sprint(y)
	}

	return stringOfYears
}

//Create three types that will be used to create the list of purchases grouped by store.
type storeWithStats struct {
	models.Store
	LastPurchased time.Time
	NumPurchased  int
	MinPrice      float64
}
type purchasesFromStore struct {
	Store     storeWithStats
	Purchases []models.Purchase
}
type purchasesFromStores []purchasesFromStore

func (p purchasesFromStore) Len() int {
	return len(p.Purchases)
}

func (p purchasesFromStore) Less(i, j int) bool {
	return p.Purchases[i].DatePurchased.After(p.Purchases[j].DatePurchased)
}

func (p purchasesFromStore) Swap(i, j int) {
	p.Purchases[i], p.Purchases[j] = p.Purchases[j], p.Purchases[i]
}

func (ps purchasesFromStores) Len() int {
	return len(ps)
}

func (ps purchasesFromStores) Less(i, j int) bool {
	return ps[i].Store.Name < ps[j].Store.Name
}

func (ps purchasesFromStores) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

//single serves the wine page associated with the wine ID in the /wine/{id} URL.
//This function also query's the the associated Purchases records, and calculates stats and a list of purchases by store to be displayed on the wine page.
func (wc *WineControllerImpl) single(w http.ResponseWriter, r *http.Request) {

	//Check if the provided ID is in the form of a valid ID.
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		http.Error(w, "Not a valid wine ID.", http.StatusBadRequest)
		return
	}

	//Create empty wine and purchases records.
	wine := new(models.Wine)
	purchases := new(models.Purchases)

	//Get a database connection.
	db := wc.database.Get(r)

	//Populate the created wine and purchases records by querying the database.
	if err := wine.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		http.Error(w, "No such wine.", http.StatusBadRequest)
		return
	}
	if err := purchases.FindByWineID(wine.Id, db); err != nil {
		http.Error(w, "Error finding associated purchase records.", http.StatusBadRequest)
		return
	}

	//Initilize variables to be used when grouping purchases by store.
	var stores purchasesFromStores
	var match bool

	//Loop through all purchases to group them by store.
	for _, purchase := range *purchases {
		//Initialize match to false on first pass through loop.
		match = false

		//Loop through the existing array of stores.
		for storeIndex, store := range stores {
			//If Store Id's Match, add purchase
			if store.Store.Id == purchase.Store.Id {
				stores[storeIndex].Purchases = append(store.Purchases, purchase)
				match = true
				break
			}
		}
		//If no match was found, append the store and purchase.
		if !match {
			stores = append(stores, purchasesFromStore{storeWithStats{purchase.Store, time.Time{}, 0, 0.0}, []models.Purchase{purchase}})
		}
	}

	//Sort stores in alphabetical order.
	sort.Sort(stores)

	//Loop over store purchases to calculate store specific statistics.
	for storeIndex, store := range stores {
		//Set NumPurchased to the length of the purchases slice.
		stores[storeIndex].Store.NumPurchased = len(store.Purchases)
		//Sort purchases in order of most recently purchased to least recently purchased.
		sort.Sort(stores[storeIndex])

		//Loop over all purchases for store.
		for _, purchase := range store.Purchases {
			//Update LastPurchased if DatePurchased is before LastPurchased.
			if stores[storeIndex].Store.LastPurchased.Before(purchase.DatePurchased) {
				stores[storeIndex].Store.LastPurchased = purchase.DatePurchased
			}
			//Set MinPrice to a purchased price to initialize.
			if stores[storeIndex].Store.MinPrice == 0 {
				stores[storeIndex].Store.MinPrice = purchase.Price
			}
			//Update MinPrice if purchase price is less than MinPrice.
			if stores[storeIndex].Store.MinPrice > purchase.Price {
				stores[storeIndex].Store.MinPrice = purchase.Price
			}
		}
	}

	//Create statsStruct type to hold the overall statisitcs that will be displayed.
	type statsStruct struct {
		MaxPrice           float64
		MaxPurchase        models.Purchase
		MinRegularPrice    float64
		MinRegularPurchase models.Purchase
		MinSalePrice       float64
		MinSalePurchase    models.Purchase
		AvgPrice           float64
		AvgRating          float64
		BestYears          bestYears
		LastImage          string
	}

	//Initialize the variables that will be used to calculate statistics.
	stats := new(statsStruct)
	lenPurchases := len(*purchases)
	var sumPrice float64
	var sumRating int
	var maxRating int
	var lastBought time.Time

	//Loop over all purchase records to calculate statistics.
	for _, purchase := range *purchases {
		//Update maxPrice if current maxPrice is less than purchase price.
		if stats.MaxPrice < purchase.Price {
			stats.MaxPrice = purchase.Price
			stats.MaxPurchase = purchase
		}

		//Set minRegularPrice to first nonsale purchase.
		if stats.MinRegularPrice == 0 && !purchase.OnSale {
			stats.MinRegularPrice = purchase.Price
			stats.MinRegularPurchase = purchase
		}
		//Update minRegularPrice if current minRegularPrice is greater than purchase price and the purchase wasn't on sale.
		if stats.MinRegularPrice > purchase.Price && !purchase.OnSale {
			stats.MinRegularPrice = purchase.Price
			stats.MinRegularPurchase = purchase
		}

		//Set minSalePrice to first sale purchase.
		if stats.MinSalePrice == 0 && purchase.OnSale {
			stats.MinSalePrice = purchase.Price
			stats.MinSalePurchase = purchase
		}
		//Update minSalePrice if current minSalerice is greater than purchase price and the purchase was on sale.
		if stats.MinSalePrice > purchase.Price && purchase.OnSale {
			stats.MinSalePrice = purchase.Price
			stats.MinSalePurchase = purchase
		}

		//Add current purchase price to sumPrice that will be used to calculate avgPrice.
		sumPrice = sumPrice + purchase.Price

		//Update maxRating if purchase rating is greater than current maxRating.
		if maxRating < purchase.Rating {
			maxRating = purchase.Rating
		}

		//Add current purchase rating to sumRating that will be used to calculate avgRating.
		sumRating = sumRating + purchase.Rating

		//Set LastImage if it is blank and ImageOriginalURL is not.
		if stats.LastImage == "" && purchase.ImageOriginalURL != "" {
			stats.LastImage = purchase.ImageOriginalURL
		}
		//Update lastBought if lastBought is after date purchased.
		if lastBought.Before(purchase.DatePurchased) {
			lastBought = purchase.DatePurchased
			//Update LastImage if ImageOriginalURL isn't blank.
			if purchase.ImageOriginalURL != "" {
				stats.LastImage = purchase.ImageOriginalURL
			}
		}
	}

	//Create an array of unique best years
	var maxYearMatch bool
	for _, purchase := range *purchases {
		if maxRating == purchase.Rating {
			maxYearMatch = false
			for _, year := range stats.BestYears {
				if year == purchase.Year {
					maxYearMatch = true
					break
				}
			}
			if !maxYearMatch {
				stats.BestYears = append(stats.BestYears, purchase.Year)
			}
		}
	}

	//Reverse sort the list of best years
	sort.Sort(sort.Reverse(sort.IntSlice(stats.BestYears)))

	//Calculate avgPrice and avgRating.
	stats.AvgPrice = sumPrice / float64(lenPurchases)
	stats.AvgRating = float64(sumRating) / float64(lenPurchases)

	//Create an anonymous struct that will be used to pass the populated wine, purchases, stores, and stats records to the page.
	data := struct {
		Wine      *models.Wine
		Purchases *models.Purchases
		Stores    purchasesFromStores
		Stats     *statsStruct
	}{
		wine,
		purchases,
		stores,
		stats,
	}

	//Parse and execute the views/wine.html template.
	t, err := template.ParseFiles("views/layout.html", "views/wine.html")
	if err != nil {
		fmt.Printf("The following error occured when parsing the wine.html template: %v\n", err)
	}
	//Execute the views/wine.html template.
	if err := t.Execute(w, data); err != nil {
		fmt.Printf("The following error occured when compiling the wine.html template: %v\n", err)
	}
}

//form serves the form used to create new wine records. This page is found at the URL /wine/
func (wc *WineControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	//Create empty varieties struct to hold lists of all varieties.
	varieties := new(models.Varieties)

	//Get a database connection.
	db := wc.database.Get(r)

	//Query the database to populate the varieties records.
	if err := varieties.FindAll(db); err != nil {
		fmt.Printf("The following error occured when getting all varieties: %v", err)
		return
	}

	//Parse and execute the views/create_wine.html template.
	t, err := template.ParseFiles("views/layout.html", "views/create_wine.html")
	if err != nil {
		fmt.Printf("The following error occured when compiling create_wine.html template: %v", err)
	}
	if err := t.Execute(w, varieties); err != nil {
		fmt.Printf("The following error occured when executing the create_wine template: %v", err)
		return
	}
}

//create recieves and parses the form data submited to create a new wine record, and creates the record.
func (wc *WineControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	//Parse the form.
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Could not parse wine form: %v", err)
	}

	//Create a new record from the parsed data.
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

	//Get a database connection.
	db := wc.database.Get(r)

	//Lookup and add variety to wine record.
	variety := new(models.Variety)
	if err := variety.FindByID(bson.ObjectIdHex(r.FormValue("Variety")), db); err != nil {
		fmt.Printf("Could not find variety: %v", err)
	}
	wo.Variety = *variety
	wo.VarietyID = wo.Variety.Id

	//Save the record to the database.
	if err := wo.Save(db); err != nil {
		fmt.Printf("Failed to save wine with error: %v\n", err)
		fmt.Printf("Wine record = %#v\n", wo)
		return
	}

	//Redirect to a new URL.
	utils.Redirect(w, r, "/wine/"+wo.Id.Hex())
}
