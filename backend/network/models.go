package network

import (
	"net"
)

// NetworkInterface represents a network interface with detailed information.
type NetworkInterface struct {
	Name        string   `json:"Name"`
	DisplayName string   `json:"DisplayName"` // From net.Interface.Name or pcap.Interface.Description
	Description string   `json:"Description"` // From pcap.Interface.Description
		HardwareAddr net.HardwareAddr `json:"HardwareAddr"`
	MTU         int      `json:"MTU"`
	Flags       []string `json:"Flags"`
	Addrs       []string `json:"Addrs"`
	IsUp        bool     `json:"IsUp"`
	IsLoopback  bool     `json:"IsLoopback"`
	IsBroadcast bool     `json:"IsBroadcast"`
	IsPointToPoint bool  `json:"IsPointToPoint"`
	IsMulticast bool     `json:"IsMulticast"`
}

// CapturedPacket represents a captured network packet.
type CapturedPacket struct {
	Timestamp   string `json:"Timestamp"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Protocol    string `json:"Protocol"`
	Length      int    `json:"Length"`
	Summary     string `json:"Summary"` // Human-readable summary
}

// ARPEntry represents a single entry in the ARP cache.
type ARPEntry struct {
	IPAddress   string `json:"IPAddress"`
	MACAddress  string `json:"MACAddress"`
	Interface   string `json:"Interface"`
	Type        string `json:"Type"` // e.g., "dynamic", "static"
}

// NetworkConnection represents an active network connection.
type NetworkConnection struct {
	FD          uint64 `json:"FD"`
	Family      uint32 `json:"Family"` // e.g., AF_INET, AF_INET6
	Type        uint32 `json:"Type"`   // e.g., SOCK_STREAM, SOCK_DGRAM
	LocalIP     string `json:"LocalIP"`
	LocalPort   uint32 `json:"LocalPort"`
	RemoteIP    string `json:"RemoteIP"`
	RemotePort  uint32 `json:"RemotePort"`
	Status      string `json:"Status"` // e.g., ESTABLISHED, LISTEN, CLOSE_WAIT
	PID         int32  `json:"PID"`
	ProcessName string `json:"ProcessName"`
	Protocol    string `json:"Protocol"` // e.g., "tcp", "udp"
}

// CaptureTemplate defines a pre-configured BPF filter.
type CaptureTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BPFFilter   string `json:"bpfFilter"`
	Duration    int    `json:"duration"`
}
