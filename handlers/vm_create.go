package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fuddata/anyvm/models"
	"github.com/fuddata/anyvm/providers"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

// CreateVMRequest defines the request payload for creating a VM.
type CreateVMRequest struct {
	Provider          string `json:"provider"`
	ResourceGroupName string `json:"resourceGroupName"`
	VMName            string `json:"vmName"`
	Location          string `json:"location"`
	VMSize            string `json:"vmSize"`
	AdminUsername     string `json:"adminUsername"`
	AdminPassword     string `json:"adminPassword"`
	NICID             string `json:"nicId"`
}

// CreateVMHandler handles VM creation requests.
// It only supports Azure; for other providers it returns an error.
func CreateVMHandler(cm *providers.CloudManager) http.HandlerFunc {
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

		// Only Azure is supported for VM creation.
		if strings.ToLower(req.Provider) != "azure" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: false,
				Error:   "VM creation is only supported for Azure",
			})
			return
		}

		// Retrieve the Azure provider from the CloudManager.
		provider := cm.GetProvider("azure")
		if provider == nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: false,
				Error:   "Azure provider not available",
			})
			return
		}

		azureProvider, ok := provider.(*providers.AzureProvider)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: false,
				Error:   "Invalid Azure provider instance",
			})
			return
		}

		// FixMe: Read these from somewhere
		// az vm image list-offers -l westeurope -p Canonical
		// az vm image list-skus -l westeurope -p Canonical -f ubuntu-24_04-lts
		publisher := "Canonical"
		offer := "ubuntu-24_04-lts"
		sku := "server"
		version := "latest"

		// Build the VM parameters.
		vmSize := armcompute.VirtualMachineSizeTypes(req.VMSize)
		createOption := armcompute.DiskCreateOptionTypesFromImage
		vmParameters := armcompute.VirtualMachine{
			Location: &req.Location,
			Properties: &armcompute.VirtualMachineProperties{
				HardwareProfile: &armcompute.HardwareProfile{
					VMSize: &vmSize,
				},
				StorageProfile: &armcompute.StorageProfile{
					ImageReference: &armcompute.ImageReference{
						Publisher: &publisher,
						Offer:     &offer,
						SKU:       &sku,
						Version:   &version,
					},
					OSDisk: &armcompute.OSDisk{
						CreateOption: &createOption,
						// DiskSizeGB:   to.Int32Ptr(30),
					},
				},
				OSProfile: &armcompute.OSProfile{
					ComputerName:  &req.VMName,
					AdminUsername: &req.AdminUsername,
					AdminPassword: &req.AdminPassword,
				},
				NetworkProfile: &armcompute.NetworkProfile{
					NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
						{
							ID: &req.NICID,
							/*
								Properties: &armcompute.NetworkInterfaceReferenceProperties{
									Primary: to.BoolPtr(true),
								},
							*/
						},
					},
				},
			},
		}

		// Call the Azure provider's CreateVM function.
		if err := azureProvider.CreateVM(r.Context(), req.ResourceGroupName, req.VMName, vmParameters); err != nil {
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
