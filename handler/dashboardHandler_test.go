package handler_test

import (
	"assignment2/handler"
	stubs "assignment2/stubs/handler"
	"assignment2/util"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// populateDashboardsFile sets an example registration in the
// registrations.json-file (In stubs/res), to use during testing.
func populateDashboardsFile() error {
	// Register the registration to test with
	var allDashboards []util.Registration

	dashboard := util.Registration{
		ID:      "1",
		Country: "Norway",
		IsoCode: "NO",
		Features: util.Features{
			Temperature:      true,
			Precipitation:    true,
			Capital:          true,
			Coordinates:      true,
			Population:       true,
			Area:             true,
			TargetCurrencies: []string{"NOK", "EUR"},
		},
	}
	allDashboards = append(allDashboards, dashboard)

	// Write it to the JSON-file
	return util.PopulateTestFile(util.STUB_DATABASE_REGISTRATIONS, allDashboards)
}

// TestRetrievePopulatedDashboard tests if the DashboardHandler only allows
// methods of type GET, and tests if the data returned is expected based on
// an example registration and mocked data from the stub-services.
//
// Requirements:
// - Only GET requests allowed
// - The data returned from the DashboardHandler corresponds to the expected data.
func TestRetrievePopulatedDashboard(t *testing.T) {
	util.FixStubPaths()

	// Declare the stub services:
	util.Config.Stubs.Database = true
	util.Config.Stubs.Weather = true
	util.Config.Stubs.Currencies = true
	util.Config.Stubs.RestCountries = true

	// Start the test Weather Stub-service
	stubWeather := httptest.NewServer(http.HandlerFunc(stubs.StubWeatherHandler))
	defer stubWeather.Close()

	// Find the port of said Weather Stub-service
	port := strings.Split(stubWeather.URL, ":")
	portNum := port[len(port)-1]
	util.WeatherStubPort = portNum // Change the global variable to the new port.

	// Start the test Country Stub-service
	stubCountry := httptest.NewServer(http.HandlerFunc(stubs.StubCountryHandler))
	defer stubCountry.Close()

	// Find the port of said Country Stub-service
	port = strings.Split(stubCountry.URL, ":")
	portNum = port[len(port)-1]
	util.CountryStubPort = portNum // Change the global variable to the new port.

	// Start the test Currency Stub-service
	stubCurrency := httptest.NewServer(http.HandlerFunc(stubs.StubCurrencyHandler))
	defer stubCurrency.Close()

	// Find the port of said Currency Stub-service
	port = strings.Split(stubCurrency.URL, ":")
	portNum = port[len(port)-1]
	util.CurrenciesStubPort = portNum // Change the global variable to the new port.

	// A Dashboard to compare with the response of the service. (Expected response)
	dashboard := util.DashboardResponse{
		Name:    "Norway",
		Isocode: "NO",
		Features: util.DashboardFeatures{
			Temperature:   13.276785714285708,
			Precipitation: 0.0375,
			Capital:       "Oslo",
			Coordinates: util.Coordinates{
				Latitude:  62.0,
				Longitude: 10.0,
			},
			Population: 5379475,
			Area:       323802,
			TargetCurrencies: map[string]float64{
				"NOK": 1.0,
				"EUR": 0.086289,
			},
		},
		LastRetrieval: time.Now().Format("2006-01-02 15:04"),
	}

	// Populate JSON-file for testing (Give an example registration)
	if err := populateDashboardsFile(); err != nil {
		t.Fatalf("Failed to populate the test data. %v\n", err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler.DashboardHandler))
	defer server.Close()

	client := http.Client{}

	// -----
	// Test illegal methods
	// -----
	if err := util.TestMethod(client, server.URL, http.MethodPost, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+util.DASHBOARD_PATH, http.MethodPut, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+util.DASHBOARD_PATH, http.MethodDelete, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+util.DASHBOARD_PATH, http.MethodPatch, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+util.DASHBOARD_PATH, http.MethodConnect, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+util.DASHBOARD_PATH, http.MethodHead, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+util.DASHBOARD_PATH, http.MethodOptions, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}

	// -----
	// Create client and GET from server
	// -----
	id := "1" // A ID to search after (Given in the populated dashboard)
	res, err := getFromServer(server.URL + util.DASHBOARD_PATH + id)
	if err != nil {
		t.Fatalf("Failed to instantiate a new request.\n%v\n", err)
	}

	defer func(Body io.ReadCloser) {
		if Body != nil {
			if err := Body.Close(); err != nil {
				log.Println("Failed to close response body")
			}
		}
	}(res.Body)

	// Decode to struct
	var returningDashboardResponse util.DashboardResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&returningDashboardResponse); err != nil {
		t.Fatalf(
			"The stub database has no response, or the response cannot be decoded to a notification struct.\n%v\n",
			err)
	}

	// Set the LastRetrieval to an empty string, because of potential delay.
	dashboard.LastRetrieval = ""
	returningDashboardResponse.LastRetrieval = ""
	// Check whether the returning object has the properties we look for
	isObjectCorrect := reflect.DeepEqual(returningDashboardResponse, dashboard)
	if false == isObjectCorrect {
		t.Error("The data are not as expected!")
	}
}
