package main

import (
	"log"
	"net/http"

	"github.com/alanowatson/LeadGenAPI/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/playlisters", handlers.GetPlaylisters).Methods("GET")
    r.HandleFunc("/playlisters", handlers.CreatePlaylister).Methods("POST")
    r.HandleFunc("/playlisters/{id}", handlers.GetPlaylister).Methods("GET")
    r.HandleFunc("/playlisters/{id}", handlers.UpdatePlaylister).Methods("PUT")
    r.HandleFunc("/playlisters/{id}", handlers.DeletePlaylister).Methods("DELETE")

    log.Println("Starting server on :8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}
