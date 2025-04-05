package main

import (
	"log"
	"net/http"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/handlers"
	"github.com/fuddata/anyvm/providers"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize cloud manager
	cm := providers.NewCloudManager()

	cm.RegisterProvider("azure", providers.NewAzureProvider(cfg))
	cm.RegisterProvider("aws", providers.NewAWSProvider(cfg))
	cm.RegisterProvider("gcp", providers.NewGCPProvider(cfg))

	// Set up router
	r := mux.NewRouter()

	// API routes with auth middleware
	api := r.PathPrefix("/api/v1").Subrouter()

	// FixMe: Enable authentication
	// api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/vms", handlers.ListVMsHandler(cm)).Methods("GET")

	// Start server
	log.Printf("Server starting on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
