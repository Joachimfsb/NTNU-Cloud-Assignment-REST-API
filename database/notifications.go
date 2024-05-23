// Notifications.go is a module that consists of data fetching from our database
// related to notifications and webhooks
package database

import (
	"assignment2/models"
	"assignment2/util"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/iterator"
)

// AddNotification adds a new notification object to the database.
//
// Parameters:
// - model: the new notification object to add to the database.
//
// Returns:
// An error object is returned if the notification could not be stored. This is normally caused by developer mistakes.
func AddNotification(model models.NotificationDatabaseModel) error {
	if util.Config.Stubs.Database == false {
		// Uppercase Event type as a good practise
		model.Event = strings.ToUpper(model.Event)
		model.Country = strings.ToUpper(model.Country)

		_, _, err := Client.Collection(util.COLLECTION_NOTIFICATIONS).Add(Ctx, model)
		if err == nil {
			return err
		}
		return fmt.Errorf("unable to add notification %v", err)
	} else {
		client := http.Client{}

		encodedModel, err := json.Marshal(&model)
		if err != nil {
			return fmt.Errorf("Unable to marshal notiCorrectCasing. This is a developer error.\n%v\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "http://localhost:"+util.DatabaseStubPort+util.NOTIFICATION_PATH, bytes.NewBuffer(encodedModel))
		if err != nil {
			return fmt.Errorf("request not compatible with client.Do %v", err)
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending delete request:", err)
			return err
		}
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				log.Println("Failed to close repsonse body")
			}
		}(res.Body)
		if err == nil {
			return err
		}
		return fmt.Errorf("unable to add notification %v", err)
	}
}

// GetAllNotifications retrieves an array of notifications from the database.
// The callee must be aware that if no notifications were found, then an empty array is returned.
// An error object is returned if something wrong happens. This is normally caused by developer mistakes.

// Credit:
// https://firebase.google.com/docs/firestore/query-data/get-data#go
func GetAllNotifications() ([]models.NotificationDatabaseModel, error) {
	var out []models.NotificationDatabaseModel
	if util.Config.Stubs.Database == false {

		iter := Client.Collection(util.COLLECTION_NOTIFICATIONS).Documents(Ctx)
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return nil, err
			}
			model := models.NotificationDatabaseModel{}
			model.PopulateFromMap(doc.Data())

			// If the ID is '123123', we must not return it. This is only used for unit testing.
			if model.Id != "123123" {
				out = append(out, model)
			}
		}
	} else {
		client := http.Client{}

		// Retrieve content from server
		res, err := client.Get("http://localhost:" + util.DatabaseStubPort + util.NOTIFICATION_PATH)
		if res == nil {
			return []models.NotificationDatabaseModel{}, fmt.Errorf("res is nil %v", err)
		}
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				log.Println("Failed to close repsonse body")
			}
		}(res.Body)
		if err != nil {
			return nil, err
		}

		// TODO if the bug is malformed, then an error occurs. This is usually a sign that no notifications were found, although
		// it does mean that the file's format is invalid.

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&out); err != nil {
			return nil, fmt.Errorf(
				"the stub database has no response, or the response cannot be decoded to a notification struct.\n%v\n",
				err)
		}
	}

	return out, nil
}

// GetSingleNotification retrieves a single notification from the database. The notification is queried by its id.
// If a notification was found, then a notification object is returned. If not, then an error is returned.
// An error is also returned if other bad events occur. These events are normally developer mistakes.
//
// Parameters:
// - id: the ID to query by.
func GetSingleNotification(id string) (models.NotificationDatabaseModel, error) {
	var out models.NotificationDatabaseModel
	if util.Config.Stubs.Database == false {
		iter := Client.Collection(util.COLLECTION_NOTIFICATIONS).Documents(Ctx)
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return models.NotificationDatabaseModel{}, err
			}

			// If the input ID matches the ID found in the database, then return this model object.
			if doc.Data()["Id"] != nil && doc.Data()["Id"].(string) == id {
				out.PopulateFromMap(doc.Data())
				break
			}
		}
	} else {
		client := http.Client{}

		// Retrieve content from server
		res, err := client.Get("http://localhost:" + util.DatabaseStubPort + util.NOTIFICATION_PATH + id)
		if res != nil {

			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					log.Println("Failed to close repsonse body")
				}
			}(res.Body)
		}
		if err != nil {
			return models.NotificationDatabaseModel{}, err
		}

		// TODO potential nil pointer error with res.Body
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&out); err != nil {
			return models.NotificationDatabaseModel{}, fmt.Errorf(
				"The stub database has no response, or the response cannot be decoded to a notification struct.\n%v\n",
				err)
		}
	}

	return out, nil
}

// DeleteNotification deletes/removes a single notification from the database. The notification is queried by its id.
// If a notification was found, then a notification object is deleted. If not, then an error is returned.
// An error is also returned if other bad events occur. These events are normally developer mistakes.
//
// Parameters:
// id: the Registration object's ID to search for.
//
// Returns:
// error: any errors which may occur
func DeleteNotification(id string) error {
	if util.Config.Stubs.Database == false {
		iter := Client.Collection(util.COLLECTION_NOTIFICATIONS).Documents(Ctx)
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return err
			}

			// If the input ID matches the ID found in the database, then delete the document
			if doc.Data()["Id"] != nil && doc.Data()["Id"].(string) == id {
				_, err := doc.Ref.Delete(Ctx)
				return err
			}
		}
	} else {
		var out models.NotificationDatabaseModel
		client := http.Client{}

		// Retrieve content from server
		req, err := http.NewRequest("DELETE", "http://localhost:"+util.DatabaseStubPort+util.NOTIFICATION_PATH+id, nil)

		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending delete request:", err)
			return err
		}
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				log.Println("Failed to close repsonse body")
			}
		}(res.Body)

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&out); err != nil {
			return fmt.Errorf(
				"The stub database has no response, or the response cannot be decoded to a notification struct.\n%v\n",
				err)
		}
	}
	return fmt.Errorf("unable to delete notification by id: %v", id)
}
