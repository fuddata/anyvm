package providers

import (
	"errors"
	"fmt"
	"os"

	"github.com/fuddata/anyvm/models"
)

type NutanixProvider struct {
	apiURL   string
	username string
	password string
}

func NewNutanixProvider(cfg interface{}) (*NutanixProvider, bool) {
	// Load Nutanix credentials from environment variables.
	apiURL := os.Getenv("NUTANIX_API_URL")
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")

	if apiURL == "" || username == "" || password == "" {
		fmt.Printf("Nutanix credentials not configured. Will continue without it.\r\n")
		return nil, false
	}

	return &NutanixProvider{
		apiURL:   apiURL,
		username: username,
		password: password,
	}, true
}

func (p *NutanixProvider) ListVMs() ([]models.VM, error) {
	// Stub implementation. In production, use Nutanix API to retrieve VMs.
	if p.apiURL == "" || p.username == "" || p.password == "" {
		return nil, errors.New("Nutanix credentials not configured")
	}

	vms := []models.VM{
		{
			ID:       "nut-001",
			Name:     "NutanixVM1",
			Provider: "nutanix",
			Region:   "cluster1",
			Status:   "running",
		},
	}
	return vms, nil
}
