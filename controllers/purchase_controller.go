package controllers

import (
	"fmt"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/models"
	"github.com/joshansen/WineDatabase/utils"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

//PurchaseControllerImpl is a stuct that holds a reference to the database to be used by the purchase controller.
type PurchaseControllerImpl struct {
	database utils.DatabaseAccessor
}

//NewPurchaseController returns a reference to a new PurchaseControllerImpl.
func NewPurchaseController(database utils.DatabaseAccessor) *PurchaseControllerImpl {
	return &PurchaseControllerImpl{
		database: database,
	}
}

//Register registers the routes controlled by the purchase controller:
// /purchase/{id}
// /purchase (Get)
// /purchase (Post)
func (bc *PurchaseControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/purchase/{id}", bc.single)
	router.HandleFunc("/purchase", bc.form).Methods("GET")
	router.HandleFunc("/purchase", bc.create).Methods("POST")
}

//single serves the purchase page associated with the purchase id in the /purchase/{id} URL.
//This function also query's the the associated Wine, Store, and Variety records to display further information.
func (bc *PurchaseControllerImpl) single(w http.ResponseWriter, r *http.Request) {

	//Check if the provided ID is in the form of a valid ID.
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		http.Error(w, "Not a valid purchase ID.", http.StatusBadRequest)
		return
	}

	//Create empty purchase, wine, store, and variety records.
	purchase := new(models.Purchase)
	wine := new(models.Wine)
	store := new(models.Store)
	variety := new(models.Variety)

	//Populated the created purchase, wine, store, and variety records by querying the database.
	db := bc.database.Get(r)
	if err := purchase.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		http.Error(w, "No such purchase.", http.StatusBadRequest)
		return
	}
	if err := wine.FindByID(purchase.Wine, db); err != nil {
		http.Error(w, "No such associated wine.", http.StatusBadRequest)
		return
	}
	if err := variety.FindByID(wine.Variety, db); err != nil {
		http.Error(w, "No such associated variety.", http.StatusBadRequest)
		return
	}
	if err := store.FindByID(purchase.Store, db); err != nil {
		http.Error(w, "No such associated store.", http.StatusBadRequest)
		return
	}

	//Create an anonymous struct that will be used to pass the populated purchase, wine, store, and variety records to the page.
	data := struct {
		Purchase *models.Purchase
		Wine     *models.Wine
		Store    *models.Store
		Variety  *models.Variety
	}{
		purchase,
		wine,
		store,
		variety,
	}

	//Parse and execute the views/purchase.html template.
	t, err := template.ParseFiles("views/layout.html", "views/purchase.html")
	if err != nil {
		fmt.Printf("The following error occurred when parsing purchase.html: %v", err)
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Printf("The following error occurred when executing purchase.html: %v", err)
	}
}

//form serves the form used to create new purchase records. This page is found at the URL /purchase/
func (bc *PurchaseControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	//Create empty wines and stores structs to hold lists of all wines and stores.
	wines := new(models.Wines)
	stores := new(models.Stores)

	//Query the database to populate the wines and stores records.
	db := bc.database.Get(r)
	if err := wines.FindAll(db); err != nil {
		http.Error(w, "Could not retrieve all wines.", http.StatusBadRequest)
		return
	}
	if err := stores.FindAll(db); err != nil {
		http.Error(w, "Could not retrieve all stores.", http.StatusBadRequest)
		return
	}

	//Create an anonymous struct that will be used to pass the populated wines and stores records to the page.
	data := struct {
		Wines  *models.Wines
		Stores *models.Stores
	}{
		wines,
		stores,
	}

	//Parse and execute the views/create_purchase.html template.
	t, err := template.ParseFiles("views/layout.html", "views/create_purchase.html")
	if err != nil {
		fmt.Printf("The following error occurred when parsing purchase.html: %v", err)
	}
	if err := t.Execute(w, data); err != nil {
		fmt.Printf("The following error occured when compiling the create_purchase.html template: %v", err)
	}
}

//create recieves and parses the form data submited to create a new purchase record, and creates the record.
func (bc *PurchaseControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	//Max image size.
	const MAX_MEMORY = 5 * 1024 * 1024 //5MB

	//Parse the form
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		fmt.Printf("Could not parse purchase form: %v", err)
	}

	//Create a new record from the pased data.
	bo := models.Purchase{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		Wine:         bson.ObjectIdHex(r.FormValue("Wine")),
		Store:        bson.ObjectIdHex(r.FormValue("Store")),
		MemoryCue:    r.FormValue("MemoryCue"),
		Notes:        r.FormValue("Notes"),
	}

	//Add BuyAgain to the record.
	switch r.FormValue("BuyAgain") {
	case "true", "on", "1":
		bo.BuyAgain = true
	case "false", "off", "0":
		bo.BuyAgain = false
	}

	//Add Onsale to the record.
	switch r.FormValue("OnSale") {
	case "true", "on", "1":
		bo.OnSale = true
	case "false", "off", "0":
		bo.OnSale = false
	}

	//Convert strings so they can be added to the record.
	var err error

	bo.Price, err = strconv.ParseFloat(r.FormValue("Price"), 64)
	if err != nil {
		fmt.Printf("The following error occured when converting the Price field to a float: %v", err)
	}
	bo.DatePurchased, err = time.Parse("2006-01-02", r.FormValue("DatePurchased"))
	if err != nil {
		fmt.Printf("The following error occured when converting the DatePurchased field to time.Time: %v", err)
	}
	bo.DateDrank, err = time.Parse("2006-01-02", r.FormValue("DateDrank"))
	if err != nil {
		fmt.Printf("The following error occured when converting the DateDrank field to time.Time: %v", err)
	}
	bo.Year, err = strconv.Atoi(r.FormValue("Year"))
	if err != nil {
		fmt.Printf("The following error occured when converting the Year field to an integer: %v", err)
	}
	bo.Rating, err = strconv.Atoi(r.FormValue("Rating"))
	if err != nil {
		fmt.Printf("The following error occured when converting the Rating field to an integer: %v", err)
	}

	//Get and upload the Image file to Amazon S3, add the name to the record.
	file, header, err := r.FormFile("Image")
	if err != nil {
		fmt.Printf("Error retreving image: %v", err)
	} else {
		defer file.Close()
		url := uploadToS3(file, header)
		bo.ImageOriginalURL = url
	}

	//Save the record to the database.
	db := bc.database.Get(r)
	if err := bo.Save(db); err != nil {
		fmt.Printf("Failed to save purchase with error: %v\n", err)
		return
	}

	//Add the purchase ID and store ID to the associated wine record.
	wine := new(models.Wine)
	if err := wine.FindByID(bo.Wine, db); err != nil {
		fmt.Printf("Could not find wine: %v", err)
		return
	}
	if err := wine.AddPurchaseStore(bo.Id, bo.Store, db); err != nil {
		fmt.Printf("Could not add purchase and store to wine: %v", err)
		return
	}

	//Add the purchase ID to the store record.
	store := new(models.Store)
	if err := store.FindByID(bo.Store, db); err != nil {
		fmt.Printf("Could not find store: %v", err)
		return
	}
	if err := store.AddPurchase(bo.Id, db); err != nil {
		fmt.Printf("Could not add wine to purchase: %v", err)
		return
	}

	//Redirect to a new URL.
	utils.Redirect(w, r, "/purchase/"+bo.Id.Hex())
}

//Upload a file to Amazon S3 and return the filename.
//See http://stackoverflow.com/questions/32152005/golang-upload-http-request-formfile-to-amazon-s3
func uploadToS3(file multipart.File, header *multipart.FileHeader) string {
	//Create the file name based on the number of nanosecods since the Unix Epoc.
	filename := strconv.FormatInt(time.Now().UnixNano(), 10)

	//Find file length by seeking to the end of the file.
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		fmt.Printf("Issue getting file size: %v", err)
		return ""
	}

	//Seek back to the begining of the file.
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("Issue seeking to begining of file: %v", err)
		return ""
	}

	//Get the Amazon S3 credentials from the enviroment.
	auth, err := aws.EnvAuth()
	if err != nil {
		fmt.Printf("Error getting AWS authentication: %v", err)
		return ""
	}

	//Set the S3 client and bucket.
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket(os.Getenv("S3_BUCKET"))

	//Upload the file to the bucket.
	err = bucket.PutReader(
		filename,
		file,
		fileSize,
		header.Header.Get("Content-Type"),
		s3.PublicRead,
		s3.Options{},
	)
	if err != nil {
		fmt.Printf("Error uploading file: %v", err)
		return ""
	}

	//Return the image name
	return filename
}
