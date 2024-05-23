package database

import (
	"assignment2/util"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"net/http"
	"os"
)

// GetAllRegistrations is a helper function used get all registrations from databases.
//
// Returns:
// - an array of util.Registration objects is returned. This array may be empty.
// - an error object is returned if an error occurred. The callee should check if the error object is not nil.
func GetAllRegistrations() ([]util.Registration, error) {
	var out []util.Registration
	if util.Config.Stubs.Database == false {
		//Gets all documents currently in the "dashboards" collection.
		fireDocs, err := Client.Collection("dashboards").Documents(context.Background()).GetAll()
		if err != nil {
			return []util.Registration{}, fmt.Errorf("Error, could not retrieve data")
		}

		//Splits then iterates over each document in the firestore
		dashboards := make([]util.Registration, 0)
		for _, fireDoc := range fireDocs { //Goes over all documents
			//Creates a structs from the documents
			var registration util.Registration
			//Gives data from firestore to struct
			if err := fireDoc.DataTo(&registration); err != nil {
				log.Println("Failed to decode document data from firebase")
				return []util.Registration{}, fmt.Errorf("Error, unable to process requests")
			}
			docFields := fireDoc.Data()
			registration.ID = docFields["ID"].(string)
			dashboards = append(dashboards, registration) //Appends firestore struct with data to dashboards
			out = dashboards
		}
	} else {
		client := http.Client{}

		// Retrieve content from server
		res, err := client.Get("http://localhost:" + util.DatabaseStubPort + util.REGISTRATION_PATH)
		if res == nil {
			return []util.Registration{}, fmt.Errorf("res is nil %v", err)
		}
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				log.Println("Failed to close repsonse body")
			}
		}(res.Body)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&out); err != nil {
			return nil, fmt.Errorf(
				"The stub database has no response, or the response cannot be decoded to a registration struct.\n%v\n",
				err)
		}
	}

	return out, nil
}

// GetSingleRegistrationByID returns a copy of util.Registration from the database.
// If the Registration object's id exists, then a populated object is returned. If not, an error occurs.
// This error is STRICTLY only for internal use.
//
// Parameters:
// registrationId: the Registration object's ID to search for.
//
// Returns:
// util.Registration: the populated registration object
// error: any errors which may occur
func GetSingleRegistrationByID(registrationId string) (util.Registration, error) {
	var registration util.Registration

	if util.Config.Stubs.Database == false {
		// Made an iterator for going throgh all documents in dashboards collection in database
		iter := Client.Collection(util.DASHBOARDS).Documents(Ctx)

		var fireDoc *firestore.DocumentSnapshot

		// Find document by the Registration object's registration_id
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return util.Registration{}, err
			}

			foundId := doc.Data()["ID"].(string) //gives all document data to docFields

			// Document was found
			if registrationId == foundId {
				fireDoc = doc
				break
			}
		}

		// Document was not found by the Registration object's registration_id
		if fireDoc == nil {
			return util.Registration{}, errors.New("registration not found by registration_id")
		}

		//Converts firestore document into Registration object

		err := fireDoc.DataTo(&registration)
		if err != nil {
			return util.Registration{}, errors.New("could not convert registration data from database to our internal Registration struct")
		}
	} else {
		// Open the JSON file
		file, err := os.Open(util.STUB_DATABASE_REGISTRATIONS)
		if err != nil {
			fmt.Println("Error:", err)
			return util.Registration{}, errors.New("could not find file")
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println("could not close file")
			}
		}(file)

		var allRegistrations []util.Registration
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&allRegistrations)
		if err != nil {
			fmt.Println("Error:", err)
			return util.Registration{}, errors.New("could not decode")
		}

		foundWanted := false
		for _, reg := range allRegistrations {
			if reg.ID == registrationId {
				registration = reg
				foundWanted = true
				break
			}
		}
		if foundWanted != true {
			fmt.Println("could not find a match")
			return util.Registration{}, errors.New("could not find a match")
		}
	}
	return registration, nil
}

// DeleteDashboardById deletes a dashboard from the database by the id parameter.
//
// Parameters:
// - id: dashboard's id to delete by
//
// Return:
// An error object is returned if something happened during the deletion of the dashboard.
// The callee should check whether the an error occurred.
func DeleteDashboardById(id string) error {
	if util.Config.Stubs.Database == false {
		// Search for documents with the specific field value
		matchingQuery := searchForQueryWithID(id)

		// Made an iterator for going throgh all documents that match
		iter := matchingQuery.Documents(Ctx)

		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return err
			}

			// Deletes the document
			_, err = doc.Ref.Delete(Ctx)
			if err == nil {
				return nil
			}
		}
	} else {
		client := http.Client{}

		// Retrieve content from server
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:"+util.DatabaseStubPort+util.REGISTRATION_PATH+id, nil)

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

		return nil

	}
	return fmt.Errorf("unable to delete notification by id: %v", id)
}

// UpdateRegistration updates an existing util.Registration object in the database.
//
// Parameters:
// - registration: The registration object to update
//
// Returns:
// An error object is returned if something went wrong.
func UpdateRegistration(registration util.Registration) error {
	if util.Config.Stubs.Database == false {
		iter := Client.Collection(util.DASHBOARDS).Documents(Ctx)

		var fireDocs *firestore.DocumentSnapshot

		// Find document by the Registration object's registration_id
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return err
			}

			foundId := doc.Data()["ID"].(string) //gives all document data to docFields

			// Document was found
			if registration.ID == foundId {
				fireDocs = doc
				break
			}
		}

		if fireDocs == nil {
			return errors.New("registration not found by id")
		}

		//Saves ID, so it doesn't get overwritten by .Set
		saveID := fireDocs.Data()
		registration.ID = saveID["ID"].(string)

		//Update firestore document with newly inputed data
		_, err := fireDocs.Ref.Set(Ctx, registration)
		if err != nil {
			return errors.New("unable to update registration in the database")
		}
	} else {
		file, err := os.ReadFile(util.STUB_DATABASE_REGISTRATIONS)
		if err != nil {
			return fmt.Errorf("unable to read database file")
		}

		var allRegistrations []util.Registration
		if err := json.Unmarshal(file, &allRegistrations); err != nil {
			return fmt.Errorf("error unmarshalling data")
		}

		//CHATGPT: Check if the ID exists in the array
		found := false
		for i, reg := range allRegistrations {
			if reg.ID == registration.ID {
				found = true
				// Update the registration in place
				allRegistrations[i] = registration
				break
			}
		}

		if !found {
			return errors.New("registration not found by ID")
		}
		client := http.Client{}

		encodedReg, err := json.Marshal(&registration)
		if err != nil {
			return fmt.Errorf("Unable to marshal notiCorrectCasing. This is a developer error.\n%v\n", err)
		}

		req, err := http.NewRequest(http.MethodPut,
			"http://localhost:"+util.DatabaseStubPort+util.REGISTRATION_PATH+registration.ID,
			bytes.NewBuffer(encodedReg))
		if err != nil {
			return fmt.Errorf("request not compatible with client.Do %v", err)
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending put request:", err)
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
		return fmt.Errorf("unable to put registration %v", err)
	}

	return nil
}

// PatchDashboardByID updates specific fields of a dashboard, based on specified ID.
// FUnction searches for an existing document to find an ID match, then applies updates
// specified within the patchData map.
//
// Parameters:
// - id: dashboard's id to match & patch
// - patchData: Map with fields to update & their new data
//
// Return:
// - Error object if something went wrong.
func PatchDashboardByID(id string, patchData map[string]interface{}) error {
	//Find firestore document from given ID

	var fireDoc *firestore.DocumentSnapshot

	// Find document by the Registration object's registration_id
	iter := Client.Collection(util.DASHBOARDS).Documents(Ctx)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		foundId := doc.Data()["ID"] //gives all document data to docFields

		// Document was found
		if id == foundId {
			fireDoc = doc
			break
		}
	}

	if fireDoc == nil {
		return errors.New("failed to update the database. Unable to find a registration object by id")
	}

	//Variable for "slice" to update
	var updates []firestore.Update
	//Iterate over the patchData map
	for key, value := range patchData {
		//Specifies path and new value and appends
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}

	//Updates into firestore
	// ID may complain about potential nil pointer dereference. This is already handled
	_, err := fireDoc.Ref.Update(Ctx, updates)
	if err != nil {
		return errors.New("Error: could not update")
	}

	//Successful, no error
	return nil
}

// searchForQueryWithID searches for a document that has a field ID matching to the parameter id.
// This lets us search for a document by the field ID.
//
// Parameters:
// - id: the id to query
//
// Returns:
// firestore.Query object
func searchForQueryWithID(id string) firestore.Query {
	return Client.Collection(util.DASHBOARDS).Where("ID", "==", id)
}

// AddNewDashboard adds a new dashboard configuration to the database.
//
// Parameters:
// registration: the new dashboard configuration to add.
//
// Returns:
// An error object is returned if something wrong happened.
// The callee should check for this error.
func AddNewDashboard(registration util.Registration, id string) error {
	if util.Config.Stubs.Database == false {
		_, _, err := Client.Collection(util.DASHBOARDS).Add(Ctx, registration)
		if err == nil {
			return err
		}
		return fmt.Errorf("unable to add registration %v", err)
	} else {
		client := http.Client{}

		fmt.Println(id)

		encodedReg, err := json.Marshal(&registration)
		if err != nil {
			return fmt.Errorf("Unable to marshal notiCorrectCasing. This is a developer error.\n%v\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "http://localhost:"+util.DatabaseStubPort+util.REGISTRATION_PATH, bytes.NewBuffer(encodedReg))
		if err != nil {
			return fmt.Errorf("request not compatible with client.Do %v", err)
		}

		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending post request:", err)
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
		return fmt.Errorf("unable to add registration %v", err)
	}
}
