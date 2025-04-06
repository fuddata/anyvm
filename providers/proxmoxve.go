package providers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/fuddata/anyvm/models"
)

type ProxmoxVEProvider struct {
	client *proxmox.Client
	node   string
}

func NewProxmoxVEProvider(cfg interface{}) (*ProxmoxVEProvider, bool) {
	apiURL := os.Getenv("PROXMOX_API_URL")
	username := os.Getenv("PROXMOX_USERNAME")
	password := os.Getenv("PROXMOX_PASSWORD")
	node := os.Getenv("PROXMOX_NODE") // This must be specified

	if apiURL == "" || username == "" || password == "" || node == "" {
		fmt.Printf("ProxmoxVE credentials not configured. Will continue without it.\r\n")
		return nil, false
	}

	// Create an HTTP client for use by the Proxmox client.
	httpClient := &http.Client{}

	// For Telmate's NewClient, we need: (apiURL, *http.Client, realm, *tls.Config, ticket, port)
	// Use an empty realm and ticket. Typical port for Proxmox is 8006.
	client, err := proxmox.NewClient(apiURL, httpClient, "", nil, "", 30)
	if err != nil {
		panic(err)

	}

	// Login using a context. The fourth parameter (OTP) is empty.
	err = client.Login(context.Background(), username, password, "")
	if err != nil {
		fmt.Printf("ProxmoxVE login failed. Will continue without it. Error: %v\r\n", err)
		return nil, false
	}

	return &ProxmoxVEProvider{
		client: client,
		node:   node,
	}, true
}

func (p *ProxmoxVEProvider) ListVMs() ([]models.VM, error) {
	// Use ListGuests to get the list of VMs for the specified node.
	guests, err := proxmox.ListGuests(context.Background(), p.client)
	if err != nil {
		return nil, err
	}

	var vms []models.VM
	for _, guest := range guests {
		vmID := strconv.FormatUint(uint64(guest.Id), 10)
		vms = append(vms, models.VM{
			ID:       vmID,
			Name:     guest.Name,
			Provider: "proxmoxve",
			Region:   p.node,
			Status:   guest.Status,
		})
	}
	return vms, nil
}
