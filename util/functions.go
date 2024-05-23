package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// GetIdFromUrl attempts to get the registration or notification ID from a URL.
// It is very difficult to create a general function to get ID, regardless of what the base URl is.
// This is because when no ID is in the URL, the function may accidentally parse the URL's last segment as the ID.
// To combat this issue, we force all IDs to have numerical values.
func GetIdFromUrl(clientUrl string) (string, error) {
	// Add trailing '/'
	if clientUrl[len(clientUrl)-1] == '/' {
		clientUrl = clientUrl[:len(clientUrl)-1]
	}

	// GET request means to return notification(s). If an ID is provided, then we return only a single notification.
	// Otherwise, all notifications will be returned.
	url := strings.Split(clientUrl, "/")

	id := url[len(url)-1]

	if id == "notifications" || id == "registrations" {
		return "", nil
	}

	urlID := id

	for i := 0; i < len(urlID); i++ {
		_, err := strconv.Atoi(strconv.Itoa(int(urlID[i])))
		if err != nil {
			fmt.Println(urlID)
			fmt.Printf("Unexpected type of %v", urlID[i])
			return "error", fmt.Errorf("%q looks not like a number.\n", id)
		}
	}

	return id, nil
}

// MakeGetRequest creates a GET-request to a URL and decodes the body of the response into content.
func MakeGetRequest(url string, content any) error {
	// Make and issue a new GET-request
	res, err1 := http.Get(url)
	if err1 != nil {
		return err1
	}

	err2 := decodeBody(res.Body, content) // Decode the body, and store it in the original variable/struct
	if err2 != nil {
		return err2
	}
	// Everything is OK
	return nil
}

// decodeBody decodes the body of a given response into a given value (content).
func decodeBody(body io.ReadCloser, content any) error {
	decoder := json.NewDecoder(body) // Initialize the decoder
	// Decode the body into the data-type
	if err1 := decoder.Decode(content); err1 != nil {
		return err1
	}
	err2 := body.Close() // Closing the body.
	if err2 != nil {
		return err2
	}
	// Everything is OK
	return nil
}

// ParseFile is to read through a file and return the data.
func ParseFile(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error reading file: ", err)
		os.Exit(1)
	}
	return data
}
