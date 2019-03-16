package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rickschubert/ebay-camel-camel-camel/database"
	"github.com/satori/go.uuid"
)

type trackItemResponse struct {
	Success bool   `json:"success"`
	UUID    string `json:"uuid"`
}

func trackItem(w http.ResponseWriter, r *http.Request) {
	var article database.ItemToTrack
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&article)
	article.UUID = uuid.Must(uuid.NewV4()).String()

	db.CreateTracking(article)

	response := trackItemResponse{
		Success: true,
		UUID:    article.UUID,
	}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(js)
}

func removeTracking(w http.ResponseWriter, r *http.Request) {
	trackingId := mux.Vars(r)["trackingUUID"]
	db.DeleteTracking(trackingId)
	w.WriteHeader(204)
}

var db database.Database

func main() {
	db = database.New()
	router := mux.NewRouter()
	router.HandleFunc("/api/track-item", trackItem).Methods("POST")
	router.HandleFunc("/api/untrack/{trackingUUID}", removeTracking).Methods("DELETE")
	http.ListenAndServe(":8079", router)
}
