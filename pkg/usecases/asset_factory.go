package usecases

import (
	"sync"

	"github.com/eislab-cps/mbaigo/pkg/components"
)

// AssetFactory is a function that creates a UnitAsset from configuration
type AssetFactory func(ConfigurableAsset) components.UnitAsset

var (
	assetFactories = make(map[string]AssetFactory)
	factoryMutex   sync.RWMutex
)

// RegisterAssetFactory registers a factory function for a given asset type
func RegisterAssetFactory(assetType string, factory AssetFactory) {
	factoryMutex.Lock()
	defer factoryMutex.Unlock()
	assetFactories[assetType] = factory
}

// GetAssetFactory retrieves a registered factory function
func GetAssetFactory(assetType string) (AssetFactory, bool) {
	factoryMutex.RLock()
	defer factoryMutex.RUnlock()
	factory, ok := assetFactories[assetType]
	return factory, ok
}

// GetRegisteredAssetTypes returns all registered asset type names
func GetRegisteredAssetTypes() []string {
	factoryMutex.RLock()
	defer factoryMutex.RUnlock()

	types := make([]string, 0, len(assetFactories))
	for assetType := range assetFactories {
		types = append(types, assetType)
	}
	return types
}
