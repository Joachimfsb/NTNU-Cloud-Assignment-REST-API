package handler_test

import (
	"assignment2/handler"
	"assignment2/models"
	stubs "assignment2/stubs/handler"
	"assignment2/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// populateNotificationsFile ensures that tests are predictable. It replaces the database stub's notification JSON file with a fixed set of
// notification objects. This lets tests hardcode expected values in the objects.
func populateNotificationsFile() error {
	// Register the notifications to test with
	notifications := []models.NotificationDatabaseModel{
		{
			Id:      "1",
			Url:     "https://1.no/1",
			Event:   "INVOCATION",
			Country: "RU",
		},

		{
			Id:      "2",
			Url:     "https://2.no/2",
			Event:   "CREATE",
			Country: "NO",
		},

		{
			Id:      "3",
			Url:     "https://3.no/3",
			Event:   "CREATE",
			Country: "NO",
		},
	}

	return util.PopulateTestFile(util.STUB_DATABASE_NOTIFICATIONS, notifications)
}

// getFromServer creates a GET request to the server and returns the response
// Callee must manually close the response body.
//
// Parameters:
// - serverURL: The whole URL to the endpoint. Example: http://localhost:8080/dashboards
//
// Returns:
// - A pointer to the response is returned if the request was successful.
// - An error is returned if the request failed.
//
// Notes:
// The callee is responsible for closing the response body manually. Example:
//
// res, err := getFromServer("http://endpoint.net")
// defer res.Body.Close()
func getFromServer(serverURL string) (*http.Response, error) {
	// Create client instance
	client := http.Client{}

	// Retrieve content from server
	res, err := client.Get(serverURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications from the server. %v", err)
	}

	return res, nil
}

// postToServer creates a POST to the server and sends an object to it.
// Callee must manually close the response body
//
// Parameters:
// - serverURL: The whole URL to the endpoint. Example: http://localhost:8080/dashboards
// - obj: The struct to send to the server. This is automatically encoded to a supported format by this function.
//
// Returns:
// - A pointer to the response is returned if the request was successful.
// - An error is returned if the request failed.
//
// Notes:
// The callee is responsible for closing the response body manually. Example:
//
// res, err := getFromServer("http://endpoint.net")
// defer res.Body.Close()
func postToServer(serverURL string, obj interface{}) (*http.Response, error) {
	// Create client instance
	client := http.Client{}

	// Encode object to JSON
	jsonTxt, err := json.Marshal(&obj)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal notiCorrectCasing. This is a developer error.\n%v\n", err)
	}

	// Register the object
	// https://stackoverflow.com/a/24455606
	req, err := http.NewRequest(http.MethodPost, serverURL, bytes.NewBuffer(jsonTxt))
	if err != nil {
		// We somehow get an error. Developer mistake?
		return nil, fmt.Errorf("Failed to instantiate a new request. Developer mistake?\n%v\n", err)
	}
	res, err := client.Do(req)
	if err != nil {
		// Connection error
		return nil, fmt.Errorf("Failed to connect to the server.\n%v\n", err)
	}

	return res, nil
}

// deleteToServer creates a DELETE to the server.
// Callee must manually close the response body.
//
// Parameters:
// - serverURL: The whole URL to the endpoint. Example: http://localhost:8080/dashboards
//
// Returns:
// - A pointer to the response is returned if the request was successful.
// - An error is returned if the request failed.
//
// Notes:
// The callee is responsible for closing the response body manually. Example:
//
// res, err := getFromServer("http://endpoint.net")
// defer res.Body.Close()
func deleteToServer(serverURL string) (*http.Response, error) {
	// Create client instance
	client := http.Client{}

	// Register the object
	// https://stackoverflow.com/a/24455606
	req, err := http.NewRequest(http.MethodDelete, serverURL, nil)
	if err != nil {
		// We somehow get an error. Developer mistake?
		return nil, fmt.Errorf("Failed to instantiate a new request. Developer mistake?\n%v\n", err)
	}
	res, err := client.Do(req)
	if err != nil {
		// Connection error
		return nil, fmt.Errorf("Failed to connect to the server.\n%v\n", err)
	}

	return res, nil
}

// TestNotificationsIllegalMethods tests whether the server handles illegal HTTP methods correctly
func TestNotificationsIllegalMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handler.NotificationHandler))
	defer server.Close()

	client := http.Client{}

	// -----
	// Test illegal methods
	// -----
	if err := util.TestMethod(client, server.URL+"/notifications", http.MethodPut, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+"/notifications", http.MethodPatch, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+"/notifications", http.MethodConnect, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+"/notifications", http.MethodHead, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
	if err := util.TestMethod(client, server.URL+"/notifications", http.MethodOptions, http.StatusMethodNotAllowed); err != nil {
		t.Error(err)
	}
}

// TestGetOneNotification tests whether the server can return a single notification.
// The test checks for each field in a notification. They must match exactly to the test parameters.
func TestGetOneNotification(t *testing.T) {
	util.FixStubPaths()

	if err := populateNotificationsFile(); err != nil {
		t.Fatalf("Failed to populate the test data. %v\n", err)
	}

	util.Config.Stubs.Database = true
	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseNotificationHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	server := httptest.NewServer(http.HandlerFunc(handler.NotificationHandler))
	if server != nil {
		defer server.Close()
	} else {
		t.Fatalf("failed to create an instance of the server. It is nil")
	}

	res, err := getFromServer(server.URL + "/2")
	if err != nil {
		t.Fatalf("Failed to instantiate a new request.\n%v\n", err)
	}

	defer func(Body io.ReadCloser) {
		if Body != nil {
			if err := Body.Close(); err != nil {
				log.Println("Failed to close repsonse body")
			}
		}
	}(res.Body)

	// The server should return OK
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v", res.Status)
	}

	// Decode to structs array
	var notification models.NotificationDatabaseModel
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&notification); err != nil {
		t.Fatalf(
			"The stub database has no response, or the response cannot be decoded to a notification struct.\n%v\n",
			err)
	}

	objectToCheckAgainst := models.NotificationDatabaseModel{
		Id:      "2",
		Url:     "https://2.no/2",
		Event:   "CREATE",
		Country: "NO",
	}

	// Check whether the returning object has the properties we look for
	isObjectCorrect := reflect.DeepEqual(objectToCheckAgainst, notification)

	if false == isObjectCorrect {
		t.Error("notification object's fields were incorrect.")
	}
}

// TestGetAllNotifications tests whether the server has implemented the functionality to retrieve all notifications.
// The test will check for three notification samples, and all the fields must match.
func TestGetAllNotifications(t *testing.T) {
	util.FixStubPaths()

	util.Config.Stubs.Database = true

	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseNotificationHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	if err := populateNotificationsFile(); err != nil {
		t.Fatalf("Failed to populate the test data. %v\n", err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler.NotificationHandler))
	defer server.Close()

	// Get response from server
	res, err := getFromServer(server.URL)
	if err != nil {
		t.Fatalf("Failed to instantiate a new request.\n%v\n", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("Failed to close repsonse body")
		}
	}(res.Body)

	// The server should return status OK
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.Status)
	}

	// Decode to structs array
	var notifications []models.NotificationDatabaseModel
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&notifications); err != nil {
		t.Fatalf(
			"The stub database has no response, or the response cannot be decoded to a notification struct.\n%v\n",
			err)
	}

	// There must be at least three notifications.
	if len(notifications) < 3 {
		t.Fatalf("There are less than 3 notifications. There should be at least 3 for testing purposes.")
	}

	notificationsCorrect := []models.NotificationDatabaseModel{
		{
			Id:      "1",
			Url:     "https://1.no/1",
			Event:   "INVOCATION",
			Country: "RU",
		},

		{
			Id:      "2",
			Url:     "https://2.no/2",
			Event:   "CREATE",
			Country: "NO",
		},

		{
			Id:      "3",
			Url:     "https://3.no/3",
			Event:   "CREATE",
			Country: "NO",
		},
	}

	// Test the first returning object
	isObjectCorrect1 := reflect.DeepEqual(notifications[0], notificationsCorrect[0])
	isObjectCorrect2 := reflect.DeepEqual(notifications[1], notificationsCorrect[1])
	isObjectCorrect3 := reflect.DeepEqual(notifications[2], notificationsCorrect[2])

	if !isObjectCorrect3 || !isObjectCorrect2 || !isObjectCorrect1 {
		t.Fatal("notification object's fields were incorrect")
	}

	// Test what happens when GetAllNotifications handler is run when there are no notifications in the database.
	// Expected error code is 200 OK and an empty string
}

// TestRegisterNotification tests whether a notification can successfully be registered.
// It tests the following:
// - The fields "Event" and "Country" are always uppercased
// - Whether the objects are actually stored correctly
//
// TODO test when an object without URL is registered. Expected status code: 422 Unprocessable Entity
func TestRegisterNotification(t *testing.T) {
	util.FixStubPaths()

	util.Config.Stubs.Database = true

	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseNotificationHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	// Notification to register. The server handle incorrect casing
	notiWrongCasing := models.NotificationDTO{
		Url:     "https://4.no/4",
		Event:   "invoke",
		Country: "rU",
	}

	notiCorrectCasing := models.NotificationDTO{
		Url:     "https://5.no/5",
		Event:   "REGISTER",
		Country: "NO",
	}

	server := httptest.NewServer(http.HandlerFunc(handler.NotificationHandler))
	defer server.Close()

	// POST and register notification with incorrect casing
	res, err := postToServer(server.URL, notiWrongCasing)
	if err != nil {
		t.Fatalf("Failed to register notification with wrong casing.\n%v\n", err)
	}

	// Response should have HTTP status 201 Created
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201 Created, got %v\n", res.Status)
	}

	// Check if the server successfully registered the new notification
	// POST and register notification with correct casing
	res1, err := postToServer(server.URL, notiCorrectCasing)
	if err != nil {
		t.Fatalf("Failed to register notification with correct casing.\n%v\n", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("Failed to close response body")
		}
	}(res1.Body)

	// Response should have HTTP status 201 Created
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201 Created, got %v\n", res.Status)
	}
}

// TestDeleteNofication deletes a notification from the database if it exists.
// The following are tested:
// - Deleting a notification removes it from the database
// - Deleting a non-existing notification still returns error code 204 No Content
func TestDeleteNotification(t *testing.T) {
	util.FixStubPaths()

	util.Config.Stubs.Database = true

	stubServer := httptest.NewServer(http.HandlerFunc(stubs.DatabaseNotificationHandler))
	defer stubServer.Close()

	port := strings.Split(stubServer.URL, ":")
	portNum := port[len(port)-1]
	util.DatabaseStubPort = portNum

	server := httptest.NewServer(http.HandlerFunc(handler.NotificationHandler))
	defer server.Close()

	// -----
	// Ensure that the notification object exists in the database
	// -----
	res, err := getFromServer(server.URL + "/1")
	if err != nil {
		t.Fatalf("Failed to instantiate a new request.\n%v\n", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("Failed to close repsonse body")
		}
	}(res.Body)

	// The server should return OK
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %v. Thus, unable to check deletion of a notification", res.Status)
	}

	// -----
	// Delete notification
	// -----
	if _, err = deleteToServer(server.URL + "/1"); err != nil {
		t.Errorf("Failed to delete a notification by ID = 1.\n%v\n", err)
	}

	// -----
	// Delete the same object again. The server should still return status code 201
	// -----
	if res, err = deleteToServer(server.URL + "/1"); err != nil {
		t.Errorf("Failed to delete a notification by ID = 1.\n%v\n", err)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Deleting a notification that does not exist should yield status code 204 Not Created. Expected %v got %v", http.StatusNoContent, res.StatusCode)
	}
}
