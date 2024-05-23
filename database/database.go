// database package is the our service's internal Database API. It is an abstraction layer to database services.
// The Database API contributes to separation of concern in our service, allowing it to focus on serving clients with data without
// caring about where the data originates from. This means that the service can ask for whatever data is needs, and the Database
// API will make sure to get data from the correct sources.
//
// Currently, the Database API supports retrieving data from Google Firestore and local stub services. The stub services also
// ensure testability for our service without connecting to external services.
package database

import (
	"assignment2/util"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
)

var Ctx context.Context
var Client *firestore.Client

// InitializeDatabase is a public API that must be called, but only once.
// This function ensures that all connections to external databases are handled correctly.
// If the Database API is used without initializing databases, there will be nil pointer errors.
// This is not handled, as it is considered a developer error.
//
// Returns:
// If an error occurs, an error object is returned.
func InitializeDatabase() error {
	return initializeFirebase()
}

// initializeFirebase is an internal function that initializes the service's connection to Google Firestore.
// It ensures that the credentials key is sent to Firestore and verified.
//
// Returns:
// An error object is returned if the connection was unsuccessful. It might be because the key does not exist,
// it is not provided in config.yaml, or because something in the key is incorrect. There is also a chance of network
// errors.
func initializeFirebase() error {
	Ctx = context.Background() //initializes Firebase

	// Find key from environment variable
	firebaseKey := util.Config.Secrets.FirebaseKey
	if firebaseKey == "" {
		return fmt.Errorf("firebaseKey has not been set")
	}

	sa := option.WithCredentialsFile(firebaseKey) //path to key
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return err
	}

	Client, err = app.Firestore(Ctx)
	if err != nil {
		return err
	}

	return nil
}

// CloseDatabase closes the connection to databases. This should only be called once.
func CloseDatabase() error {
	// Don't close the client if it's nil... If it's nil, then the database was never initialized anyway.
	if Client == nil {
		return nil
	}

	err := Client.Close()
	if err != nil {
		return err
	} else {
		Client = nil
		return nil
	}
}
