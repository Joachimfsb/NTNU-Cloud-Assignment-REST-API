package stubs

import (
	"assignment2/util"
	"log"
	"net/http"
)

// StubCurrencyHandler returns a mocked response from the content of the file currency.json.
// It handles the following methods:
// - GET: Returns mocked Currency information.
func StubCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Println("Received " + r.Method + " request on Currency stub handler. Returning mocked information")
		w.Header().Add("content-type", "application/json")
		response := util.ParseFile(util.STUB_CURRENCIES_RESPONSE) // Get the content of the file.
		http.Error(w, string(response), http.StatusOK)
	default:
		http.Error(w, "Method not supported!", http.StatusNotImplemented)
	}
}
