package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type articleToTrack struct {
	SearchTerm string `json:"searchTerm"`
	Price      int    `json:"price"`
}

func postyMcPost(w http.ResponseWriter, r *http.Request) {
	var article articleToTrack
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&article)

	fmt.Println(article)

	js, err := json.Marshal(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(js)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/track-product", postyMcPost).Methods("POST")
	http.ListenAndServe(":8079", router)
}
