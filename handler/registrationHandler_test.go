package handler_test

import (
	"assignment2/handler"
	stubs "assignment2/stubs/handler"
	"assignment2/util"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// populateRegistrationFile ensures that tests are predictable. It replaces the database stub's registration
// JSON file with a registration object. This lets tests hardcode expected values in the objects.
func populateRegistrationFile() error {
	var allRegistrations []util.Registration
	registration := util.Registration{
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
			TargetCurrencies: []string{"EUR", "USD", "SEK"},
		},
		LastChange: "2024-04-10 14:09", // Time is hardcoded as per the input
	}
	allRegistrations = append(allRegistrations, registration)

	return util.PopulateTestFile(util.STUB_DATABASE_REGISTRATIONS, allRegistrations)
}

// TestRegistrationPostHandler tests POST method for registration of a dashboard.
// It verifies:
// 1. Correct handling for a valid input.
// 2. Appropriate responses when met with an invalid input.
//
// Test checks that the registration of a dashboard is correctly processed, returning 201 (created)
// and ID + LastChange in the response. Also test invalid input, returning 400 (bad request).
func TestRegistrationPostHandler(t *testing.T) {
	util.FixStubPaths()

	//Enable database stub
	util.Config.Stubs.Database = true
	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseDashboardHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	//Test HTTP server setup with route to RegistrationHandler
	server := httptest.NewServer(http.HandlerFunc(handler.RegistrationHandler))
	defer server.Close()

	//*******VALID TESTING*******
	//New object with expected BODY for a request
	newDashboard := util.Registration{
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
	//Marshal object into JSON & create request/response
	body, _ := json.Marshal(newDashboard)
	request := httptest.NewRequest(http.MethodPost, server.URL+util.REGISTRATION_PATH, bytes.NewBuffer(body))
	responseRecorder := httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	//Chekc for 201 (status created)
	if status := responseRecorder.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	//Check for expected response returned (ID and lastChange)
	var result map[string]string
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &result); err != nil {
		t.Fatal("Could not decode response:", err)
	}
	if result["id"] == "" || result["lastChange"] == "" {
		t.Error("Missing id or lastChange in response")
	}

	//*******INVALID TESTING*******
	// Create a malformed registration object (invalid data)
	invalid := []byte(`{
        "Country": 123,
        "IsoCode": ["NO"],
        "Features": {
            "Temperature": "NotABool",
            "Precipitation": "NotABool",
            "Capital": "NotABool",
            "Coordinates": "NotABool",
            "Population": "NotABool",
            "Area": "NotABool",
            "TargetCurrencies": "NOK, EUR"
        }
    }`)

	// Marshal object into JSON & create request/response
	badBody, _ := json.Marshal(invalid)
	badRequest := httptest.NewRequest(http.MethodPost, server.URL+util.REGISTRATION_PATH, bytes.NewBuffer(badBody))
	badResponseRecorder := httptest.NewRecorder()
	handler.RegistrationHandler(badResponseRecorder, badRequest)

	// Check for expected failure status code (e.g., 400 Bad Request)
	if status := badResponseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid input, got %d", http.StatusBadRequest, status)
	}
}

// TestRegistrationGetHandlerAll tests the GET method for getting ALL dashboards (empty path /).
// Verifies that server responds with retrieval of all stored dashboards.
//
// Test checks that a successful retrieval returns status code 200 (OK), and error if no dashboards
// are registered.
func TestRegistrationGetHandlerAll(t *testing.T) {
	util.FixStubPaths()

	if err := populateRegistrationFile(); err != nil {
		t.Fatalf("Failed to populate the test data. %v\n", err)
	}

	//Enable database stub
	util.Config.Stubs.Database = true
	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseDashboardHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	//Test HTTP server setup with route to RegistrationHandler
	server := httptest.NewServer(http.HandlerFunc(handler.RegistrationHandler))
	defer server.Close()

	//Get request for all registrations (empty path /)
	request := httptest.NewRequest(http.MethodGet, server.URL+util.REGISTRATION_PATH, nil)
	responseRecorder := httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	//Checks for statusOK
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d for getting all registrations, got %d", http.StatusOK, status)
	}

	//Checks body
	var dashboards []util.Registration
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &dashboards); err != nil {
		t.Fatal("Could not decode the response for all dashboards:", err)
	}

	//No dashboards
	if len(dashboards) == 0 {
		t.Errorf("No dashboards registered")
	}
}

// TestRegistrationGetHandlerSpecific tests the GET method for specified ID (/ID)
// It verifies:
// 1. Correct handling for a valid ID.
// 2. Appropriate responses when met with an invalid ID.
//
// Test checks that the retrieval of a dashboard is correctly processed, returning the dashboard &
// 200 (OK) in the response. For unrecognized id 404 is returned (not found).
func TestRegistrationGetHandlerSpecified(t *testing.T) {
	util.FixStubPaths()

	if err := populateRegistrationFile(); err != nil {
		t.Fatalf("Failed to populate the test data. %v\n", err)
	}

	//Enable database stub
	util.Config.Stubs.Database = true
	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseDashboardHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	//Test HTTP server setup with route to RegistrationHandler
	server := httptest.NewServer(http.HandlerFunc(handler.RegistrationHandler))
	defer server.Close()

	//*******VALID TESTING*******
	//Assumes valid ID from POST, checks path /+ID for to get info from specified
	testID := "1"
	request := httptest.NewRequest(http.MethodGet, server.URL+util.REGISTRATION_PATH+testID, nil)
	responseRecorder := httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	//Checks for status 200
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d for getting registration by ID, got %d", http.StatusOK, status)
	}

	//Check response
	var registration util.Registration
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &registration); err != nil {
		t.Fatal("Could not decode response for registration by ID:", err)
	}

	//*******INVALID TESTING*******
	badID := "Random123"
	request = httptest.NewRequest(http.MethodGet, server.URL+util.REGISTRATION_PATH+badID, nil)
	responseRecorder = httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for invalid ID, got %d", http.StatusNotFound, responseRecorder.Code)
	}
}

// TestRegistrationDeleteHandler tests the DELETE method for specified ID.
// It verifies:
// 1. Correct handling for a valid ID.
// 2. Idempotency
//
// Whether the document exists =/= is found 204 (No content) is returned as delete is idempotent, this would also
// cover multiple requests to the same ID returning 204 (No content) which is checked in the test.
func TestRegistrationDeleteHandler(t *testing.T) {
	util.FixStubPaths()

	//Enable database stub
	util.Config.Stubs.Database = true
	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseDashboardHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	//Test HTTP server setup with route to RegistrationHandler
	server := httptest.NewServer(http.HandlerFunc(handler.RegistrationHandler))
	defer server.Close()

	//Assuming testID exists
	testID := "1"

	//Check if document exists before deleting
	preDelete := httptest.NewRequest(http.MethodGet, server.URL+util.REGISTRATION_PATH+testID, nil)
	preDeleteRecorder := httptest.NewRecorder()
	handler.RegistrationHandler(preDeleteRecorder, preDelete)

	if preDeleteRecorder.Code != http.StatusOK {
		t.Errorf("Check failed expected: %d, got %d", http.StatusOK, preDeleteRecorder.Code)
	}

	request := httptest.NewRequest(http.MethodDelete, server.URL+util.REGISTRATION_PATH+testID, nil)
	responseRecorder := httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	//Checks for status
	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("Expected status code %d for successful deletion, got %d", http.StatusNoContent, status)
	}

	//Second DELETE request to test idempotency
	responseRecorder2 := httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder2, request)

	if responseRecorder2.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d for second deletion attempt, got %d", http.StatusNoContent, responseRecorder2.Code)
	}
}

// TestRegistrationPutHandler tests the PUT method for change of data in a specified dashboard.
// It verifies:
// 1. Correct handling for a valid ID.
// 2. Appropriate responses when met with an invalid ID.
// 3. Successfully updating when given valid data changes.
// 4. Appropriate error when a request is sent containing an invalid body.
//
// Upon a successful update containing a valid ID & JSON input 200 (OK) is returned, if the ID
// is unrecognized 404 (not found) is returned and for bad JSON 400 (bad request is returned with
// appropriate messages.
func TestRegistrationPutHandler(t *testing.T) {
	util.FixStubPaths()

	if err := populateRegistrationFile(); err != nil {
		t.Fatalf("Failed to populate the test data. %v\n", err)
	}

	//Enable database stub
	util.Config.Stubs.Database = true
	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseDashboardHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	//Test HTTP server setup with route to RegistrationHandler
	server := httptest.NewServer(http.HandlerFunc(handler.RegistrationHandler))
	defer server.Close()

	//*******VALID TESTING*******
	//Assumes valid ID from POST
	testID := "1"
	//Changed values from post request
	changedDashboard := util.Registration{
		Country: "Norway",
		IsoCode: "NO",
		Features: util.Features{
			Temperature:      false,
			Precipitation:    false,
			Capital:          false,
			Coordinates:      false,
			Population:       false,
			Area:             false,
			TargetCurrencies: []string{"NOK", "SWE"},
		},
	}

	//Marshal object into JSON & create request/response
	body, _ := json.Marshal(changedDashboard)
	request := httptest.NewRequest(http.MethodPut, server.URL+util.REGISTRATION_PATH+testID, bytes.NewBuffer(body))
	responseRecorder := httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	//Check for status 200 (OK)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	//*******INVALID TESTING*******
	//Invalid ID
	badID := "random123" //Unregistered ID
	request = httptest.NewRequest(http.MethodPut, server.URL+util.REGISTRATION_PATH+badID, bytes.NewBuffer(body))
	responseRecorder = httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	//Check for status 404 (Not Found)
	if status := responseRecorder.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d for unrecognized ID, got %d", http.StatusNotFound, status)
	}

	//Invalid JSON
	//Test with invalid JSON
	invalid := []byte(`{"This should not work"`) // Bad JSON
	request = httptest.NewRequest(http.MethodPut, server.URL+util.REGISTRATION_PATH+testID, bytes.NewBuffer(invalid))
	responseRecorder = httptest.NewRecorder()
	handler.RegistrationHandler(responseRecorder, request)

	//Check for status 400 (Bad Request)
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid JSON, got %d", http.StatusBadRequest, status)
	}
}
