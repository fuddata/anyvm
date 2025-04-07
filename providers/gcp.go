package providers

import (
	"context"
	"fmt"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/models"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type GCPProvider struct {
	Client    *compute.Service
	projectID string
}

func NewGCPProvider(cfg *config.Config) (*GCPProvider, bool) {
	ctx := context.Background()
	client, err := compute.NewService(ctx, option.WithCredentialsFile(cfg.GCPCreds.CredentialsFile))
	if err != nil {
		fmt.Printf("Failed to active GCP provider. Will continue without it. Error: %v\r\n", err)
		return nil, false
	}
	return &GCPProvider{Client: client, projectID: cfg.GCPCreds.ProjectID}, true
}

// GET  https://compute.googleapis.com/compute/v1/projects/<project id>/aggregated/instances?alt=json&prettyPrint=false
func (p *GCPProvider) ListVMs() ([]models.VM, error) {
	ctx := context.Background()
	var vms []models.VM

	req := p.Client.Instances.AggregatedList(p.projectID)
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
