//go:build windows

package network

import (
	"fmt"

	"github.com/google/gopacket/pcap"
	gopsnet "github.com/shirou/gopsutil/v3/net"

	anynetwork "privacy-buddy/backend/network"
)

// WindowsNetworkInterfaceService provides Windows-specific implementation for network interface listing.
type WindowsNetworkInterfaceService struct{}

// ListInterfaces lists all network interfaces on Windows.
func (s *WindowsNetworkInterfaceService) ListInterfaces() ([]anynetwork.NetworkInterface, error) {
	interfaces, err := gopsnet.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	pcapDevices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("failed to find pcap devices: %w", err)
	}

	var netInterfaces []anynetwork.NetworkInterface
	for _, iface := range interfaces {
		// Debug output for MAC address
		fmt.Printf("Interface: %s, Flags: %v, MAC Address: %s\n", iface.Name, iface.Flags, iface.HardwareAddr)

		addrs := iface.Addrs

		var ipAddrs []string
		for _, addr := range addrs {
			ipAddrs = append(ipAddrs, addr.Addr)
		}

		flags := []string{}
		if containsFlag(iface.Flags, "up") {
			flags = append(flags, "up")
		}
		if containsFlag(iface.Flags, "broadcast") {
			flags = append(flags, "broadcast")
		}
		if containsFlag(iface.Flags, "loopback") {
			flags = append(flags, "loopback")
		}
		if containsFlag(iface.Flags, "pointtopoint") {
			flags = append(flags, "pointtopoint")
		}
		if containsFlag(iface.Flags, "multicast") {
			flags = append(flags, "multicast")
		}

		description := ""
		for _, dev := range pcapDevices {
			if dev.Name == iface.Name {
				description = dev.Description
				break
			}
		}

		if iface.HardwareAddr == "" {
			fmt.Printf("Warning: MAC address is empty for interface %s\n", iface.Name)
		}
		netInterfaces = append(netInterfaces, anynetwork.NetworkInterface{
			Name:           iface.Name,
			DisplayName:    iface.Name, // gopsutil doesn't provide DisplayName directly
			Description:    description,
			HardwareAddr:   iface.HardwareAddr, // Default from gopsutil
			// Attempt to get HardwareAddr from pcap device if available and gopsutil's is empty
			// Note: pcap.Interface.HardwareAddr is a []byte, convert to string
			// This part needs to be carefully integrated into the loop where `dev` is matched.
			// For now, let's assume we'll get it from the matched `dev`.
			MTU:            iface.MTU,
			Flags:          flags,
			Addrs:          ipAddrs,
			IsUp:           containsFlag(iface.Flags, "up"),
			IsLoopback:     containsFlag(iface.Flags, "loopback"),
			IsBroadcast:    containsFlag(iface.Flags, "broadcast"),
			IsPointToPoint: containsFlag(iface.Flags, "pointtopoint"),
			IsMulticast:    containsFlag(iface.Flags, "multicast"),
		})
	}
	return netInterfaces, nil
}

func containsFlag(flags []string, flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}