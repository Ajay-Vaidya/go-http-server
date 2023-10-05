package util

import (
	"encoding/json"
	"net/http"
)

const (
	JSON string = "application/json"
	HTML string = "text/html"
)

// Filter slice by running the function given as the second
// argument for each element and accept ones that return true
func Filter[K comparable](data []K, f func(K) bool) []K {
	fltd := make([]K, 0)

	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}

// Return Response body as HTML or JSON based on input contentType
func ResponseBody[K any](data K, contentType string, w http.ResponseWriter, htmlRender func(data K, w http.ResponseWriter) error) {
	switch contentType {
	case JSON:
		w.Header().Set("Content-Type", JSON)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		w.Header().Set("Content-Type", HTML)
		err := htmlRender(data, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
