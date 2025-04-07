package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/providers"
	"google.golang.org/api/compute/v1"
)

// Helper function for GCP VM creation.
func createGCPVM(ctx context.Context, req CreateVMRequest, cm *providers.CloudManager, cfg *config.Config) error {
	prov := cm.GetProvider("gcp")
	if prov == nil {
		return fmt.Errorf("GCP provider not available")
	}
	gcpProvider, ok := prov.(*providers.GCPProvider)
	if !ok {
		return fmt.Errorf("invalid GCP provider instance")
	}
	actualMachineType := req.MachineType
	if mapped, ok := cfg.Mappings.GCP.CustomVMSizes[strings.ToLower(req.MachineType)]; ok {
		actualMachineType = mapped
	}
	actualSourceImage := req.SourceImage
	if mapped, ok := cfg.Mappings.GCP.CustomImages[strings.ToLower(req.SourceImage)]; ok {
		actualSourceImage = mapped
	}
	zone := req.Zone
	if zone == "" {
		zone = cfg.Mappings.GCP.DefaultZone
	}
	projectID := req.ProjectID
	if projectID == "" {
		projectID = cfg.Mappings.GCP.DefaultProject
	}
	instance := &compute.Instance{
		Name:        req.VMName,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", zone, actualMachineType),
		Disks: []*compute.AttachedDisk{
			{
				Boot: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: actualSourceImage,
					DiskSizeGb:  10,
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := gcpProvider.Client.Instances.Insert(projectID, zone, instance).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create GCP instance: %w", err)
	}
	_ = op
	return nil
}
