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
func (pc *PurchaseControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/purchase/{id}", pc.single)
	router.HandleFunc("/purchase", pc.form).Methods("GET")
	router.HandleFunc("/purchase", pc.create).Methods("POST")
}

//single serves the purchase page associated with the purchase id in the /purchase/{id} URL.
//This function also query's the the associated Wine, Store, and Variety records to display further information.
func (pc *PurchaseControllerImpl) single(w http.ResponseWriter, r *http.Request) {

	//Check if the provided ID is in the form of a valid ID.
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		http.Error(w, "Not a valid purchase ID.", http.StatusBadRequest)
		return
	}

	//Create empty purchase record.
	purchase := new(models.Purchase)

	//Populated the created purchase record by querying the database.
	db := pc.database.Get(r)
	if err := purchase.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		http.Error(w, "No such purchase.", http.StatusBadRequest)
		return
	}

	//Parse and execute the views/purchase.html template.
	t, err := template.ParseFiles("views/layout.html", "views/purchase.html")
	if err != nil {
		fmt.Printf("The following error occurred when parsing purchase.html: %v", err)
	}
	err = t.Execute(w, purchase)
	if err != nil {
		fmt.Printf("The following error occurred when executing purchase.html: %v", err)
	}
}

//form serves the form used to create new purchase records. This page is found at the URL /purchase/
func (pc *PurchaseControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	//Create empty wines and stores structs to hold lists of all wines and stores.
	wines := new(models.Wines)
	stores := new(models.Stores)

	//Query the database to populate the wines and stores records.
	db := pc.database.Get(r)
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
func (pc *PurchaseControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	//Max image size.
	const MAX_MEMORY = 5 * 1024 * 1024 //5MB

	//Parse the form
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		fmt.Printf("Could not parse purchase form: %v", err)
	}

	//Create a new record from the pased data.
	po := models.Purchase{
		Id:           bson.NewObjectId(),
		CreatedDate:  time.Now(),
		ModifiedDate: time.Now(),
		MemoryCue:    r.FormValue("MemoryCue"),
		Notes:        r.FormValue("Notes"),
	}

	//Get database connection
	db := pc.database.Get(r)

	//Lookup and add wine to purchase record
	wine := new(models.Wine)
	if err := wine.FindByID(bson.ObjectIdHex(r.FormValue("Wine")), db); err != nil {
		fmt.Printf("Could not find wine: %v", err)
		return
	}
	po.Wine = *wine

	//Lookup and add store to purchase record
	store := new(models.Store)
	if err := store.FindByID(bson.ObjectIdHex(r.FormValue("Store")), db); err != nil {
		fmt.Printf("Could not find store: %v", err)
		return
	}
	po.Store = *store

	//Add BuyAgain to the record.
	switch r.FormValue("BuyAgain") {
	case "true", "on", "1":
		po.BuyAgain = true
	case "false", "off", "0":
		po.BuyAgain = false
	}

	//Add Onsale to the record.
	switch r.FormValue("OnSale") {
	case "true", "on", "1":
		po.OnSale = true
	case "false", "off", "0":
		po.OnSale = false
	}

	//Convert strings so they can be added to the record.
	var err error

	//Convert price if value is present
	if r.FormValue("Price") == "" {
		po.Price = 0
	} else{
		po.Price, err = strconv.ParseFloat(r.FormValue("Price"), 64)
		if err != nil {
			fmt.Printf("The following error occured when converting the Price field to a float: %\n", err)
		}
	}

	//Convert DatePurchased if value is present
	if r.FormValue("DatePurchased") == "" {
		po.DatePurchased = *new(time.Time)
	} else{
		po.DatePurchased, err = time.Parse("2006-01-02", r.FormValue("DatePurchased"))
		if err != nil {
			fmt.Printf("The following error occured when converting the DatePurchased field to time.Time: %v\n", err)
		}
	}

	//Convert DateDrank if value is present
	if r.FormValue("DateDrank") == "" {
		po.DateDrank = *new(time.Time)
	} else{
		po.DateDrank, err = time.Parse("2006-01-02", r.FormValue("DateDrank"))
		if err != nil {
			fmt.Printf("The following error occured when converting the DateDrank field to time.Time: %v\n", err)
		}
	}

	//Convert Year if value is present
	if r.FormValue("Year") == "" {
		po.Year = 0
	} else {
		po.Year, err = strconv.Atoi(r.FormValue("Year"))
		if err != nil {
			fmt.Printf("The following error occured when converting the Year field to an integer: %v\n", err)
		}
	}

	//Convert Rating if value is present
	if r.FormValue("Rating") == "" {
		po.Rating = 0
	} else{
		po.Rating, err = strconv.Atoi(r.FormValue("Rating"))
		if err != nil {
			fmt.Printf("The following error occured when converting the Rating field to an integer: %v\n", err)
		}
	}

	//Get and upload the Image file to Amazon S3, add the name to the record.
	file, header, err := r.FormFile("Image")
	if err != nil {
		fmt.Printf("Error retreving image: %v\n", err)
	} else {
		defer file.Close()
		url := uploadToS3(file, header)
		po.ImageOriginalURL = url
	}

	//Save the record to the database.
	if err := po.Save(db); err != nil {
		fmt.Printf("Failed to save purchase with error: %v\n", err)
		fmt.Printf("Purchase record = %#v", po)
		return
	}

	//Add the purchase ID and store ID to the associated wine record.
	if err := wine.AddPurchaseStore(po.Id, po.Store.Id, db); err != nil {
		fmt.Printf("Could not add purchase and store to wine: %v", err)
		return
	}

	//Add the purchase and wine ID to the store record.
	if err := store.AddPurchaseWine(po.Id, po.Wine.Id, db); err != nil {
		fmt.Printf("Could not add wine and purchase to purchase: %v", err)
		return
	}

	//Redirect to a new URL.
	utils.Redirect(w, r, "/purchase/"+po.Id.Hex())
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
