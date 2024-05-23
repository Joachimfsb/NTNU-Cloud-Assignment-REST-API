package handler

import (
	"assignment2/crypto"
	"assignment2/database"
	"assignment2/models"
	"assignment2/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// NotificationHandler is the main entry point for the notification endpoint.
// It handles the following:
// - Retrieving a single notification by its id
// - Retrieving all notifications in the database
// - Deleting a notification
// - Registering a notification
//
// If an illegal method is used, an appropriate message is sent to the client.
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, _ := util.GetIdFromUrl(r.URL.Path)

		if id == "" {
			getAllNotifications(w)
		} else {
			// Pass the ID provided by the client.
			getSingleNotification(w, id)
		}
	case http.MethodPost:
		registerNotification(w, r)
	case http.MethodDelete:
		id, _ := util.GetIdFromUrl(r.URL.Path)

		if id == "" {
			util.HttpError(w, "method 'DELETE' requires parameter 'ID'", http.StatusUnprocessableEntity)
			return
		} else {
			deleteNotification(w, id)
		}

	default:
		util.HttpError(w,
			"This method is not supported! Only POST, GET and DELETE are supported",
			http.StatusMethodNotAllowed)
		return
	}
}

// getAllNotifications returns all notifications to the client.
// If no notifications are found, then an empty array is returned.
func getAllNotifications(w http.ResponseWriter) {
	notifications, err := database.GetAllNotifications()

	if err != nil {
		log.Printf("Error getting all notifications: %v\n", err)
		util.HttpError(w, "no notifications were found", http.StatusInternalServerError)
		return
	}

	// There might not be any registered notifications at all... return an empty array
	if len(notifications) == 0 {
		w.Header().Set(util.CONTENT_TYPE, util.MIMETYPE_JSON)
		w.Header().Set(util.X_CONTENT_TYPE_OPTION, "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "[]")
		return
	}

	// Marshall the returning struct from the database
	marshalled, err := json.Marshal(notifications)
	if err != nil {
		log.Println("Failed to marshall notifications structs in getAllNotifications().")
		log.Println(err)
		util.HttpError(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set(util.CONTENT_TYPE, util.MIMETYPE_JSON)
	w.Header().Set(util.X_CONTENT_TYPE_OPTION, "nosniff")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, string(marshalled))
	if err != nil {
		log.Printf("unable to encode data from database. %v\n", err)
		util.HttpError(w, "unable to get all notifications", http.StatusInternalServerError)
		return
	}
}

// getSingleNotification retrieves a single notification by its id.
// If the notification was not found, 404 not found is returned to the client.
func getSingleNotification(w http.ResponseWriter, id string) {
	notification, err := database.GetSingleNotification(id)

	// Error from the database
	if err != nil {
		log.Println(err)
		util.HttpError(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}

	// The notification might be empty. It means no notificaiton was found
	if notification.Id == "" {
		util.HttpError(w, "no notification is found", http.StatusNotFound)
		return
	}

	// Marshall the returning struct from the database
	marshalled, err := json.Marshal(notification)
	if err != nil {
		log.Println("Failed to marshall notification structs in getAllNotifications().")
		log.Println(err)
		util.HttpError(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}

	// There is a notification stored in the database with the ID: "123123".
	// This is only for unit testing. Ideally we should make a check for this ID and prevent it from being returned to
	// the client. This is (probably) out of scope for this assignment, so this check is not implemented.
	// It is still relevant for the portfolio, so therefore there will be a TODO for this.
	// TODO implement a check for ID "123123". This is only used for unit testing.

	w.Header().Set(util.CONTENT_TYPE, util.MIMETYPE_JSON)
	w.Header().Set(util.X_CONTENT_TYPE_OPTION, "nosniff")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, string(marshalled))
	if err != nil {
		log.Printf("could not return marshalled object.%v", err)
		util.HttpError(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}
}

// registerNotification lets the client register a new notification.
//
// The client provides us a URL, so we can notify them whenever a certain event has occurred.
func registerNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(util.CONTENT_TYPE, util.MIMETYPE_JSON)

	// Expects incoming body in terms of WebhookRegistration struct
	dto := models.NotificationDTO{}
	err := json.NewDecoder(r.Body).Decode(&dto)

	// https://stackoverflow.com/a/32718077
	switch {
	case err == io.EOF:
		// empty body
		util.HttpError(w,
			"The request body was empty. Refer to the API documentation.",
			http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check for empty fields
	validationMessage := dto.ValidateFromClient()
	if validationMessage != nil {
		util.HttpError(w,
			validationMessage.Error(),
			http.StatusUnprocessableEntity)
		return
	}

	dto.Event = strings.ToUpper(dto.Event)
	dto.Country = strings.ToUpper(dto.Country)

	// Validate event
	validationMessage = util.ValidateEvents(dto.Event)
	if validationMessage != nil {
		util.HttpError(w,
			"field 'event' is invalid. "+validationMessage.Error(),
			http.StatusUnprocessableEntity)
		return
	}

	// Convert DTO fields to database model
	databaseModel := models.NotificationDatabaseModel{
		// ID cannot be provided by the client. If this still happens, we clear it. This is strictly not necessary, but
		// prevents potential bugs in this project's future development.
		Id:      myCrypto.GetMD5Hash(dto.Url + dto.Event + time.Now().String()),
		Url:     dto.Url,
		Event:   dto.Event,
		Country: dto.Country,
	}

	// Store the new notification to the database
	dbError := database.AddNotification(databaseModel)
	if dbError != nil {
		log.Println(dbError)
		util.HttpError(w, "unable to register the notification", http.StatusInternalServerError)
		return
	}

	// Write to client
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "{\"id\": \"%v\"}", databaseModel.Id)
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	} else {
		log.Println("Webhook " + databaseModel.Url + " has been registered.")
	}
}

// deleteNotification removes a notification by id from the database.
// It also protects notification by id '123123' to be deleted. This particular notification is used for unit testing.
func deleteNotification(w http.ResponseWriter, id string) {
	// The ID '123123' is protected and must never be deleted. This is used for unit
	// testing and is always stored on Firestore.
	if id == "123123" {
		log.Println("The client tried to delete notification with ID '123123'. This is a protected ID and must not be deleted at all.")
		http.Error(w, "", http.StatusNoContent)
		return
	}

	err := database.DeleteNotification(id)

	// Log any errors occurring, but don't tell the client about it. If a resource is deleted, they should safely
	// assume that it is deleted. Also, if a resource already is deleted, we should still return an OK status
	if err != nil {
		log.Println(err)
	}

	log.Printf("Deleted notification by id %v\n", id)

	http.Error(w, "", http.StatusNoContent)
}
