package utils

import (
	"net/http"
	"net/url"
)

func Redirect(w http.ResponseWriter, r *http.Request, defualtUrl string) {

	//queryMap := r.URL.Query()

	http.Redirect(w, r, defualtUrl, http.StatusSeeOther)

}
