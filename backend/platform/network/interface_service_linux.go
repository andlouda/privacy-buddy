//go:build linux

package network

import (
	"fmt"
	"net"

	"github.com/google/gopacket/pcap"

	anynetwork "privacy-buddy/backend/network"
)

// LinuxNetworkInterfaceService provides Linux-specific implementation for network interface listing.
type LinuxNetworkInterfaceService struct{}

// ListInterfaces lists all network interfaces on Linux.
func (s *LinuxNetworkInterfaceService) ListInterfaces() ([]anynetwork.NetworkInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	pcapDevices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("failed to find pcap devices: %w", err)
	}

	var netInterfaces []anynetwork.NetworkInterface
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			// Log the error but continue with other interfaces
			fmt.Printf("warning: failed to get addresses for interface %s: %v\n", iface.Name, err)
			continue
		}

		var ipAddrs []string
		for _, addr := range addrs {
			ipAddrs = append(ipAddrs, addr.String())
		}

		flags := []string{}
		if iface.Flags&net.FlagUp != 0 {
			flags = append(flags, "up")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "broadcast")
		}
		if iface.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "loopback")
		}
		if iface.Flags&net.FlagPointToPoint != 0 {
			flags = append(flags, "pointtopoint")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "multicast")
		}

		description := ""
		for _, dev := range pcapDevices {
			if dev.Name == iface.Name {
				description = dev.Description
				break
			}
		}

		netInterfaces = append(netInterfaces, anynetwork.NetworkInterface{
			Name:           iface.Name,
			DisplayName:    description, // Use pcap description as DisplayName
			Description:    description,
			HardwareAddr:   iface.HardwareAddr,
			MTU:            iface.MTU,
			Flags:          flags,
			Addrs:          ipAddrs,
			IsUp:           iface.Flags&net.FlagUp != 0,
			IsLoopback:     iface.Flags&net.FlagLoopback != 0,
			IsBroadcast:    iface.Flags&net.FlagBroadcast != 0,
			IsPointToPoint: iface.Flags&net.FlagPointToPoint != 0,
			IsMulticast:    iface.Flags&net.FlagMulticast != 0,
		})
	}
	return netInterfaces, nil
}
