package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/models"
	"github.com/joshansen/WineDatabase/utils"
	"github.com/monoculum/formam"
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
	router.HandleFunc("/wine/", wc.form).Methods("GET")
	router.HandleFunc("/wine/", wc.create).Methods("POST")
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
	t.Execute(w, nil)
}

func (wc *WineControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	wo := models.Wine{Id: bson.NewObjectId(), CreatedDate: time.Now(), ModifiedDate: time.Now()}
	r.ParseForm()
	if err := formam.Decode(r.Form, &wo); err != nil {
		fmt.Printf("Failed to decode form with error: %v\n", err)
		return
	}

	wo.Save(wc.database.Get(r))
	http.Redirect(w, r, "/wine/"+wo.Id.Hex(), http.StatusSeeOther)
}
