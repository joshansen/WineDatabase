package utils

import (
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request, defualtUrl string) {
	urlMap := map[string]string{
		"home":     "",
		"variety":  "variety",
		"purchase": "purchase",
		"store":    "store",
		"wine":     "wine",
	}

	url, ok := urlMap[r.URL.Query().Get("prev")]

	if ok {
		http.Redirect(w, r, "/"+url, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, defualtUrl, http.StatusSeeOther)
	}

}
