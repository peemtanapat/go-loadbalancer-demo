package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Status		string		`json:"status,omitempty"`
	Instance 	string		`json:"instance,omitempty`
	Port      	string      `json:"port,omitempty"`
	Timestamp 	time.Time   `json:"timestamp,omitempty"`
	Users     	[]string    `json:"users,omitempty"`
	ServedBy  	string      `json:"servedBy,omitempty"`
	Message   	string      `json:"message,omitempty"`
	User      	interface{} `json:"user,omitempty"`
	ProcessingTime int64  	`json:"processingTime,omitempty"`
}

func main() {
	port := getEnv("PORT", "8080")
	instanceName := getEnv("INSTANCE_NAME", "Go-Unknown-Name")

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		response := Response{
			Status: "healthy",
			Instance: instanceName,
			Port: port,
			Timestamp: time.Now(),
		}

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	fmt.Printf("🚀 API Service (%s) starting on port %s\n", instanceName, port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}