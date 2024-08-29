package htmx

import "net/http"

func IsHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func Reswap(w http.ResponseWriter, value string) {
	w.Header().Add("HX-Reswap", value)
}

func Retarget(w http.ResponseWriter, selector string) {
	w.Header().Add("HX-Retarget", selector)
}
