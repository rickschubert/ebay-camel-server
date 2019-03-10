package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Route has been accessed")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

func postyMcPost(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	// decoder := json.NewDecoder(r.Body)
	// fmt.Println(decoder)
	fmt.Println(user)

	profile := User{"Alex", "asdadsf", 43564563}

	js, err := json.Marshal(profile)
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
	router.HandleFunc("/articles", ArticlesCategoryHandler)
	router.HandleFunc("/api/track-product", postyMcPost).Methods("POST")
	http.ListenAndServe(":8079", router)
}
