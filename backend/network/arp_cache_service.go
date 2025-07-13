package network

import (
)

// ARPCacheService defines the interface for reading the ARP cache.
type ARPCacheService interface {
	GetARPEntries() ([]ARPEntry, error)
}

// NewARPCacheService creates a new instance of ARPCacheService based on the operating system.
// This function is implemented in platform-specific files (e.g., arp_cache_service_darwin.go, etc.).
func NewARPCacheService() ARPCacheService {
	return nil // Should never be reached
}