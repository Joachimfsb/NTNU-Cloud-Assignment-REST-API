package handler

import (
	"assignment2/crypto"
	"assignment2/database"
	"assignment2/notifications"
	"assignment2/util"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// RegistrationHandler is the main entry point for the registration endpoint.
// Routes requests based on their HTTP method (GET, POST, PUT, DELETE & PATCH)
// to their respective handlers.
//
// It handles the following:
// - GET: Retrieve a specific ID or if no ID is specified retrieve all stored ID's.
// - POST: Creates a new dashboard.
// - PUT: Updates a dashboard.
// - DELETE: Removes a dashboard.
// - PATCH: Apply only specified updates to a dashboard.
//
// If an unrecognized method is detected, an appropriate error is returned.
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method { //a switch for the supported methods.
	case http.MethodGet:
		HandleRegistrationGetRequest(w, r)
	case http.MethodPost:
		HandleRegistrationPostRequest(w, r)
	case http.MethodPut:
		HandleRegistrationPutRequest(w, r)
	case http.MethodDelete:
		HandleRegistrationDeleteRequest(w, r)
	case http.MethodPatch:
		HandleRegistrationPatchRequest(w, r)
	default: //Error message if GET method is not used
		http.Error(w, "This method is not supported! Only POST, GET, PUT, DELETE and PATCH are supported", http.StatusNotImplemented)
	}
}

// HandleRegistrationPostRequest creates a new dashboard. First it decodes the request body into a
// Registration struct, generates an INT hashed ID, sets the structs last change time to current, and saves
// the dashboard to the database.
func HandleRegistrationPostRequest(w http.ResponseWriter, r *http.Request) {
	//Decode JSON into registration struct
	var registration util.Registration
	err := json.NewDecoder(r.Body).Decode(&registration)
	if err != nil {
		http.Error(w, "Error, could not parse body", http.StatusBadRequest)
		return
	}

	//Generates hashed ID (mashing country name and time.now)
	hashID := myCrypto.GetMD5Hash(registration.Country + time.Now().String())
	registration.ID = hashID //Sets registration.ID (ID in struct) to hashed ID

	//Time when document is created (later changed in PUT function for last changed time).
	registration.LastChange = time.Now().Format("2006-01-02 15:04")

	//Adds a new document to firestore in "dashboards" collection.
	err = database.AddNewDashboard(registration, hashID)
	if err != nil {
		http.Error(w, "Error, could not store information", http.StatusInternalServerError)
		return
	}

	//Response for request, only shows ID and last change/creation timestamp to the user.
	response, err := json.Marshal(map[string]interface{}{
		"id":         registration.ID,
		"lastChange": registration.LastChange,
	})
	if err != nil {
		http.Error(w, "Error could not genereate response", http.StatusInternalServerError)
		return
	}

	// Invoke registration notification
	err = notifications.InvokeEvent(registration.IsoCode, util.EVENT_REGISTER)
	if err != nil {
		log.Println("Error invoking event:", err)
	}

	//Sets header to JSON type and with creation status message and writes bakc id/time response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(response); err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}

// HandleRegistrationGetRequest retrieves either a specified dashboard or ALL dashboard if no ID
// is given. Decodes documents into Registration structs and returns in JSON format.
func HandleRegistrationGetRequest(w http.ResponseWriter, r *http.Request) {
	//Splits URL path at "/"
	id, _ := util.GetIdFromUrl(r.URL.Path)
	if id == "" {
		dashboards, err := database.GetAllRegistrations()

		// There might not be any registered notifications at all...
		if len(dashboards) == 0 {
			util.HttpError(w, "no registered dashboards found", http.StatusNotFound)
			return
		}

		// Error from the database
		if err != nil {
			log.Println(err)
			util.HttpError(w, "failed to get dashboards", http.StatusInternalServerError)
			return
		}

		//Set content type and header with error message
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dashboards); err != nil { // Use the renamed slice here
			http.Error(w, "Error encoding response JSON", http.StatusInternalServerError)
			return
		}
	} else { //Path has a value, get request for specified id

		//Gets specified document by ID from firestore.
		reg, err := database.GetSingleRegistrationByID(id)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(reg); err != nil {
				http.Error(w, "Error encoding response JSON", http.StatusInternalServerError)
				return
			}
		}

		if err != nil {
			http.Error(w, "Error: could not find specified ID", http.StatusNotFound)
			return
		}
	}
}

// HandleRegistrationPutRequest updates a dashboard by given ID. Decodes body into
// Registration struct, updates the document and updates lastChange to current.
func HandleRegistrationPutRequest(w http.ResponseWriter, r *http.Request) {
	//Get ID from URL
	id, _ := util.GetIdFromUrl(r.URL.Path)
	if id == "" {
		http.Error(w, "Error, no ID declared", http.StatusBadRequest)
		return
	}

	//Decode the request body into a Registration struct (excluding ID and lastChange)
	var registration util.Registration
	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(w, "Error, could not decode body", http.StatusBadRequest)
		return
	}

	registration.ID = id

	//Change the timestamp to last changed
	registration.LastChange = time.Now().Format("2006-01-02 15:04")

	// Update registration on the database
	// Update registration on the database
	err := database.UpdateRegistration(registration)
	if err != nil {
		if err.Error() == "registration not found by ID" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	// Invoke notification event for changing a registration
	if err := notifications.InvokeEvent(registration.IsoCode, util.EVENT_CHANGE); err != nil {
		log.Println("Error invoking event:", err)
	}

}

// HandleRegistrationDeleteRequest deletes a specified document by given ID. Is idempotent
// so may return 204 regardless of whether a document was actually deleted or not.
func HandleRegistrationDeleteRequest(w http.ResponseWriter, r *http.Request) {
	//Get ID from URL
	id, _ := util.GetIdFromUrl(r.URL.Path)
	if id == "" {
		http.Error(w, "Error, no ID declared", http.StatusBadRequest)
		return
	}

	//Get ISO code from registration
	existingRegistration, err := database.GetSingleRegistrationByID(id)
	if err != nil {
		http.Error(w, "Error: could not find specified ID to patch", http.StatusNoContent)
		return
	}

	//Apply ISO so it can be invoked
	isoCode := existingRegistration.IsoCode
	if err := notifications.InvokeEvent(isoCode, util.EVENT_DELETE); err != nil {
		log.Println("Error invoking event:", err)
	}

	//Delete specified document from firestore
	err = database.DeleteDashboardById(id)
	if err != nil {
		log.Println(err)
	}

	//204 No content status
	w.WriteHeader(http.StatusNoContent)
}

// HandleRegistrationPatchRequest applies partial updates to document by given ID.
// Uses a map to only update the fields given in request body (unlike PUT where you give all fields).
func HandleRegistrationPatchRequest(w http.ResponseWriter, r *http.Request) {
	//Get ID from URL
	id, _ := util.GetIdFromUrl(r.URL.Path)
	if id == "" {
		http.Error(w, "Error, no ID declared", http.StatusBadRequest)
		return
	}

	//Map to hold new data from body
	var patchData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		http.Error(w, "Error, could not decode body", http.StatusBadRequest)
		return
	}

	//Get ISO code from registration
	existingRegistration, err := database.GetSingleRegistrationByID(id)
	if err != nil {
		http.Error(w, "Error: could not find specified ID to patch", http.StatusNotFound)
		return
	}

	//Apply ISO so it can be invoked
	isoCode := existingRegistration.IsoCode
	if err := notifications.InvokeEvent(isoCode, util.EVENT_CHANGE); err != nil {
		log.Println("Error invoking event:", err)
	}

	//Update specified dashboard in firestore
	if err := database.PatchDashboardByID(id, patchData); err != nil {
		http.Error(w, "Error, could not patch", http.StatusInternalServerError)
		return
	}
}
