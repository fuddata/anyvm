package config

import "os"

type Config struct {
	Port       string
	JWTSecret  string
	AzureCreds AzureCredentials
	AWSCreds   AWSCredentials
	GCPCreds   GCPCredentials
}

type AzureCredentials struct {
	TenantID     string
	ClientID     string
	ClientSecret string
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

func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		AzureCreds: AzureCredentials{
			TenantID:     getEnv("AZURE_TENANT_ID", ""),
			ClientID:     getEnv("AZURE_CLIENT_ID", ""),
			ClientSecret: getEnv("AZURE_CLIENT_SECRET", ""),
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
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
