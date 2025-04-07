package config

import "os"

type AzureCredentials struct {
	TenantID       string
	ClientID       string
	ClientSecret   string
	SubscriptionID string // Added subscription ID for Azure
}

type AWSCredentials struct {
	AccessKey string
	SecretKey string
	Region    string
}

type GCPCredentials struct {
	ProjectID       string
	CredentialsFile string
}

// Add these type definitions below your existing credential types.

type CloudMappings struct {
	Azure AzureMapping `json:"azure"`
	AWS   AWSMapping   `json:"aws"`
	GCP   GCPMapping   `json:"gcp"`
}

type AzureMapping struct {
	CustomVMSizes        map[string]string `json:"customVmSizes"`        // e.g. "small": "Standard_DS1_v2"
	CustomImages         map[string]string `json:"customImages"`         // e.g. "ubuntu24": "Canonical:UbuntuServer:18.04-LTS:latest"
	DefaultResourceGroup string            `json:"defaultResourceGroup"` // fallback resource group name
	DefaultLocation      string            `json:"defaultLocation"`      // fallback location (e.g. "eastus")
}

type AWSMapping struct {
	CustomVMSizes           map[string]string `json:"customVmSizes"` // e.g. "small": "t2.micro"
	CustomImages            map[string]string `json:"customImages"`  // e.g. "ubuntu24": "ami-0abcdef1234567890"
	DefaultKeyName          string            `json:"defaultKeyName"`
	DefaultSecurityGroupIDs []string          `json:"defaultSecurityGroupIds"`
	DefaultRegion           string            `json:"defaultRegion"`
}

type GCPMapping struct {
	CustomVMSizes  map[string]string `json:"customVmSizes"` // e.g. "small": "n1-standard-1"
	CustomImages   map[string]string `json:"customImages"`  // e.g. "ubuntu24": "projects/ubuntu-os-cloud/global/images/family/ubuntu-1804-lts"
	DefaultZone    string            `json:"defaultZone"`
	DefaultProject string            `json:"defaultProject"`
}

// Then add a new field to your Config struct:
type Config struct {
	Port       string
	JWTSecret  string
	AzureCreds AzureCredentials
	AWSCreds   AWSCredentials
	GCPCreds   GCPCredentials
	Mappings   CloudMappings // <--- new field for unified mappings
}

// Finally, update LoadConfig to set default mappings (or load them from environment variables as needed):
func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		AzureCreds: AzureCredentials{
			TenantID:       getEnv("AZURE_TENANT_ID", ""),
			ClientID:       getEnv("AZURE_CLIENT_ID", ""),
			ClientSecret:   getEnv("AZURE_CLIENT_SECRET", ""),
			SubscriptionID: getEnv("AZURE_SUBSCRIPTION_ID", ""),
		},
		AWSCreds: AWSCredentials{
			AccessKey: getEnv("AWS_ACCESS_KEY", ""),
			SecretKey: getEnv("AWS_SECRET_KEY", ""),
			Region:    getEnv("AWS_REGION", "us-east-1"),
		},
		GCPCreds: GCPCredentials{
			ProjectID:       getEnv("GCP_PROJECT_ID", ""),
			CredentialsFile: getEnv("GCP_CREDENTIALS_FILE", ""),
		},
		Mappings: CloudMappings{
			Azure: AzureMapping{
				CustomVMSizes: map[string]string{
					"small":  "Standard_DS1_v2",
					"medium": "Standard_DS2_v2",
					"large":  "Standard_DS3_v2",
				},
				CustomImages: map[string]string{
					"ubuntu24": "Canonical:ubuntu-24_04-lts:server:latest",
				},
				DefaultResourceGroup: getEnv("AZURE_DEFAULT_RESOURCE_GROUP", "script-test"),
				DefaultLocation:      getEnv("AZURE_DEFAULT_LOCATION", "westeurope"),
			},
			AWS: AWSMapping{
				CustomVMSizes: map[string]string{
					"small":  "t2.micro",
					"medium": "t2.small",
					"large":  "t2.medium",
				},
				CustomImages: map[string]string{
					"ubuntu24": "ami-0644165ab979df02d",
				},
				DefaultKeyName:          getEnv("AWS_DEFAULT_KEYNAME", "default-key"),
				DefaultSecurityGroupIDs: []string{getEnv("AWS_DEFAULT_SECURITY_GROUP", "sg-01234567")},
				DefaultRegion:           getEnv("AWS_DEFAULT_REGION", "eu-west-3"),
			},
			GCP: GCPMapping{
				CustomVMSizes: map[string]string{
					"small":  "t2d-standard-1",
					"medium": "t2d-standard-2",
					"large":  "t2d-standard-4",
				},
				CustomImages: map[string]string{
					"ubuntu24": "projects/ubuntu-os-cloud/global/images/ubuntu-2404-noble-amd64-v20250313",
				},
				DefaultZone:    getEnv("GCP_DEFAULT_ZONE", "europe-west9-c"),
				DefaultProject: getEnv("GCP_DEFAULT_PROJECT", ""),
			},
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
