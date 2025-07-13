package tools

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	anynetwork "privacy-buddy/backend/network"

	"github.com/go-ping/ping"
	"github.com/google/gopacket/pcap" // Hinzugefügt für Netzwerkschnittstellen-Erkennung
)

// NetworkToolsService bietet Funktionen für Netzwerkdiagnose-Tools.
type NetworkToolsService struct {
	tracerouteSvc TracerouteService
}

// TracerouteService definiert die Schnittstelle für plattformspezifische Traceroute-Implementierungen.
type TracerouteService interface {
	Traceroute(host string) ([]TracerouteHop, error)
}

// TracerouteHop strukturiert einen Hop in einem Traceroute-Ergebnis.
type TracerouteHop struct {
	N       int    `json:"n"`
	Host    string `json:"host"`
	Address string `json:"address"`
	RTT     string `json:"rtt"`
}

// NewNetworkToolsService erstellt eine neue Instanz des NetworkToolsService.
func NewNetworkToolsService(tracerouteSvc TracerouteService) *NetworkToolsService {
	return &NetworkToolsService{
		tracerouteSvc: tracerouteSvc,
	}
}

// PingResult strukturiert die Ergebnisse eines Ping-Befehls.
type PingResult struct {
	Host    string `json:"host"`
	IP      string `json:"ip"`
	Packets int    `json:"packets"`
	Loss    string `json:"loss"`
	MinRtt  string `json:"minRtt"`
				AvgRtt  string `json:"avgRtt"`
	MaxRtt  string `json:"maxRtt"`
	Error   string `json:"error,omitempty"`
}

// runNativePing führt den nativen Ping-Befehl des Betriebssystems aus und parst dessen Ausgabe.
// Dies dient als Fallback, wenn die go-ping Bibliothek aufgrund von Berechtigungsproblemen fehlschlägt.
func runNativePing(host string) (*PingResult, error) {
	var cmd *exec.Cmd
	var args []string

	switch runtime.GOOS {
	case "windows":
		args = []string{"-n", "4", host} // -n for count on Windows
		cmd = exec.Command("ping", args...)
	case "linux", "darwin":
		args = []string{"-c", "4", host} // -c for count on Linux/macOS
		cmd = exec.Command("ping", args...)
	default:
		return nil, fmt.Errorf("natives Ping wird auf dieser Plattform (%s) nicht unterstützt", runtime.GOOS)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("nativer Ping-Befehl fehlgeschlagen: %w, %s", err, stderr.String())
	}

	output := out.String()
	result := &PingResult{
		Host: host,
		// IP, RTTs, Loss will be parsed from output
	}

	// Simplified parsing for native ping output
	// This parsing is highly dependent on OS and OS, and might need refinement.
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "bytes from") || strings.Contains(line, "Antwort von") { // Linux/macOS and Windows success
			// Example: 64 bytes from 8.8.8.8: icmp_seq=1 ttl=117 time=12.3 ms
			// Example: Antwort von 8.8.8.8: Bytes=32 Zeit=12ms TTL=117
			if strings.Contains(line, "time=") { // Linux/macOS
				if match := regexp.MustCompile(`time=(\d+\.?\d*) ms`).FindStringSubmatch(line); len(match) > 1 {
					if rtt, _ := strconv.ParseFloat(match[1], 64); rtt > 0 {
						// For simplicity, we'll just take the last RTT as AvgRtt for now.
						// A more robust parser would collect all RTTs and calculate min/avg/max.
						result.AvgRtt = fmt.Sprintf("%.1fms", rtt)
						result.MinRtt = result.AvgRtt // Placeholder
						result.MaxRtt = result.AvgRtt // Placeholder
					}
				}
			} else if strings.Contains(line, "Zeit=") { // Windows
				if match := regexp.MustCompile(`Zeit=(\d+)ms`).FindStringSubmatch(line); len(match) > 1 {
					if rtt, _ := strconv.Atoi(match[1]); rtt > 0 {
						result.AvgRtt = fmt.Sprintf("%dms", rtt)
						result.MinRtt = result.AvgRtt // Placeholder
						result.MaxRtt = result.AvgRtt // Placeholder
					}
				}
			}
			result.Packets++ // Count successful packets
		} else if strings.Contains(line, "Packets: Sent =") || strings.Contains(line, "Pakete: Gesendet =") {
			// Linux/macOS: 4 packets transmitted, 4 received, 0% packet loss, time 3005ms
			// Windows: Pakete: Gesendet = 4, Empfangen = 4, Verloren = 0 (0% Verlust)
			if strings.Contains(line, "transmitted") { // Linux/macOS summary
				if match := regexp.MustCompile(`(\d+) packets transmitted, (\d+) received, (\d+)% packet loss`).FindStringSubmatch(line); len(match) > 3 {
					result.Packets, _ = strconv.Atoi(match[1])
					// result.PacketsRecv, _ = strconv.Atoi(match[2]) // Not directly in PingResult
					result.Loss = fmt.Sprintf("%s%%", match[3])
				}
			} else if strings.Contains(line, "Gesendet =") { // Windows summary
				if match := regexp.MustCompile(`Gesendet = (\d+), Empfangen = (\d+), Verloren = (\d+) \((\d+)% Verlust\)`).FindStringSubmatch(line); len(match) > 4 {
					result.Packets, _ = strconv.Atoi(match[1])
					// result.PacketsRecv, _ = strconv.Atoi(match[2]) // Not directly in PingResult
					result.Loss = fmt.Sprintf("%s%%", match[4])
				}
			}
		}
	}

	// Fallback for IP if not parsed from output (e.g., if only host was provided)
	if result.IP == "" {
		if addrs, err := net.LookupIP(host); err == nil && len(addrs) > 0 {
			result.IP = addrs[0].String()
		} else {
			result.IP = "unbekannt"
		}
	}

	// If no packets were received, assume 100% loss
	if result.Packets == 0 || (result.Packets > 0 && result.Loss == "100.00%") {
		result.Error = "Ziel nicht erreichbar oder Timeout (nativer Ping)."
	}

	return result, nil
}

// Ping führt einen Ping-Befehl aus und gibt die Ergebnisse zurück.
func (s *NetworkToolsService) Ping(host string) (*PingResult, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		// Fehler beim Erstellen des Pingers, z.B. ungültiger Hostname
		return nil, fmt.Errorf("fehler beim Erstellen des Pingers für %s: %w", host, err)
	}

	pinger.Count = 4                 // Sende 4 Pakete
	pinger.Timeout = time.Second * 5 // 5 Sekunden Timeout
	pinger.SetPrivileged(true)       // ⚠️ Wichtig für Windows + Linux ohne sysctl/cap_net_raw

	// Führe den Ping aus. Die Run()-Methode blockiert und gibt Fehler zurück,
	// die während des Ping-Vorgangs auftreten (z.B. Berechtigungsprobleme).
	err = pinger.Run() // Assign error to 'err'

	stats := pinger.Statistics()

	result := &PingResult{
		Host:    stats.Addr,
		IP:      stats.IPAddr.String(),
		Packets: stats.PacketsSent,
		MinRtt:  stats.MinRtt.String(),
		AvgRtt:  stats.AvgRtt.String(),
		MaxRtt:  stats.MaxRtt.String(),
	}

	// Behandle Paketverlust, insbesondere den NaN-Fall
	if stats.PacketsSent > 0 {
		result.Loss = fmt.Sprintf("%.2f%%", stats.PacketLoss)
	} else {
		// Wenn keine Pakete gesendet wurden, gib N/A für den Verlust an
		result.Loss = "N/A"
	}

	// Priorisiere Fehler, die direkt von pinger.Run() zurückgegeben wurden
	if err != nil { // Use 'err' here
		// Überprüfe auf spezifische häufige Fehler von go-ping im Zusammenhang mit Berechtigungen/Sockets
		errMsg := err.Error() // Use 'err' here
		if strings.Contains(errMsg, "socket: The requested protocol has not been configured into the system") ||
			strings.Contains(errMsg, "permission denied") ||
			strings.Contains(errMsg, "operation not permitted") {
			// Fallback to native ping if go-ping fails due to privilege/socket issues
			nativeResult, nativeErr := runNativePing(host)
				if nativeErr != nil {
				result.Error = fmt.Sprintf("Ping-Fehler (go-ping): %v. Fallback (native) fehlgeschlagen: %v", errMsg, nativeErr)
			} else {
				// If native ping succeeded, use its results
				result = nativeResult
			}
		} else {
			// Generische Fehlermeldung für andere Fehler von go-ping
			result.Error = fmt.Sprintf("Ping-Fehler (go-ping): %v", errMsg)
		}
	} else if stats.PacketsSent == 0 {
		// Keine Pakete gesendet, aber kein direkter Run-Fehler (z.B. Host nicht auflösbar)
		result.Error = "Ping konnte keine Pakete senden. Überprüfen Sie Berechtigungen oder Netzwerkverbindung."
	} else if stats.PacketsRecv == 0 && stats.PacketLoss == 100 {
		// 100% Paketverlust, Ziel nicht erreichbar oder Timeout
		result.Error = "Ziel nicht erreichbar oder Timeout."
	}

	return result, nil
}

// Traceroute führt einen Traceroute-Befehl aus und gibt die Ergebnisse zurück.
func (s *NetworkToolsService) Traceroute(host string) ([]TracerouteHop, error) {
	return s.tracerouteSvc.Traceroute(host)
}

// GetNetworkInterfaces listet alle verfügbaren Netzwerkschnittstellen auf.
func (s *NetworkToolsService) GetNetworkInterfaces() ([]anynetwork.NetworkInterface, error) {
	var interfaces []anynetwork.NetworkInterface

	// Get interfaces from pcap for the actual device name needed by pcap.OpenLive
	pcapDevs, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("error finding pcap devices: %w", err)
	}

	// Get interfaces from net package for IP addresses, MAC, and status
	netDevs, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error finding net interfaces: %w", err)
	}

	// Create a map for net.Interfaces for quick lookup by name
	netDevMap := make(map[string]net.Interface)
	for _, dev := range netDevs {
		netDevMap[dev.Name] = dev
	}

	for _, pcapDev := range pcapDevs {
		log.Printf("PCAP Device - Name: %s, Description: %s", pcapDev.Name, pcapDev.Description)
		var addresses []string
		var hardwareAddr string
		var mtu int
		var flagsList []string
		var isUp, isLoopback, isBroadcast, isPointToPoint, isMulticast bool

		// Try to find the corresponding net.Interface
		if netDev, found := netDevMap[pcapDev.Name]; found {
			addrs, _ := netDev.Addrs()
			for _, addr := range addrs {
				addresses = append(addresses, addr.String())
			}
			hardwareAddr = netDev.HardwareAddr.String()
			mtu = netDev.MTU
			flagsList = strings.Split(netDev.Flags.String(), ",")
			isUp = (netDev.Flags & net.FlagUp) != 0
			isLoopback = (netDev.Flags & net.FlagLoopback) != 0
			isBroadcast = (netDev.Flags & net.FlagBroadcast) != 0
			isPointToPoint = (netDev.Flags & net.FlagPointToPoint) != 0
			isMulticast = (netDev.Flags & net.FlagMulticast) != 0
		} else {
			// Fallback for addresses if net.Interface not found (less common for pcap devices)
			for _, addr := range pcapDev.Addresses {
				addresses = append(addresses, addr.IP.String())
			}
		}

		interfaces = append(interfaces, anynetwork.NetworkInterface{
			Name:           pcapDev.Name,        // This is the name pcap.OpenLive expects
			DisplayName:    pcapDev.Description, // User-friendly name
			Description:    pcapDev.Description,
			HardwareAddr:   hardwareAddr,
			MTU:            mtu,
			Flags:          flagsList,
			Addrs:          addresses,
			IsUp:           isUp,
			IsLoopback:     isLoopback,
			IsBroadcast:    isBroadcast,
			IsPointToPoint: isPointToPoint,
			IsMulticast:    isMulticast,
		})
	}

	return interfaces, nil
}