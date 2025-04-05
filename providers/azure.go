package providers

import (
	"context"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/models"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

type AzureProvider struct {
	client *armcompute.VirtualMachinesClient
}

func NewAzureProvider(cfg *config.Config) *AzureProvider {
	cred, err := azidentity.NewClientSecretCredential(cfg.AzureCreds.TenantID, cfg.AzureCreds.ClientID, cfg.AzureCreds.ClientSecret, nil)
	if err != nil {
		panic(err) // In production, handle gracefully
	}
	subscriptionID := cfg.AzureCreds.SubscriptionID
	if subscriptionID == "" {
		panic("Azure subscription ID is not provided")
	}
	client, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, nil)
	if err != nil {
		panic(err)
	}
	return &AzureProvider{client: client}
}

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
