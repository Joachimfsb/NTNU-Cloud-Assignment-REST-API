package handler

import (
	"assignment2/database"
	"assignment2/models"
	"assignment2/util"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var startTime = time.Now() //Used for uptime (how long since service was started)

// StatusHandler is the main entry point for the status (diagnostics) endpoint.
//
// It handles the following:
// - GET: Retrieves status codes for the different services and the total number of webhooks used for the assignment.
//
// If another method than GET is used, an Error occurs.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method { //a switch for the supported methods. Only GET method is supported
	case http.MethodGet:
		handleStatusGetRequest(w, r)
	default: //Error message if GET method is not used
		http.Error(w, "This method is not supported! Only GET are supported", http.StatusMethodNotAllowed)
	}
}

// handleDashboardGetRequest retrieves status codes for the different services and the total
// number of webhooks used for the assignment.
func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var err error
	// Get status code for the REST Countires API:
	var countryStatusCode int
	if util.Config.Stubs.RestCountries == true {
		countryStatusCode, err = getStatusCode(util.LOCALHOST + util.REST_COUNTRIES_PORT + "/") // Use the stub-service
	} else {
		countryStatusCode, err = getStatusCode(util.COUNTRY_URL + "/alpha/no") // Use the real-service
	}
	if err != nil {
		http.Error(w, "Error during connection to RESTCountry API", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Get status code for the Open Meteo (Weather) API:
	var metoStatusCode int
	if util.Config.Stubs.Weather == true {
		metoStatusCode, err = getStatusCode(util.LOCALHOST + util.OPENWEATHER_PORT + "/") // Use the stub-service
	} else {
		metoStatusCode, err = getStatusCode(util.OPEN_METO_URL) // Use the real-service
	}
	if err != nil {
		http.Error(w, "Error during connection to Open Meteo API", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Get status code for the Currencies API:
	var currencyStatusCode int
	if util.Config.Stubs.Currencies == true {
		currencyStatusCode, err = getStatusCode(util.LOCALHOST + util.CURRENCIES_PORT + "/") // Use the stub-service
	} else {
		currencyStatusCode, err = getStatusCode(util.CURRENCY_URL + "/NOK") // Use the real-service
	}
	if err != nil {
		http.Error(w, "Error during connection to Currency API", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Get status code from the Notification database.
	notificationStatusCode, err := getStatusCode("http://localhost:8080/dashboard/v1/notifications/123123")
	if err != nil {
		http.Error(w, "Error during connection to Notification DB", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Get the number of webhooks:
	var notifications []models.NotificationDatabaseModel
	notifications, err = database.GetAllNotifications()
	if err != nil {
		http.Error(w, "Error retrieving webhooks", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	webhookNumber := len(notifications) // Length of the array of all webhooks is the amount of webhooks.

	// Fix the reponse
	diagnosticsMessage := util.Diagnostics{
		Countriesapi:    countryStatusCode,                    // Statuscode for Country API
		Meteoapi:        metoStatusCode,                       // Statuscode for Weather api
		Currencyapi:     currencyStatusCode,                   // Statuscode for Currencies api
		Notificationapi: notificationStatusCode,               // Statuscode for the Notification database
		NumWebhooks:     webhookNumber,                        // Number of Webhooks
		Version:         "v1",                                 //version 1
		Uptime:          int(time.Since(startTime).Seconds()), //calculated seconds since start
	}

	// Encode the response:
	encoder := json.NewEncoder(w)
	err1 := encoder.Encode(diagnosticsMessage)
	if err1 != nil {
		http.Error(w, "Error during JSON encoding.", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusOK)
}

// getStatusCode retrieves the status codes for a chosen url and returns it.
func getStatusCode(url string) (int, error) {
	// Make and issue a new GET-request
	res, err := http.Get(url)
	if err != nil {
		return res.StatusCode, err
	}
	res.Body.Close() // Close the body.

	// Everything is OK
	return res.StatusCode, nil
}
