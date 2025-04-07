package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/providers"
)

// Helper function for Azure VM creation.
func createAzureVM(ctx context.Context, req CreateVMRequest, cm *providers.CloudManager, cfg *config.Config) error {
	prov := cm.GetProvider("azure")
	if prov == nil {
		return fmt.Errorf("Azure provider not available")
	}
	azureProvider, ok := prov.(*providers.AzureProvider)
	if !ok {
		return fmt.Errorf("invalid Azure provider instance")
	}

	// Map custom VM size.
	actualSize := req.VMSize
	if mapped, ok := cfg.Mappings.Azure.CustomVMSizes[strings.ToLower(req.VMSize)]; ok {
		actualSize = mapped
	}

	// Use defaults if resource group or location are not provided.
	resourceGroup := req.ResourceGroupName
	if resourceGroup == "" {
		resourceGroup = cfg.Mappings.Azure.DefaultResourceGroup
	}
	location := req.Location
	if location == "" {
		location = cfg.Mappings.Azure.DefaultLocation
	}
	// Use a default image key ("ubuntu18") for this example.
	imageKey := "ubuntu18"
	actualImageRef := cfg.Mappings.Azure.CustomImages[imageKey]

	vmSize := armcompute.VirtualMachineSizeTypes(actualSize)
	createOption := armcompute.DiskCreateOptionTypesFromImage
	vmParameters := armcompute.VirtualMachine{
		Location: &location,
		Properties: &armcompute.VirtualMachineProperties{
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: &vmSize,
			},
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: parseAzureImageReference(actualImageRef),
				OSDisk: &armcompute.OSDisk{
					CreateOption: &createOption,
					DiskSizeGB:   to.Int32Ptr(30),
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
						Properties: &armcompute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
		},
	}

	// Use a timeout for creation.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	return azureProvider.CreateVM(ctx, resourceGroup, req.VMName, vmParameters)
}

// Helper function to parse an Azure image reference string in the format "Publisher:Offer:SKU:Version".
func parseAzureImageReference(ref string) *armcompute.ImageReference {
	parts := strings.Split(ref, ":")
	if len(parts) != 4 {
		// Fallback defaults.
		return &armcompute.ImageReference{
			Publisher: to.StringPtr("Canonical"),
			Offer:     to.StringPtr("UbuntuServer"),
			SKU:       to.StringPtr("18.04-LTS"),
			Version:   to.StringPtr("latest"),
		}
	}
	return &armcompute.ImageReference{
		Publisher: to.StringPtr(parts[0]),
		Offer:     to.StringPtr(parts[1]),
		SKU:       to.StringPtr(parts[2]),
		Version:   to.StringPtr(parts[3]),
	}
}
