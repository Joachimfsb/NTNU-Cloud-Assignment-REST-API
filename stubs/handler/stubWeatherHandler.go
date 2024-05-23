package stubs

import (
	"assignment2/util"
	"log"
	"net/http"
)

// StubWeatherHandler returns a mocked response from the content of the file weather.json.
// It handles the following methods:
// - GET; Returns mocked Weather information.
func StubWeatherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on Weather stub handler. Returning mocked information")
		w.Header().Add("content-type", "application/json")
		response := util.ParseFile(util.STUB_WEATHER_REPONSE) // Get the content of the file.
		http.Error(w, string(response), http.StatusOK)
	default:
		http.Error(w, "Method not supported!", http.StatusNotImplemented)
	}
}
