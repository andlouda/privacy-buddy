package network

// NetworkInterface represents a network interface with detailed information.
type NetworkInterface struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"` // From net.Interface.Name or pcap.Interface.Description
	Description string   `json:"description"` // From pcap.Interface.Description
	HardwareAddr string   `json:"hardwareAddr"`
	MTU         int      `json:"mtu"`
	Flags       []string `json:"flags"`
	Addrs       []string `json:"addrs"`
	IsUp        bool     `json:"isUp"`
	IsLoopback  bool     `json:"isLoopback"`
	IsBroadcast bool     `json:"isBroadcast"`
	IsPointToPoint bool  `json:"isPointToPoint"`
	IsMulticast bool     `json:"isMulticast"`
}

// CapturedPacket represents a captured network packet.
type CapturedPacket struct {
	Timestamp   string `json:"timestamp"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Protocol    string    `json:"protocol"`
	Length      int       `json:"length"`
	Payload     []byte    `json:"payload"` // Raw packet data
	Summary     string    `json:"summary"` // Human-readable summary
}

// ARPEntry represents a single entry in the ARP cache.
type ARPEntry struct {
	IPAddress   string `json:"ipAddress"`
	MACAddress  string `json:"macAddress"`
	Interface   string `json:"interface"`
	Type        string `json:"type"` // e.g., "dynamic", "static"
}

// NetworkConnection represents an active network connection.
type NetworkConnection struct {
	FD          uint64 `json:"fd"`
	Family      uint32 `json:"family"` // e.g., AF_INET, AF_INET6
	Type        uint32 `json:"type"`   // e.g., SOCK_STREAM, SOCK_DGRAM
	LocalIP     string `json:"localIP"`
	LocalPort   uint32 `json:"localPort"`
	RemoteIP    string `json:"remoteIP"`
	RemotePort  uint32 `json:"remotePort"`
	Status      string `json:"status"` // e.g., ESTABLISHED, LISTEN, CLOSE_WAIT
	PID         int32  `json:"pid"`
	ProcessName string `json:"processName"`
	Protocol    string `json:"protocol"` // e.g., "tcp", "udp"
}

// CaptureTemplate defines a pre-configured BPF filter.
type CaptureTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BPFFilter   string `json:"bpfFilter"`
	Duration    int    `json:"duration"`
}
