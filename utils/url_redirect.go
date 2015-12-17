package utils

import (
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request, defualtUrl string) {

	queryMap := r.URL.Query()
	prev := queryMap.Get("prev")

	if prev != "" {
		http.Redirect(w, r, "/"+prev, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, defualtUrl, http.StatusSeeOther)
	}

}
