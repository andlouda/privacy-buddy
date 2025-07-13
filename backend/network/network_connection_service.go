package network

import (
)

// NetworkConnectionService defines the interface for listing network connections.
type NetworkConnectionService interface {
	GetConnections() ([]NetworkConnection, error)
}

// NewNetworkConnectionService creates a new instance of NetworkConnectionService based on the operating system.
// This function is implemented in platform-specific files (e.g., network_connection_service_darwin.go, etc.).
func NewNetworkConnectionService() NetworkConnectionService {
	return nil // Should never be reached
}