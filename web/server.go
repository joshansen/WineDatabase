package web

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/joshansen/WineDatabase/controllers"
	"github.com/joshansen/WineDatabase/utils"
	"github.com/joshansen/WineDatabase/web/middleware"
	"github.com/unrolled/secure"
	"html/template"
	"net/http"
)

type Server struct {
	*negroni.Negroni
}

func NewServer(dba utils.DatabaseAccessor, sessionSecret string, isDevelopment bool) *Server {
	s := Server{negroni.Classic()}

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("views/layout.html", "views/index.html")
		if err != nil{
			fmt.Printf("Index template wastn't parsed with error: %v", err)
		}
		err = t.Execute(w, nil)
		if err != nil{
			fmt.Printf("Index template failed to execute with error: %v", err)
		}
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("views/layout.html", "views/404.html")
		if err != nil{
			fmt.Printf("404 template wastn't parsed with error: %v", err)
		}
		err = t.Execute(w, nil)
		if err != nil{
			fmt.Printf("404 template failed to execute with error: %v", err)
		}
	})

	storeController := controllers.NewStoreController(dba)
	storeController.Register(router)

	s.Use(negroni.HandlerFunc(secure.New(secure.Options{
		//TODO add allowed hosts
		//AllowedHosts:       []string{},
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		FrameDeny:          true,
		IsDevelopment:      isDevelopment,
	}).HandlerFuncWithNext))
	s.Use(sessions.Sessions("wineapp", cookiestore.New([]byte(sessionSecret))))
	s.Use(middleware.NewDatabase(dba).Middleware())
	s.UseHandler(router)

	return &s
}
