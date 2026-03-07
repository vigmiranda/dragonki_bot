package bootstrap

import (
	"strings"
	"sync"
)

var (
	featureFlagResolverMu sync.RWMutex
	featureFlagResolver   func(name string) (enabled bool, found bool)
)

func SetFeatureFlagResolver(resolver func(name string) (enabled bool, found bool)) {
	featureFlagResolverMu.Lock()
	defer featureFlagResolverMu.Unlock()
	featureFlagResolver = resolver
}

func IsFeatureEnabled(name string, defaultValue bool) bool {
	normalizedName := normalizeFeatureName(name)
	if normalizedName == "" {
		return defaultValue
	}

	featureFlagResolverMu.RLock()
	resolver := featureFlagResolver
	featureFlagResolverMu.RUnlock()

	if resolver != nil {
		enabled, found := resolver(normalizedName)
		if found {
			return enabled
		}
	}

	return defaultValue
}

func normalizeFeatureName(name string) string {
	normalized := strings.TrimSpace(name)
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	normalized = strings.ToUpper(normalized)
	return normalized
}
