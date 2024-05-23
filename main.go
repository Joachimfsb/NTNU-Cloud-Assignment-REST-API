package main

import (
	"assignment2/database"
	"assignment2/handler"
	"assignment2/stubs"
	"assignment2/util"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize build configurations
	err := util.InitializeConfig()
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	} else {
		log.Println("Initialized config")
	}

	// Initialize database
	if util.Config.Stubs.Database == true {
		// Database stub
		go stubs.DatabaseStub()
	} else {
		// Live database
		defer func() {
			err := database.CloseDatabase()
			if err != nil {
				log.Fatal("Unable to close the database connection")
			} else {
				log.Println("Closed database connection")
			}
		}()

		dbError := database.InitializeDatabase()
		if dbError != nil {
			log.Fatalf("Failed to initialize the database connection: %v", dbError)
		} else {
			log.Println("Database initialized")
		}
	}

	// Start stub service if it's environment variable is present.
	if util.Config.Stubs.Weather == true {
		go stubs.Weather_stub()
	}

	if util.Config.Stubs.Currencies == true {
		go stubs.Currency_stub()
	}

	if util.Config.Stubs.RestCountries == true {
		go stubs.Country_stub()
	}

	// Run the main server
	go func() {
		port := os.Getenv("PORT")
		if port == "" { //sets port to 8080 as a default
			log.Println(("$PORT has not been set. Default 8080"))
			port = "8080"
		}
		http.HandleFunc(util.REGISTRATION_PATH, handler.RegistrationHandler)
		http.HandleFunc(util.DASHBOARD_PATH, handler.DashboardHandler)
		http.HandleFunc(util.NOTIFICATION_PATH, handler.NotificationHandler)
		http.HandleFunc(util.STATUS_PATH, handler.StatusHandler)

		log.Println("Service is listening on port: " + port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	// https://stackoverflow.com/a/66834066
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	for {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			fmt.Println("sigint")
			return
		case syscall.SIGTERM:
			fmt.Println("sigterm")
			return
		}
	}
}
