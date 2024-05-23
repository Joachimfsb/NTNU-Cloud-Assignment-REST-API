package stubs

import (
	stubHandler "assignment2/stubs/handler"
	"assignment2/util"
	"log"
	"net/http"
)

// Currency_stub starts a stub-service, listening on a specific port
// to return mocked information, which is stored in a JSON-file and is
// an old response-body from the Currency API, to the client instead of
// sending a request to the actual service.
func Currency_stub() {
	currencyMux := http.NewServeMux()
	currencyMux.HandleFunc("/", stubHandler.StubCurrencyHandler)

	log.Println("Currency Stub Service is listening on port: " + util.CURRENCIES_PORT)
	log.Fatal(http.ListenAndServe(":"+util.CURRENCIES_PORT, currencyMux))
}
