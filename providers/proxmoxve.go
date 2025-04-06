package providers

import (
	"errors"
	"fmt"
	"os"

	"github.com/fuddata/anyvm/models"
)

type ProxmoxVEProvider struct {
	apiURL   string
	username string
	password string
}

func NewProxmoxVEProvider(cfg interface{}) (*ProxmoxVEProvider, bool) {
	// For simplicity, credentials are loaded directly from environment variables.
	apiURL := os.Getenv("PROXMOX_API_URL")
	username := os.Getenv("PROXMOX_USERNAME")
	password := os.Getenv("PROXMOX_PASSWORD")

	if apiURL == "" || username == "" || password == "" {
		fmt.Printf("Proxmox credentials not configured. Will continue without it.\r\n")
		return nil, false
	}

	return &ProxmoxVEProvider{
		apiURL:   apiURL,
		username: username,
		password: password,
	}, true
}

func (p *ProxmoxVEProvider) ListVMs() ([]models.VM, error) {
	// Stub implementation. In production, implement API calls to ProxmoxVE.
	if p.apiURL == "" || p.username == "" || p.password == "" {
		return nil, errors.New("Proxmox credentials not configured")
	}

	vms := []models.VM{
		{
			ID:       "pmx-001",
			Name:     "ProxmoxVM1",
			Provider: "proxmoxve",
			Region:   "local",
			Status:   "running",
		},
	}
	return vms, nil
}
