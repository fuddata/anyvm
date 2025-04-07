package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/providers"
)

// Helper function for AWS VM creation.
func createAWSVM(ctx context.Context, req CreateVMRequest, cm *providers.CloudManager, cfg *config.Config) error {
	prov := cm.GetProvider("aws")
	if prov == nil {
		return fmt.Errorf("AWS provider not available")
	}
	awsProvider, ok := prov.(*providers.AWSProvider)
	if !ok {
		return fmt.Errorf("invalid AWS provider instance")
	}

	// Supply defaults if not provided.
	if req.ImageID == "" {
		req.ImageID = "ami-0644165ab979df02d"
	}
	if req.InstanceType == "" {
		req.InstanceType = "small"
	}

	// Map custom instance type and image ID.
	actualInstanceType := req.InstanceType
	if mapped, ok := cfg.Mappings.AWS.CustomVMSizes[strings.ToLower(req.InstanceType)]; ok {
		actualInstanceType = mapped
	}
	actualImageID := req.ImageID
	if mapped, ok := cfg.Mappings.AWS.CustomImages[strings.ToLower(req.ImageID)]; ok {
		actualImageID = mapped
	}

	// Validate that we have valid values.
	if actualImageID == "" {
		return fmt.Errorf("no valid image ID provided for AWS")
	}
	if actualInstanceType == "" {
		return fmt.Errorf("no valid instance type provided for AWS")
	}

	if req.KeyName == "" {
		req.KeyName = cfg.Mappings.AWS.DefaultKeyName
	}
	if len(req.SecurityGroupIDs) == 0 {
		req.SecurityGroupIDs = cfg.Mappings.AWS.DefaultSecurityGroupIDs
	}

	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(actualImageID),
		InstanceType: aws.String(actualInstanceType),
		KeyName:      aws.String(req.KeyName),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		//SecurityGroupIds: aws.StringSlice(req.SecurityGroupIDs),
	}

	_, err := awsProvider.Client.RunInstances(input)
	return err
}
