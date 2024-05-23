package handler

import (
	"assignment2/database"
	"assignment2/util"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// DashboardHandler is the main entry point for the Dashboards endpoint.
//
// It handles the following:
// - GET: Retrieves information about a specific registration in the database using APIs or Stub-services.
//
// If another method than GET is used, an Error occurs.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method { //a switch for the supported methods. Only GET method is supported
	case http.MethodGet:
		handleDashboardGetRequest(w, r)
	default: //Error message if GET method is not used
		http.Error(w, "This method is not supported! Only GET are supported", http.StatusMethodNotAllowed)
	}
}

// handleDashboardGetRequest retrieves a specified registration from the database,
// and finds information about it from APIs or Stub-services.
func handleDashboardGetRequest(w http.ResponseWriter, r *http.Request) {
	// Find ID in the URL.
	id, err := util.GetIdFromUrl(r.URL.Path)
	if err != nil {
		return
	}

	// Find registratration with the help of the ID:
	var reg util.Registration
	reg, err = database.GetSingleRegistrationByID(id)
	if err != nil {
		http.Error(w, "Error: could not find specified ID", http.StatusNotFound)
		return
	}

	// Find info about the country using REST Countries API or the Stub service.
	var url string
	if util.Config.Stubs.RestCountries == true {
		url = util.LOCALHOST + util.CountryStubPort + "/" // Use the Stub service.
	} else {
		url = util.COUNTRY_URL + "/alpha/" + reg.IsoCode // Use the real service.
	}
	// Create GET-request and decode response:
	var countries []util.Country
	err = util.MakeGetRequest(url, &countries)
	if err != nil {
		http.Error(w, "Error in reponse: "+err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	country := countries[0]

	// Find the coordinates of the given contry:
	latitude := country.LatitudeAndLongitude[0]
	longitude := country.LatitudeAndLongitude[1]

	// Convert the float to strings:
	lat := strconv.FormatFloat(latitude, 'f', 2, 64)
	lon := strconv.FormatFloat(longitude, 'f', 2, 64)

	// Find the Weather forecast for the given coordinates. Either with stub or real service.
	if util.Config.Stubs.Weather == true {
		url = util.LOCALHOST + util.WeatherStubPort + "/" // Use the Stub service.
	} else {
		url = util.OPEN_METO_URL + "?latitude=" + lat + "&longitude=" + lon + "&hourly=temperature_2m,precipitation" // Use the real service.
	}
	// Create GET-request and decode response:
	var weatherForcast util.Weather
	err = util.MakeGetRequest(url, &weatherForcast)
	if err != nil {
		http.Error(w, "Error in reponse: "+err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Find mean-values for temperature and precipitation
	meanTemp := findMean(weatherForcast.Hourly.Temperature)
	meanPrec := findMean(weatherForcast.Hourly.Precipitation)

	// Find the targeted currencies.
	targetedCurrencies := reg.Features.TargetCurrencies

	// Find Currency-code from the Country.
	var currencyCode string
	for key := range country.Currencies {
		if key != "" && len(key) == 3 {
			currencyCode = key
			break // IF there are more then one, just use the first one:
		}
	}

	// Get currency rates from the given country. Either with stub or real service.
	if util.Config.Stubs.Currencies == true {
		url = util.LOCALHOST + util.CurrenciesStubPort + "/" // Use the Stub service.
	} else {
		url = util.CURRENCY_URL + currencyCode // Use the real service.
	}
	// Create GET-request and decode response:
	var currency util.Currency
	err = util.MakeGetRequest(url, &currency)
	if err != nil {
		http.Error(w, "Error in reponse: "+err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Find the currency rates for the target currencies.
	targetCurrencyRate := make(map[string]float64)
	for _, val := range targetedCurrencies {
		targetCurrencyRate[val] = currency.Rates[val]
	}

	// Fix the body of the response:
	var response util.DashboardResponse

	// Fix the coordinates
	if reg.Features.Coordinates {
		response.Features.Coordinates.Latitude = latitude
		response.Features.Coordinates.Longitude = longitude
	}

	// Fix the other features:
	if reg.Features.Area {
		response.Features.Area = country.Area
	}
	if reg.Features.Capital {
		response.Features.Capital = country.CapitalCity[0]
	}
	if reg.Features.Temperature {
		response.Features.Temperature = meanTemp
	}
	if reg.Features.Precipitation {
		response.Features.Precipitation = meanPrec
	}
	if reg.Features.Population {
		response.Features.Population = country.Population
	}
	response.Features.TargetCurrencies = targetCurrencyRate

	// Fix the rest of the response struct:
	response.Name = reg.Country
	response.Isocode = reg.IsoCode
	response.LastRetrieval = time.Now().Format("2006-01-02 15:04")

	// Add the content type to the reponsewriter.
	w.Header().Add("content-type", "application/json")

	// Encode the reponse
	encoder := json.NewEncoder(w)
	err = encoder.Encode(response)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return

	}

	http.Error(w, "", http.StatusOK)

}

// findMean is a function that calculates the mean value of a list
// containing floats, and returns it as a float64
func findMean(list []float64) float64 {
	var sum float64 = 0
	for _, value := range list {
		sum += value
	}
	return sum / float64(len(list))
}
