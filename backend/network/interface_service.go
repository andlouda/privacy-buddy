package network

import (
)

// NetworkInterfaceService defines the interface for listing network interfaces.
type NetworkInterfaceService interface {
	ListInterfaces() ([]NetworkInterface, error)
}

// NewNetworkInterfaceService creates a new instance of NetworkInterfaceService based on the operating system.
// This function is implemented in platform-specific files (e.g., interface_service_darwin.go, etc.).
func NewNetworkInterfaceService() NetworkInterfaceService {
	return nil // Should never be reached
}