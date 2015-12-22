package controllers

import (
	//"errors"
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

type PurchaseControllerImpl struct {
	database utils.DatabaseAccessor
}

func NewPurchaseController(database utils.DatabaseAccessor) *PurchaseControllerImpl {
	return &PurchaseControllerImpl{
		database: database,
	}
}

//create a function to register the purchase urls
func (bc *PurchaseControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/purchase/{id}", bc.single)
	router.HandleFunc("/purchase", bc.form).Methods("GET")
	router.HandleFunc("/purchase", bc.create).Methods("POST")
}

//servethe purchase.html page
func (bc *PurchaseControllerImpl) single(w http.ResponseWriter, r *http.Request) {

	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		//errors.New("Not a valid ID.")
	}

	purchase := new(models.Purchase)
	wine := new(models.Wine)
	store := new(models.Store)
	variety := new(models.Variety)

	db := bc.database.Get(r)
	if err := purchase.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		//return new(models.Purchase), errors.New("No such purchase.")
	}
	if err := wine.FindByID(purchase.Wine, db); err != nil {
		//return new(models.Purchase), errors.New("No such purchase.")
	}
	if err := variety.FindByID(wine.Variety, db); err != nil {
		//return new(models.Purchase), errors.New("No such purchase.")
	}
	if err := store.FindByID(purchase.Store, db); err != nil {
		//return new(models.Purchase), errors.New("No such purchase.")
	}

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

	// if err != nil {
	// 	//TODO Fix this so it doesn't respond with only text
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// } else {
	t, _ := template.ParseFiles("views/layout.html", "views/purchase.html")
	t.Execute(w, data)
	//}
}

//serve the create_purchase page
func (bc *PurchaseControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/layout.html", "views/create_purchase.html")

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

	if err := t.Execute(w, WinesStores{Wines: wines, Stores: stores}); err != nil {
		fmt.Printf("\nThe following error occured when compiling the create_purchase.html template: %v\n", err)
	}
}

//create a new purchase record
func (bc *PurchaseControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	//parse form with max file size
	const MAX_MEMORY = 5 * 1024 * 1024 //5MB

	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		fmt.Printf("Could not parse purchase form: %v", err)
	}

	bo := models.Purchase{
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

	file, header, err := r.FormFile("Image")

	if err != nil {
		fmt.Printf("Error retreving image: %v", err)
	} else {
		defer file.Close()
		url := uploadToS3(file, header)
		bo.ImageOriginalURL = url
	}

	//get the database
	db := bc.database.Get(r)

	//save the purchase
	if err := bo.Save(db); err != nil {
		fmt.Printf("Failed to save purchase with error: %v\n", err)
		return
	}

	//add purchase and store to the wine
	wine := new(models.Wine)
	if err := wine.FindByID(bo.Wine, db); err != nil {
		fmt.Printf("Could not find wine: %v", err)
		return
	}
	if err := wine.AddPurchaseStore(bo.Id, bo.Store, db); err != nil {
		fmt.Printf("Could not add purchase and store to wine: %v", err)
		return
	}

	//add purchase to store
	store := new(models.Store)
	if err := store.FindByID(bo.Store, db); err != nil {
		fmt.Printf("Could not find store: %v", err)
		return
	}
	if err := store.AddPurchase(bo.Id, db); err != nil {
		fmt.Printf("Could not add wine to purchase: %v", err)
		return
	}

	//redirect to a new url
	utils.Redirect(w, r, "/purchase/"+bo.Id.Hex())
}

//See http://stackoverflow.com/questions/32152005/golang-upload-http-request-formfile-to-amazon-s3
func uploadToS3(file multipart.File, header *multipart.FileHeader) (url string) {
	filename := strconv.FormatInt(time.Now().UnixNano(), 10)

	//find file length
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		fmt.Printf("Issue getting file size: %v", err)
		return
	}

	//return to begining of file
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("Issue seeking to begining of file: %v", err)
		return
	}

	//get credentials from the enviroment
	auth, err := aws.EnvAuth()
	if err != nil {
		fmt.Printf("Error getting AWS authentication: %v", err)
		return
	}

	//set client and bucket
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket(os.Getenv("S3_BUCKET"))

	//upload file
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
		return
	}

	//return the image url
	url = filename
	return
}
