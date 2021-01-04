package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"

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

	log.Println("Serving on port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
