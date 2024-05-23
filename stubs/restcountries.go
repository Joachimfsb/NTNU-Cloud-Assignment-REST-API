package stubs

import (
	stubHandler "assignment2/stubs/handler"
	"assignment2/util"
	"log"
	"net/http"
)

// Country_stub starts a stub-service, listening on a specific port
// to return mocked information, which is stored in a JSON-file and is
// an old response-body from the RESTCountries API, to the client instead of
// sending a request to the actual service.
func Country_stub() {
	countryMux := http.NewServeMux()
	countryMux.HandleFunc("/", stubHandler.StubCountryHandler)

	log.Println("Country Stub Service is listening on port: " + util.REST_COUNTRIES_PORT)
	log.Fatal(http.ListenAndServe(":"+util.REST_COUNTRIES_PORT, countryMux))
}
