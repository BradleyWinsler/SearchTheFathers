package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/BradleyWinsler/SearchTheFathers/handlers"
	"github.com/BradleyWinsler/SearchTheFathers/store"
)

const (
	mongoURI = "mongodb://fathers_mongo:27017/fathers"
)

func main() {
	r := mux.NewRouter()

	// Mongo setup
	mongoClient, err := store.NewClient(context.Background(), mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = store.Close(mongoClient)
	}()

	// Setup the handlers with the mongo client
	citationHandlers := handlers.NewCitationHandlers(mongoClient)

	// Routes
	r.HandleFunc("/api/citations", citationHandlers.GetCitations).Methods("GET")

	log.Println("Serving on port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
