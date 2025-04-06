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

	azureProvider, azureEnable := providers.NewAzureProvider(cfg)
	if azureEnable {
		cm.RegisterProvider("azure", azureProvider)
	}
	awsProvider, awsEnable := providers.NewAWSProvider(cfg)
	if awsEnable {
		cm.RegisterProvider("aws", awsProvider)
	}
	gcpProvider, gcpEnable := providers.NewGCPProvider(cfg)
	if gcpEnable {
		cm.RegisterProvider("gcp", gcpProvider)
	}

	hypervProvider, hypervEnable := providers.NewHyperVProvider(cfg)
	if hypervEnable {
		cm.RegisterProvider("hyperv", hypervProvider)
	}
	nutanixProvider, nutanixEnable := providers.NewNutanixProvider(cfg)
	if nutanixEnable {
		cm.RegisterProvider("nutanix", nutanixProvider)
	}
	proxmoxProvider, proxmoxEnable := providers.NewProxmoxVEProvider(cfg)
	if proxmoxEnable {
		cm.RegisterProvider("proxmox", proxmoxProvider)
	}
	vsphereProvider, vsphereEnable := providers.NewVSphereProvider(cfg)
	if vsphereEnable {
		cm.RegisterProvider("vsphere", vsphereProvider)
	}

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
