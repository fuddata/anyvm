package providers

import (
	"context"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/models"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type GCPProvider struct {
	client    *compute.Service
	projectID string
}

func NewGCPProvider(cfg *config.Config) *GCPProvider {
	ctx := context.Background()
	client, err := compute.NewService(ctx, option.WithCredentialsFile(cfg.GCPCreds.CredentialsFile))
	if err != nil {
		panic(err)
	}
	return &GCPProvider{client: client, projectID: cfg.GCPCreds.ProjectID}
}

func (p *GCPProvider) ListVMs() ([]models.VM, error) {
	ctx := context.Background()
	var vms []models.VM

	req := p.client.Instances.AggregatedList(p.projectID)
	if err := req.Pages(ctx, func(page *compute.InstanceAggregatedList) error {
		for _, instances := range page.Items {
			for _, inst := range instances.Instances {
				vms = append(vms, models.VM{
					ID:       inst.Name,
					Name:     inst.Name,
					Provider: "gcp",
					Region:   inst.Zone,
					Status:   inst.Status,
				})
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return vms, nil
}
