package stubs

import (
	myCrypto "assignment2/crypto"
	"assignment2/models"
	"assignment2/notifications"
	"assignment2/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// DatabaseDashboardHandler is the stub entry point for the registration endpoint.
// It handles the following:
// - Retrieving a single registration by its id
// - Retrieving all registrations in the database
// - Deleting a registration
// - Adding a registration to the database
// - Put/updating a registration in the database
//
// If an illegal method is used, an appropriate message is sent to the client.
func DatabaseDashboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, _ := util.GetIdFromUrl(r.URL.Path)

		if id == "" {
			stub_handleRegistrationGetAllRequest(w)
		} else {
			// Pass the ID provided by the client.
			stub_getSingleRegistration(w, id)
		}
	case http.MethodPost:
		stub_handleRegistrationPostRequest(w, r)
	case http.MethodPut:
		id, _ := util.GetIdFromUrl(r.URL.Path)

		if id != "" {
			// Pass the ID provided by the client.
			stub_handleRegistrationPutRequest(w, r, id)
		}
	case http.MethodDelete:
		id, _ := util.GetIdFromUrl(r.URL.Path)

		if id == "" {
			util.HttpError(w, "method 'DELETE' requires parameter 'ID'", http.StatusUnprocessableEntity)
			return
		} else {
			stub_handleRegistrationDeleteRequest(w, id)
		}
	default: //Error message if GET method is not used
		http.Error(w, "This method is not supported! Only POST, GET, PUT and DELETE are supported", http.StatusNotImplemented)
	}
}

// stub_handleRegistrationGetAllRequest returns all registrations to the client.
func stub_handleRegistrationGetAllRequest(w http.ResponseWriter) {
	file, err := os.Open(util.STUB_DATABASE_REGISTRATIONS) //opens json file that operates as an offline database
	if err != nil {                                        //error if file couldnt be opened
		log.Println(err)
		return
	}
	defer func(file *os.File) { //closes file after function
		err = file.Close()
		if err != nil { //error if pogram cant close file
			log.Println("could not close file")
		}
	}(file)

	var allRegistrations []util.Registration
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&allRegistrations) //adds all registrations to an array/slice
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Marshall the returning struct from the database
	marshalled, err := json.Marshal(&allRegistrations)
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
		log.Printf("failed to return marshalled object. %v", err)
		util.HttpError(w, "could not return marshalled object", http.StatusInternalServerError)
		return
	}
}

// stub_getSingleRegistration retrieves a single registration by its id.
// If the registration was not found, 404 not found is returned to the client.
func stub_getSingleRegistration(w http.ResponseWriter, id string) {
	file, err := os.Open(util.STUB_DATABASE_REGISTRATIONS) //opens file
	if err != nil {
		log.Println(err)
		return
	}
	defer func(file *os.File) { //closes file at the end of function
		err := file.Close()
		if err != nil {
			log.Println("could not close file")
		}
	}(file)

	var allNotifications []util.Registration
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&allNotifications) //takes all registrations and put them into a slice
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var notification util.Registration

	for _, not := range allNotifications {
		if not.ID == id {
			notification = not
			break
		}
	}

	// Marshall the returning struct from the database
	marshalled, err := json.Marshal(notification)
	if err != nil {
		log.Println("Failed to marshall notification structs in getAllNotifications().")
		log.Println(err)
		util.HttpError(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set(util.CONTENT_TYPE, util.MIMETYPE_JSON)
	w.Header().Set(util.X_CONTENT_TYPE_OPTION, "nosniff")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, string(marshalled))
	if err != nil {
		return
	}
}

// stub_handleRegistrationPostRequest lets the client register a new registration.
//
// The client provides us a URL, so we can notify them whenever a certain event has occurred.
func stub_handleRegistrationPostRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(util.CONTENT_TYPE, util.MIMETYPE_JSON)

	file, err := os.ReadFile(util.STUB_DATABASE_REGISTRATIONS) //opens file
	if err != nil {
		log.Println(err)
		return
	}

	var allRegistrations []util.Registration
	if err := json.Unmarshal(file, &allRegistrations); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Expects incoming body in terms of WebhookRegistration struct
	registrationModel := util.Registration{}
	err = json.NewDecoder(r.Body).Decode(&registrationModel)

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

	if registrationModel.ID == "" { //if id of registration is empty, id is created
		registrationModel.ID = myCrypto.GetMD5Hash(registrationModel.Country + time.Now().String())
	}
	registrationModel.LastChange = time.Now().Format("2006-01-02 15:04") //changes the last change og registration

	// Store the new notification to the database/json file
	allRegistrations = append(allRegistrations, registrationModel) //adding the new registration to all the registrations
	updatedAllRegistrations, err := json.MarshalIndent(allRegistrations, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	//with the help of ChatGPT. Overwriting the old file
	if err := os.WriteFile(util.STUB_DATABASE_REGISTRATIONS, updatedAllRegistrations, 0666); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("successfully added new registration")

	// Write to client
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "{\"id\": \"%v\"}", registrationModel.ID)
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	} else {
		log.Println("Registration " + registrationModel.Country + " has been registered.")
	}
}

// stub_deleteById lets the client delete a registration with an id
func stub_deleteById(id string) error {
	file, err := os.ReadFile(util.STUB_DATABASE_REGISTRATIONS) //open file
	if err != nil {
		log.Println(err)
		return err
	}

	var allRegistrations []util.Registration
	if err := json.Unmarshal(file, &allRegistrations); err != nil { //adding all registrations in file to slice
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	var allRegistrationsUpdated []util.Registration
	foundMatch := false
	for _, reg := range allRegistrations {
		if reg.ID != id { //if registration id dont match wanted id, registration is added to the updated list
			allRegistrationsUpdated = append(allRegistrationsUpdated, reg)
		} else { //if registration id matches wanted id
			foundMatch = true
		}
	}
	if foundMatch == false { //if match was not found, return error
		log.Println("could not find/delete wanted registration")
		return err
	}

	updatedRegList, err := json.MarshalIndent(allRegistrationsUpdated, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	// Write the updated data back to the file
	if err := os.WriteFile(util.STUB_DATABASE_REGISTRATIONS, updatedRegList, 0666); err != nil {
		return err
	}

	log.Printf("Deleted notification by id %v\n", id)
	return nil
}

// stub_handleRegistrationDeleteRequest removes a notification by id from the database.
func stub_handleRegistrationDeleteRequest(w http.ResponseWriter, id string) {
	if err := stub_deleteById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// stub_handleRegistrationPutRequest lets client update or add a registration to/in the database
func stub_handleRegistrationPutRequest(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Add(util.CONTENT_TYPE, util.MIMETYPE_JSON)

	//deletes old registration
	stub_deleteById(id)

	file, err := os.ReadFile(util.STUB_DATABASE_REGISTRATIONS)
	if err != nil {
		log.Println(err)
		return
	}
	//Decode the request body into a Registration struct (excluding ID and lastChange)
	registration := util.Registration{}
	err = json.NewDecoder(r.Body).Decode(&registration)

	//sets the "new" registration to the same id as url
	registration.ID = id

	// Invoke notification event for changing a registration
	if err := notifications.InvokeEvent(registration.IsoCode, util.EVENT_CHANGE); err != nil {
		log.Println("Error invoking event:", err)
	}

	// Update registration on the database
	var allRegistrations []util.Registration
	if err := json.Unmarshal(file, &allRegistrations); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	//Change the timestamp to last changed
	registration.LastChange = time.Now().Format("2006-01-02 15:04")
	// Store the new notification to the database/json file

	allRegistrations = append(allRegistrations, registration)
	updatedAllRegistrations, err := json.MarshalIndent(allRegistrations, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	//with the help of ChatGPT. Overwriting the old file
	if err := os.WriteFile(util.STUB_DATABASE_REGISTRATIONS, updatedAllRegistrations, 0666); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("successfully put old registration")

	// Write to client
	w.WriteHeader(http.StatusAccepted)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Registration " + registration.ID + " has been put/changed.")
}

// DatabaseNotificationHandler is the main entry point for the notification endpoint.
// It handles the following:
// - Retrieving a single notification by its id
// - Retrieving all notifications in the database
// - Deleting a notification
// - Registering a notification
//
// If an illegal method is used, an appropriate message is sent to the client.
func DatabaseNotificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, _ := util.GetIdFromUrl(r.URL.Path)

		if id == "" {
			stub_handleNotificationGetAllRequest(w)
		} else {
			// Pass the ID provided by the client.
			stub_getSingleNotification(w, id)
		}
	case http.MethodPost:
		stub_registerNotification(w, r)
	case http.MethodDelete:
		id, _ := util.GetIdFromUrl(r.URL.Path)

		if id == "" {
			util.HttpError(w, "method 'DELETE' requires parameter 'ID'", http.StatusUnprocessableEntity)
			return
		} else {
			stub_deleteNotification(w, id)
		}
	default:
		util.HttpError(w, "This method is not supported! Only POST, GET and DELETE are supported", http.StatusMethodNotAllowed)
	}
}

// stub_handleNotificationGetAllRequest returns all notifications to the client.
// If no notifications are found, then an empty array is returned.
func stub_handleNotificationGetAllRequest(w http.ResponseWriter) {
	file, err := os.Open(util.STUB_DATABASE_NOTIFICATIONS) //open file
	if err != nil {
		log.Println(err)
		return
	}
	defer func(file *os.File) { //closes file at the end of function
		err := file.Close()
		if err != nil {
			log.Println("could not close file")
		}
	}(file)

	//adds all notifications to a slice
	var allNotifications []models.NotificationDatabaseModel
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&allNotifications)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Marshall the returning struct from the database
	marshalled, err := json.Marshal(&allNotifications)
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
		log.Printf("failed to return marshalled object. %v", err)
		util.HttpError(w, "could not return marshalled object", http.StatusInternalServerError)
		return
	}
}

// stub_getSingleNotification retrieves a single notification by its id.
// If the notification was not found, 404 not found is returned to the client.
func stub_getSingleNotification(w http.ResponseWriter, id string) {
	file, err := os.Open(util.STUB_DATABASE_NOTIFICATIONS) //open file
	if err != nil {
		log.Println(err)
		return
	}
	defer func(file *os.File) { //closes file at the end of a function
		err := file.Close()
		if err != nil {
			log.Println("could not close file")
		}
	}(file)

	//adds all notifications to a slice
	var allNotifications []models.NotificationDatabaseModel
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&allNotifications)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var notification models.NotificationDatabaseModel

	for _, not := range allNotifications {
		if not.Id == id { //looking for a match
			notification = not
			break
		}
	}

	// Marshall the returning struct from the database
	marshalled, err := json.Marshal(notification)
	if err != nil {
		log.Println("Failed to marshall notification structs in getAllNotifications().")
		log.Println(err)
		util.HttpError(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set(util.CONTENT_TYPE, util.MIMETYPE_JSON)
	w.Header().Set(util.X_CONTENT_TYPE_OPTION, "nosniff")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, string(marshalled))
	if err != nil {
		return
	}
}

// stub_registerNotification lets the client register a new notification.
//
// The client provides us a URL, so we can notify them whenever a certain event has occurred.
func stub_registerNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(util.CONTENT_TYPE, util.MIMETYPE_JSON)

	file, err := os.ReadFile(util.STUB_DATABASE_NOTIFICATIONS) //open file
	if err != nil {
		log.Println(err)
		return
	}

	//adds all notifications to a slice
	var allNotifications []models.NotificationDatabaseModel
	if err := json.Unmarshal(file, &allNotifications); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Expects incoming body in terms of WebhookRegistration struct
	databaseModel := models.NotificationDatabaseModel{}
	err = json.NewDecoder(r.Body).Decode(&databaseModel)

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

	//makes an id for notification
	databaseModel.Id = myCrypto.GetMD5Hash(time.Now().String() + databaseModel.Url)

	// Store the new notification to the database/json file
	allNotifications = append(allNotifications, databaseModel)
	updatedAllNotifications, err := json.MarshalIndent(allNotifications, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	//with the help of ChatGPT. Overwriting the old file
	if err := os.WriteFile(util.STUB_DATABASE_NOTIFICATIONS, updatedAllNotifications, 0666); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("successfully added new notification")

	// Write to client
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "{\"id\": \"%v\"}", databaseModel.Id)
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	} else {
		log.Println("Webhook " + databaseModel.Url + " has been registered.")
	}
}

// stub_deleteNotification removes a notification by id from the database.
func stub_deleteNotification(w http.ResponseWriter, id string) {
	file, err := os.ReadFile(util.STUB_DATABASE_NOTIFICATIONS) //open file
	if err != nil {
		log.Println(err)
		return
	}

	//adds all notifications to a slice
	var allNotifications []models.NotificationDatabaseModel
	if err := json.Unmarshal(file, &allNotifications); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	var allNotificationsUpdated []models.NotificationDatabaseModel
	foundMatch := false
	for _, not := range allNotifications {
		if not.Id != id { //adds all notifications that doesnt have a matching id
			allNotificationsUpdated = append(allNotificationsUpdated, not)
		} else { //if match was found
			foundMatch = true
		}
	}
	if foundMatch == false { //if no match was found
		log.Println("could not find/delete wanted notification")
		return
	}

	updatedNotificationList, err := json.Marshal(allNotificationsUpdated)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Write the updated data back to the file
	if err := os.WriteFile(util.STUB_DATABASE_NOTIFICATIONS, updatedNotificationList, 0666); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	log.Printf("Deleted notification by id %v\n", id)
	http.Error(w, "", http.StatusNoContent)
}
