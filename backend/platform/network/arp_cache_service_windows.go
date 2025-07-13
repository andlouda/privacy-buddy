//go:build windows

package network

import (
	"fmt"

	"github.com/mostlygeek/arp"

	anynetwork "privacy-buddy/backend/network"
)

// WindowsARPCacheService provides Windows-specific implementation for reading the ARP cache.
type WindowsARPCacheService struct{}

// GetARPEntries reads the ARP cache on Windows using gopsutil.
func (s *WindowsARPCacheService) GetARPEntries() ([]anynetwork.ARPEntry, error) {
	arpTable := arp.Table()

	if arpTable == nil {
		return nil, fmt.Errorf("failed to get ARP entries: arp.Table() returned nil")
	}

	var entries []anynetwork.ARPEntry
	for ip, mac := range arpTable {
		entries = append(entries, anynetwork.ARPEntry{
			IPAddress:  ip,
			MACAddress: mac,
			// Interface and Type are not directly available from mostlygeek/arp
			Interface:  "", 
			Type:       "", 
		})
	}

	return entries, nil
}


