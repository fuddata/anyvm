package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fuddata/anyvm/models"

	"github.com/masterzen/winrm"
)

// HyperVProvider uses WinRM to remotely execute PowerShell on a Hyper‑V host.
type HyperVProvider struct {
	client *winrm.Client
	host   string
}

// NewHyperVProvider creates and configures a HyperVProvider using environment variables.
func NewHyperVProvider(cfg interface{}) (*HyperVProvider, bool) {
	host := os.Getenv("HYPERV_HOST")
	portStr := os.Getenv("HYPERV_PORT")
	username := os.Getenv("HYPERV_USERNAME")
	password := os.Getenv("HYPERV_PASSWORD")

	if host == "" || username == "" || password == "" {
		fmt.Printf("Hyper-V credentials not configured. Will continue without it.\r\n")
		return nil, false
	}

	port := 5985
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Create the WinRM endpoint.
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	// Set up WinRM parameters with NTLM encryption.
	params := winrm.DefaultParameters
	enc, err := winrm.NewEncryption("ntlm")
	if err != nil {
		panic(fmt.Sprintf("failed to create encryption: %v", err))
	}
	params.TransportDecorator = func() winrm.Transporter { return enc }

	client, err := winrm.NewClientWithParameters(endpoint, username, password, params)
	if err != nil {
		panic(fmt.Sprintf("failed to create winrm client: %v", err))
	}

	return &HyperVProvider{
		client: client,
		host:   host,
	}, true
}

// hypervVM represents the subset of VM information returned by the PowerShell command.
type hypervVM struct {
	Id    interface{} `json:"Id"`
	Name  string      `json:"Name"`
	State string      `json:"State"`
}

// ListVMs runs a PowerShell command via WinRM to retrieve Hyper‑V VMs and parses the output.
func (p *HyperVProvider) ListVMs() ([]models.VM, error) {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run the PowerShell command to list VMs.
	cmd := `Get-WmiObject -Namespace "root\virtualization\v2" -Class "Msvm_ComputerSystem" | Where-Object { $_.Caption -eq "Virtual Machine" } | Select-Object @{l="Id";e={$_.Name.ToLower()}},@{l="Name";e={$_.ElementName}},@{l="State";e={if ($_.ProcessID){"Running"} else {"Stopped"}}} | ConvertTo-Json -Compress`
	stdOut, stdErr, exitCode, err := p.client.RunPSWithContext(ctx, cmd)
	if err != nil || exitCode != 0 {
		return nil, fmt.Errorf("failed to run command: %v, stderr: %s", err, stdErr)
	}

	// Pre-process output: if the output starts with a quote, unquote it.
	trimmed := strings.TrimSpace(stdOut)
	if len(trimmed) > 0 && trimmed[0] == '"' {
		unquoted, err := strconv.Unquote(trimmed)
		if err == nil {
			stdOut = unquoted
		}
	}

	// Parse the JSON output. Handle both array and single object cases.
	var vmsData []hypervVM
	err = json.Unmarshal([]byte(stdOut), &vmsData)
	if err != nil {
		// Try unmarshaling as a single object.
		var single hypervVM
		if err2 := json.Unmarshal([]byte(stdOut), &single); err2 != nil {
			return nil, fmt.Errorf("failed to parse JSON output: %v , std out: %s", err, stdOut)
		}
		vmsData = []hypervVM{single}
	}

	// Convert to the unified VM model.
	var vms []models.VM
	for _, hv := range vmsData {
		idStr := ""
		switch v := hv.Id.(type) {
		case string:
			idStr = v
		case float64:
			idStr = fmt.Sprintf("%.0f", v)
		default:
			idStr = fmt.Sprintf("%v", v)
		}
		vms = append(vms, models.VM{
			ID:       idStr,
			Name:     hv.Name,
			Provider: "hyperv",
			Region:   p.host,
			Status:   hv.State,
		})
	}
	return vms, nil
}
