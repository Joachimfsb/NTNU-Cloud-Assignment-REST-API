package stubs

import (
	stubHandler "assignment2/stubs/handler"
	"assignment2/util"
	"log"
	"net/http"
)

// Weather_stub starts a stub-service, listening on a specific port
// to return mocked information, which is stored in a JSON-file and is
// an old response-body from the Open Meteo API, to the client instead of
// sending a request to the actual service.
func Weather_stub() {
	weatherMux := http.NewServeMux()
	weatherMux.HandleFunc("/", stubHandler.StubWeatherHandler)

	log.Println("Weather Stub Service is listening on port: " + util.OPENWEATHER_PORT)
	log.Fatal(http.ListenAndServe(":"+util.OPENWEATHER_PORT, weatherMux))
}
