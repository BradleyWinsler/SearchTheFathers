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
	tagHandlers := handlers.NewTagHandlers(mongoClient)

	// Citation routes
	r.HandleFunc("/api/citations", citationHandlers.GetCitations).Methods("GET")
	r.HandleFunc("/api/citations/{id}", citationHandlers.GetCitation).Methods("GET")
	r.HandleFunc("/api/citations", citationHandlers.AddCitation).Methods("POST")
	r.HandleFunc("/api/citations/{id}", citationHandlers.DeleteCitation).Methods("DELETE")

	// Tag routes
	r.HandleFunc("/api/tags", tagHandlers.GetTags).Methods("GET")
	r.HandleFunc("/api/tags/{slug}", tagHandlers.AddTag).Methods("POST")
	r.HandleFunc("/api/tags/{slug}", tagHandlers.DeleteTag).Methods("DELETE")

	log.Println("Serving on port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
