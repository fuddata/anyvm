package providers

import (
	"context"
	"fmt"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/models"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

type AzureProvider struct {
	client *armcompute.VirtualMachinesClient
}

func NewAzureProvider(cfg *config.Config) (*AzureProvider, bool) {
	cred, err := azidentity.NewClientSecretCredential(cfg.AzureCreds.TenantID, cfg.AzureCreds.ClientID, cfg.AzureCreds.ClientSecret, nil)
	if err != nil {
		fmt.Printf("Failed to active Azure provider. Will continue without it. Error: %v\r\n", err)
		return nil, false
	}
	subscriptionID := cfg.AzureCreds.SubscriptionID
	if subscriptionID == "" {
		panic("Azure subscription ID is not provided")
	}
	client, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		panic(err)
	}
	return &AzureProvider{client: client}, true
}

// GET https://management.azure.com/subscriptions/<subcription id>/providers/Microsoft.Compute/virtualMachines?api-version=2022-03-01
func (p *AzureProvider) ListVMs() ([]models.VM, error) {
	ctx := context.Background()
	var vms []models.VM

	pager := p.client.NewListAllPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, vm := range page.Value {
			vms = append(vms, models.VM{
				ID:       *vm.ID,
				Name:     *vm.Name,
				Provider: "azure",
				Region:   *vm.Location,
				Status:   "running", // Simplified; real status requires additional API call
			})
		}
	}
	return vms, nil
}

// CreateVM creates a new virtual machine in the specified resource group.
// The caller must supply the VM parameters (of type armcompute.VirtualMachine).
func (p *AzureProvider) CreateVM(ctx context.Context, resourceGroupName, vmName string, parameters armcompute.VirtualMachine) error {
	poller, err := p.client.BeginCreateOrUpdate(ctx, resourceGroupName, vmName, parameters, nil)
	if err != nil {
		return fmt.Errorf("failed to start VM creation: %w", err)
	}

	// FixMe: We might need to use bigger value in here?
	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore@v1.17.0/runtime#PollUntilDoneOptions
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create VM: %w", err)
	}
	return nil
}
