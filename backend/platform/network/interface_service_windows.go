//go:build windows

package network

import (
	"fmt"
	"net"

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
		fmt.Printf("Interface: %s, Flags: %v, MAC Address (gopsutil): %s\n", iface.Name, iface.Flags, iface.HardwareAddr)

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

				displayName := iface.Name
		description := ""

		for _, dev := range pcapDevices {
			if dev.Name == iface.Name {
				description = dev.Description
				displayName = dev.Description // Use pcap description as DisplayName
				break
			}
		}

		netInterfaces = append(netInterfaces, anynetwork.NetworkInterface{
			Name:           iface.Name,
			DisplayName:    displayName,
			Description:    description,
						HardwareAddr:   func() net.HardwareAddr{
				mac, err := net.ParseMAC(iface.HardwareAddr)
				if err != nil {
					fmt.Printf("Warning: Failed to parse MAC address %s for interface %s: %v\n", iface.HardwareAddr, iface.Name, err)
					return nil
				}
				return mac
			}(),
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