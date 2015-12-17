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

type VarietyControllerImpl struct {
	database utils.DatabaseAccessor
}

func NewVarietyController(database utils.DatabaseAccessor) *VarietyControllerImpl {
	return &VarietyControllerImpl{
		database: database,
	}
}

func (vc *VarietyControllerImpl) Register(router *mux.Router) {
	router.HandleFunc("/variety/{id}", vc.single)
	router.HandleFunc("/variety/", vc.form).Methods("GET")
	router.HandleFunc("/variety/", vc.create).Methods("POST")
}

func (vc *VarietyControllerImpl) single(w http.ResponseWriter, r *http.Request) {
	//written below
	data, err := vc.get(w, r)

	if err != nil {
		//TODO Fix this so it doesn't respond with only text
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		resultString, _ := json.Marshal(data)
		t, _ := template.ParseFiles("views/layout.html", "views/variety.html")
		t.Execute(w, string(resultString))
	}
}

func (vc *VarietyControllerImpl) get(w http.ResponseWriter, r *http.Request) (*models.Variety, error) {
	if !bson.IsObjectIdHex(mux.Vars(r)["id"]) {
		return new(models.Variety), errors.New("Not a valid ID.")
	}
	variety := new(models.Variety)
	db := vc.database.Get(r)
	if err := variety.FindByID(bson.ObjectIdHex(mux.Vars(r)["id"]), db); err != nil {
		return new(models.Variety), errors.New("No such variety.")
	}

	return variety, nil
}

func (vc *VarietyControllerImpl) form(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/layout.html", "views/create_variety.html")
	t.Execute(w, nil)
}

func (vc *VarietyControllerImpl) create(w http.ResponseWriter, r *http.Request) {
	vo := models.Variety{Id: bson.NewObjectId(), CreatedDate: time.Now(), ModifiedDate: time.Now()}
	r.ParseForm()
	if err := formam.Decode(r.Form, &vo); err != nil {
		fmt.Printf("Failed to decode form with error: %v\n", err)
		return
	}

	vo.Save(vc.database.Get(r))
	utils.Redirect(w, r, "/variety/"+vo.Id.Hex())
}
