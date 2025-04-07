package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/models"
	"github.com/fuddata/anyvm/providers"
	// Azure SDK helpers
	// AWS SDK
	// GCP SDK
)

// CreateVMRequest defines the unified request payload for creating a VM.
type CreateVMRequest struct {
	Provider string `json:"provider"`
	VMName   string `json:"vmName"`

	// Azure-specific fields
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	Location          string `json:"location,omitempty"`
	VMSize            string `json:"vmSize,omitempty"`
	AdminUsername     string `json:"adminUsername,omitempty"`
	AdminPassword     string `json:"adminPassword,omitempty"`
	NICID             string `json:"nicId,omitempty"`
	// (Optionally, you can allow specifying an image key)

	// AWS-specific fields
	ImageID          string   `json:"imageId,omitempty"`
	InstanceType     string   `json:"instanceType,omitempty"`
	KeyName          string   `json:"keyName,omitempty"`
	SecurityGroupIDs []string `json:"securityGroupIds,omitempty"`

	// GCP-specific fields
	ProjectID   string `json:"projectId,omitempty"`
	Zone        string `json:"zone,omitempty"`
	MachineType string `json:"machineType,omitempty"`
	SourceImage string `json:"sourceImage,omitempty"`
}

// CreateVMHandler handles VM creation requests for Azure, AWS, and GCP.
// It uses unified mappings from the configuration to convert custom identifiers
// to the actual cloud-specific values.
func CreateVMHandler(cm *providers.CloudManager, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req CreateVMRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: false,
				Error:   "Invalid request payload",
			})
			return
		}

		provider := strings.ToLower(req.Provider)
		ctx := r.Context()
		var err error

		switch provider {
		case "azure":
			err = createAzureVM(ctx, req, cm, cfg)
		case "aws":
			err = createAWSVM(ctx, req, cm, cfg)
		case "gcp":
			err = createGCPVM(ctx, req, cm, cfg)
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: false,
				Error:   "VM creation is only supported for Azure, AWS, and GCP",
			})
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(models.APIResponse{
			Success: true,
			Data:    "VM creation initiated successfully",
		})
	}
}
