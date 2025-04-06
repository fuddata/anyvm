package providers

import (
	"fmt"

	"github.com/fuddata/anyvm/config"
	"github.com/fuddata/anyvm/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSProvider struct {
	client *ec2.EC2
}

func NewAWSProvider(cfg *config.Config) (*AWSProvider, bool) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSCreds.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AWSCreds.AccessKey, cfg.AWSCreds.SecretKey, ""),
	})
	if err != nil {
		fmt.Printf("Failed to active AWS provider. Will continue without it. Error: %v\r\n", err)
		return nil, false
	}
	return &AWSProvider{client: ec2.New(sess)}, true
}

// POST https://ec2.eu-west-3.amazonaws.com
// Action=DescribeInstances&Version=2016-11-15
func (p *AWSProvider) ListVMs() ([]models.VM, error) {
	result, err := p.client.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}
	var vms []models.VM
	for _, res := range result.Reservations {
		for _, inst := range res.Instances {
			vms = append(vms, models.VM{
				ID:       *inst.InstanceId,
				Name:     getTagValue(inst.Tags, "Name"),
				Provider: "aws",
				Region:   *inst.Placement.AvailabilityZone,
				Status:   *inst.State.Name,
			})
		}
	}
	return vms, nil
}

func getTagValue(tags []*ec2.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}
