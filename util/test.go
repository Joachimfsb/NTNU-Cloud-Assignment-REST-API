// Helper functions for testing

package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// TestMethod is a helper function that tests whether a server implements a certain method for an endpoint.
//
// Parameters:
//   - client: the client to send the request from. This enables reusing the same client, which is the best practise
//     according to the official 'net/http' documentation.
//   - endpoint: the endpoint to send the request to. Example: http://example.com/foo/bar
//   - method: the method to test. Use built-in constants, like http.MethodPost
//   - expectedStatusCode: the expected status code which is returned from the server
//
// Returns:
//   - If there was an error while connecting to the server on the endpoint, there will be an error raised
//   - If the method does not match, then an error is also raised
func TestMethod(client http.Client, endpoint string, method string, expectedStatusCode int) error {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		// We somehow get an error. Developer mistake?
		return fmt.Errorf("Failed to instantiate a new request. Incorrect URL?\n%v\n", err)
	}

	defer func(Body io.ReadCloser) {
		if Body != nil {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close request body. Incorrect URL?\n%v\n", err)
			}
		}
	}(req.Body)

	resPost, err := client.Do(req)
	if err != nil {
		// Connection error
		return fmt.Errorf("Failed to connect to the server.\n%v\n", err)
	}

	// Test status code
	if resPost.StatusCode != expectedStatusCode {
		return fmt.Errorf("Unexpected status code. Expected %d, got %d\n", expectedStatusCode, resPost.StatusCode)
	}

	return nil
}

// PopulateTestFile populates the database's stub's JSON file where all data is stored.
// This ensures that whenever running a test, the file is exactly the way the test expects it to be.
//
// Parameters:
// - filepath: location where the populated file should be. Example: /path/to/file.json
// - obj: array of objects to populate the file with.
//
// Returns:
// - An error is returned if the file failed to be populated with the array of objects.
// - Error is nil if everything was OK.
func PopulateTestFile(filepath string, obj any) error {
	// Ensure that the filepaths are valid
	if filepath != STUB_DATABASE_REGISTRATIONS && filepath != STUB_DATABASE_NOTIFICATIONS {
		return fmt.Errorf("invalid filepath. Must be util.STUB_DATABASE_REGISTRATIONS or util.STUB_DATABASE_NOTIFICATIONS")
	}

	// Encode the object array to JSON
	jsonText, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON. %v", err)
	}

	// Write the JSON to the specified filepath
	err = os.WriteFile(filepath, jsonText, 0755)
	if err != nil {
		return fmt.Errorf("unable to write file. %v", err)
	}

	return nil
}

// FixStubPaths fixes the paths of the database stub JSON paths. This is due to errors in CWD.
func FixStubPaths() {
	STUB_DATABASE_REGISTRATIONS = "../stubs/res/registrations.json"
	STUB_DATABASE_NOTIFICATIONS = "../stubs/res/notifications.json"

	STUB_WEATHER_REPONSE = "../stubs/res/weather.json"
	STUB_CURRENCIES_RESPONSE = "../stubs/res/currency.json"
	STUB_COUNTRY_RESPONSE = "../stubs/res/country.json"
}
