//go:build windows

package network

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	anynetwork "privacy-buddy/backend/network"
)

// WindowsNetworkConnectionService provides Windows-specific implementation for listing network connections.
type WindowsNetworkConnectionService struct{}

// GetConnections lists all active network connections on Windows.
func (s *WindowsNetworkConnectionService) GetConnections() ([]anynetwork.NetworkConnection, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	var netConnections []anynetwork.NetworkConnection
	for _, conn := range connections {
		processName := "N/A"
		if conn.Pid != 0 {
			proc, err := process.NewProcess(conn.Pid)
			if err == nil {
				name, err := proc.Name()
				if err == nil {
					processName = name
				}
			}
		}

		netConnections = append(netConnections, anynetwork.NetworkConnection{
			FD:          uint64(conn.Fd),
			Family:      conn.Family,
			Type:        conn.Type,
			LocalIP:     conn.Laddr.IP,
			LocalPort:   conn.Laddr.Port,
			RemoteIP:    conn.Raddr.IP,
			RemotePort:  conn.Raddr.Port,
			Status:      conn.Status,
			PID:         conn.Pid,
			ProcessName: processName,
		})
	}

	return netConnections, nil
}
