package main

import (
	"net/http"

	controller "middlewareApi/controllers"
	"middlewareApi/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize logging
	utils.InitLogger()
	utils.InitSession()
	defer utils.CloseLogger()

	// Initialize database connection
	utils.InitDBMaster() // Replace with your DB credentials

	mux := mux.NewRouter()
	mux.Use(utils.RequestIdMiddleware)
	mux.Use(utils.LoggingMiddleware)

	// Define the SOAP endpoint
	mux.HandleFunc("/service/authorization", controller.AuthController).Methods("GET")
	mux.HandleFunc("/service/store-procedure", controller.SpController).Methods("POST")
	mux.HandleFunc("/service/safe-query", controller.QueryController).Methods("POST")
	mux.HandleFunc("/service/webhook", controller.WebhookController).Methods("POST")

	// CORS configuration
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-API-KEY", "X-SIGNATURE"}),
	)

	// Start the server
	utils.Info("API Gateway Start | PORT:8081")
	err := http.ListenAndServe(":8081", corsHandler(utils.SessionManager.LoadAndSave(mux)))
	if err != nil {
		utils.Fatal("Server failed to start: %v", err)
	}
}
