package main

import (
	"encoding/json"
	"fmt"
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

	_, errWriting := db.CreateTracking(article)
	if errWriting != nil {
		panic(fmt.Sprintf("An error occured when trying to post the item to DynamoDB: %v", errWriting.Error()))
	}

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
	_, err := db.DeleteTracking(trackingId)
	if err != nil {
		panic(fmt.Sprintf("An error occured when trying to delete the item from DynamoDB: %v", err.Error()))
	}
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
