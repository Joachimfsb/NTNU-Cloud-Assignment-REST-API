package stubs

import (
	stubs "assignment2/stubs/handler"
	"assignment2/util"
	"log"
	"net/http"
)

// DatabaseStub is file based database system. It replaces the live database and is designed for offline local
// development and testing.
//
// The way it works data is persistently stored in a JSON file. Each JSON object is a struct on its own and represents
// the types that the main server expects. For example: If the server requests a util.Registration datastructure from
// the live server, then it should work the exact way on the stub and return the exact same data.
//
// Credit:
// Code for http.ServeMux: https://stackoverflow.com/a/60232262.
func DatabaseStub() {
	// New server mux so the handlers only work for this database stub
	dbMux := http.NewServeMux()
	dbMux.HandleFunc(util.REGISTRATION_PATH, stubs.DatabaseDashboardHandler)
	dbMux.HandleFunc(util.NOTIFICATION_PATH, stubs.DatabaseNotificationHandler)

	log.Println("Database Stub Service is listening on port: " + util.DATABASE_PORT)
	log.Fatal(http.ListenAndServe(":"+util.DATABASE_PORT, dbMux))
}
