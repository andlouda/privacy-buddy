//go:build darwin

package network

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"

	anynetwork "privacy-buddy/backend/network"
)

// DarwinARPCacheService provides macOS-specific implementation for reading the ARP cache.
type DarwinARPCacheService struct{}

// GetARPEntries reads the ARP cache on macOS using gopsutil.
func (s *DarwinARPCacheService) GetARPEntries() ([]anynetwork.ARPEntry, error) {
	arpEntries, err := net.ArpEntries(false) // false for no connections
	if err != nil {
		return nil, fmt.Errorf("failed to get ARP entries: %w", err)
	}

	var entries []anynetwork.ARPEntry
	for _, entry := range arpEntries {
		entries = append(entries, anynetwork.ARPEntry{
			IPAddress:  entry.IPAddress,
			MACAddress: entry.HardwareAddr,
			Interface:  entry.Interface,
			Type:       entry.Type, // gopsutil provides this
		})
	}

	return entries, nil
}


