package notifications

import (
	"assignment2/database"
	"assignment2/models"
	"assignment2/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// InvokeEvent ensures notifications are sent based on the event type.
//
// # Overview
//
// The system will find all events related to the input country. Consider using 'util.EVENT_*' for correct event type.
// Rules:
// - A notification may not be registered to a specific country.
//
// # Example
//
// notifications.InvokeEvent("Norway", util.EVENT_REGISTRATION)
//
// Output: No output, but all dashboard configurations with "NO" will be found.
func InvokeEvent(country string, event string) error {
	validateEvent := util.ValidateEvents(event)

	// This event should only be used internally. If an error is invoked, it means developers have made a mistake, and
	// it must be fixed.
	if validateEvent != nil {
		log.Println("assignment2/models/invocations.go: An error occurred in function InvokeEvent(). This should not happen.")
		log.Printf("Event to validate: %v\n", event)
		return fmt.Errorf("failed to validate event %v. This should not happen", event)
	}

	databaseModels, err := database.GetAllNotifications()
	if err != nil {
		return fmt.Errorf("unable to get all notifications. This should not happen. %v", err)
	}

	var filteredModels []models.NotificationDatabaseModel
	for _, model := range databaseModels {

		// Skip non-matching or non-empty country
		if model.Country != "" && model.Country != country {
			continue
		}

		// Skip non-matching event or non-empty event
		if model.Event != "" && event != model.Event {
			continue
		}

		filteredModels = append(filteredModels, model)
	}

	// Send notification to all registered URLs
	if len(filteredModels) > 0 {
		// Create a HTTP post request
		for _, model := range filteredModels {
			err = sendNotification(model)
			if err != nil {
				// Something went wrong with sending the notification to the client
				log.Println("assignments02/notifications/invocations.go: Function InvokeEvent() failed, to send notification to client.")
				log.Println(err)
			}
		}
	}

	// At this point, no match was found in the database. Nothing happens
	return nil
}

// sendNotification is an internal function that sends POST request to a notification's URL.
func sendNotification(n models.NotificationDatabaseModel) error {
	// Select attributes to send to the client
	out := models.InvocationNotificationModel{
		Id:      n.Id,
		Country: n.Country,
		Event:   n.Event,
		Time:    time.Now(),
	}

	marshalledOutput, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal notification by id: %v", n.Id)
	}

	// https://www.kirandev.com/http-post-golang
	r, err := http.NewRequest("POST", n.Url, bytes.NewBuffer(marshalledOutput))
	if err != nil {
		return fmt.Errorf("failed to create request. This should not happen: %v", err)
	}

	r.Header.Add(util.CONTENT_TYPE, util.MIMETYPE_JSON)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("invocation [%v] was unable to reach the client URL at: %v", n.Event, n.Url)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close connection to %v", n.Url)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		// Unable to connect to the client
		return fmt.Errorf("the connection to notification client by ID %v was unsuccessful", n.Id)
	}

	return nil
}
