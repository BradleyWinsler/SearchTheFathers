package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	log.Println("Serving on port 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
