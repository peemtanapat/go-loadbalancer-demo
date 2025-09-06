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
	Instance 	string		`json:"instance,omitempty"`
	Port      	string      `json:"port,omitempty"`
	Timestamp 	time.Time   `json:"timestamp,omitempty"`
	Users     	[]string    `json:"users,omitempty"`
	ServedBy  	string      `json:"servedBy,omitempty"`
	Message   	string      `json:"message,omitempty"`
	User      	any			`json:"user,omitempty"`
	ProcessingTime int64  	`json:"processingTimeMs,omitempty"`
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

	router.HandleFunc("/api/users", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
			case "GET":
				response := Response{
					Users: []string{"Alice", "Bird", "Charlie", "Dan"},
					ServedBy: instanceName,
					Port: port,
				}

				w.Header().Set("Content-type", "application/json")
				json.NewEncoder(w).Encode(response)
			case "POST":
				var user any
				json.NewDecoder(req.Body).Decode(&user)

				response := Response{
					Message: "User created successfully",
					User: user,
					ServedBy: instanceName,
					Port: port,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
		}
	}).Methods("GET", "POST")

	router.HandleFunc("/api/heavy-task", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		time.Sleep(2 * time.Second)

		response := Response{
			Message: "Heavy task completed",
			ProcessingTime: int64(time.Since(startTime).Milliseconds()),
			ServedBy: instanceName,
			Port: port,
		}

		w.Header().Set("Content-Type", "application/json")
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