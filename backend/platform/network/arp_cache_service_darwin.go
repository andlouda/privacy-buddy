//go:build darwin

package network

import (
	"os/exec"
	anynetwork "privacy-buddy/backend/network"
	"strings"
)

// DarwinARPCacheService provides macOS-specific implementation for reading the ARP cache.
type DarwinARPCacheService struct{}

// GetARPEntries reads the ARP cache on macOS using gopsutil.
// ArpEntries gibt strukturierte ARP-Einträge für macOS zurück
func ArpEntries() ([]anynetwork.ARPEntry, error) {
	out, err := exec.Command("arp", "-a").Output()
	if err != nil {
		return nil, err
	}

	var entries []anynetwork.ARPEntry
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		// rudimentäres Beispiel: "(192.168.0.1) at 4c:66:41:xx:xx:xx on en0 ifscope [ethernet]"
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			ip := strings.Trim(parts[0], "()")
			mac := parts[3]
			iface := parts[5]
			entries = append(entries, anynetwork.ARPEntry{
				IPAddress:  ip,
				MACAddress: mac,
				Interface:  iface,
				Type:       "ethernet", // oder aus parts[len(parts)-1]
			})
		}
	}

	return entries, nil
}
