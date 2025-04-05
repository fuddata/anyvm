package providers

import (
	"fmt"

	"github.com/fuddata/anyvm/models"
)

type CloudProvider interface {
	ListVMs() ([]models.VM, error)
}

type CloudManager struct {
	providers map[string]CloudProvider
}

func NewCloudManager() *CloudManager {
	return &CloudManager{
		providers: make(map[string]CloudProvider),
	}
}

func (cm *CloudManager) RegisterProvider(name string, provider CloudProvider) {
	fmt.Printf("Registering provider %s\r\n", name)
	cm.providers[name] = provider
}

func (cm *CloudManager) GetProvider(name string) CloudProvider {
	return cm.providers[name]
}

func (cm *CloudManager) GetAllProviders() map[string]CloudProvider {
	return cm.providers
}
