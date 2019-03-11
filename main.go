package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

const awsDatabaseRegion = "eu-west-2"

type itemToTrack struct {
	SearchTerm string `json:"searchTerm"`
	Price      int    `json:"price"`
	UUID       string
}

type trackItemResponse struct {
	Success bool   `json:"success"`
	UUID    string `json:"uuid"`
}

func trackItem(w http.ResponseWriter, r *http.Request) {
	var article itemToTrack
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&article)
	article.UUID = uuid.Must(uuid.NewV4()).String()

	av, err := dynamodbattribute.MarshalMap(article)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("trackings"),
	}
	_, err = dynamoClient.PutItem(input)
	if err != nil {
		fmt.Println("An error occured when trying to post the item to DynamoDB:")
		fmt.Println(err.Error())
		return
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
	uuidToRemove := mux.Vars(r)["trackingUUID"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UUID": {
				S: aws.String(uuidToRemove),
			},
		},
		TableName: aws.String("trackings"),
	}
	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		fmt.Println("An error occured when trying to delete the item from DynamoDB:")
		fmt.Println(err.Error())
		return
	}
}

var dynamoClient *dynamodb.DynamoDB

func connectToDynamoDB() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsDatabaseRegion)},
	)
	if err != nil {
		panic("Could not initiate new session with Dynamo DB.")
	}
	dynamoClient = dynamodb.New(sess)
}

func main() {
	connectToDynamoDB()
	router := mux.NewRouter()
	router.HandleFunc("/api/track-item", trackItem).Methods("POST")
	router.HandleFunc("/api/untrack/{trackingUUID}", removeTracking).Methods("DELETE")
	http.ListenAndServe(":8079", router)
}
