//go:build linux

package network

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	anynetwork "privacy-buddy/backend/network"
)

// LinuxARPCacheService provides Linux-specific implementation for reading the ARP cache.
type LinuxARPCacheService struct{}

// GetARPEntries reads the ARP cache from /proc/net/arp on Linux.
func (s *LinuxARPCacheService) GetARPEntries() ([]anynetwork.ARPEntry, error) {
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/net/arp: %w", err)
	}
	defer file.Close()

	var entries []anynetwork.ARPEntry
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Skip header line

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue // Malformed line
		}

		ipAddress := fields[0]
		macAddress := fields[3]
		flags := fields[2] // Flags field contains type information
		device := fields[5]

		entryType := "dynamic"
		if strings.Contains(flags, "0x2") { // AT_PERM flag for static entries
			entryType = "static"
		}

		entries = append(entries, anynetwork.ARPEntry{
			IPAddress:  ipAddress,
			MACAddress: macAddress,
			Interface:  device,
			Type:       entryType,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /proc/net/arp: %w", err)
	}

	return entries, nil
}


