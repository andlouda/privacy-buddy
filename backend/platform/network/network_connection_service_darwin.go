//go:build darwin

package network

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	anynetwork "privacy-buddy/backend/network"
)

// DarwinNetworkConnectionService provides macOS-specific implementation for listing network connections.
type DarwinNetworkConnectionService struct{}

// GetConnections lists all active network connections on macOS.
func (s *DarwinNetworkConnectionService) GetConnections() ([]anynetwork.NetworkConnection, error) {
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
			FD:          conn.Fd,
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


