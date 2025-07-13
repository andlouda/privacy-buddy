package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"strings"

	anynetwork "privacy-buddy/backend/network"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	templatesFileName = "capture_templates.json"
)

// AdvancedNetworkToolsService provides advanced network diagnostic functionalities.
type AdvancedNetworkToolsService struct {
	appCtx context.Context

	// Packet capture state
	captureMutex   sync.Mutex
	isCapturing    bool
	stopCapture    context.CancelFunc
}

// NewAdvancedNetworkToolsService creates a new instance of AdvancedNetworkToolsService.
func NewAdvancedNetworkToolsService() *AdvancedNetworkToolsService {
	return &AdvancedNetworkToolsService{}
}

// WailsInit is called by Wails when the application is starting.
func (s *AdvancedNetworkToolsService) WailsInit(ctx context.Context) {
	s.appCtx = ctx
}

// getTemplatesFilePath returns the full path to the templates file.
func (s *AdvancedNetworkToolsService) getTemplatesFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	appConfigDir := filepath.Join(configDir, "PrivacyBuddy") // Use a specific directory for your app
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app config directory: %w", err)
	}
	return filepath.Join(appConfigDir, templatesFileName), nil
}

// loadUserTemplates loads user-defined templates from a file.
func (s *AdvancedNetworkToolsService) loadUserTemplates() ([]anynetwork.CaptureTemplate, error) {
	filePath, err := s.getTemplatesFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []anynetwork.CaptureTemplate{}, nil // No file, return empty slice
		}
		return nil, fmt.Errorf("failed to read templates file: %w", err)
	}

	var templates []anynetwork.CaptureTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templates: %w", err)
	}
	return templates, nil
}

// saveUserTemplates saves user-defined templates to a file.
func (s *AdvancedNetworkToolsService) saveUserTemplates(templates []anynetwork.CaptureTemplate) error {
	filePath, err := s.getTemplatesFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal templates: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write templates file: %w", err)
	}
	return nil
}

// StartPacketCapture starts capturing packets on a given interface.
// It runs asynchronously and sends packets back to the frontend via events.
func (s *AdvancedNetworkToolsService) StartPacketCapture(iface string, bpfFilter string, durationSeconds int) error {
	s.captureMutex.Lock()
	defer s.captureMutex.Unlock()

	if s.isCapturing {
		return fmt.Errorf("a packet capture is already in progress")
	}

	log.Printf("Attempting to open device: %s with BPF filter: %s for %d seconds", iface, bpfFilter, durationSeconds)
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Printf("ERROR: Failed to open device %s: %v", iface, err)
		return fmt.Errorf("error opening device %s: %w", iface, err)
	}

	if err := handle.SetBPFFilter(bpfFilter); err != nil {
		return fmt.Errorf("error setting BPF filter: %w", err)
	}

	// Create a new context for this capture session
	var captureCtx context.Context
	captureCtx, s.stopCapture = context.WithCancel(s.appCtx)

	s.isCapturing = true

	// Run the capture in a goroutine so it doesn't block the UI
	go func() {
		defer handle.Close()
		defer func() {
			s.captureMutex.Lock()
			s.isCapturing = false
			s.stopCapture = nil
			
			// Inform the frontend that capture has stopped
			runtime.EventsEmit(s.appCtx, "packetCaptureStopped", "Capture finished or was stopped.")
			s.captureMutex.Unlock()
		}()

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		timer := time.NewTimer(time.Duration(durationSeconds) * time.Second)

		for {
			select {
			case packet := <-packetSource.Packets():
				// Process and emit the packet to the frontend
				cp := s.processPacket(packet)
				runtime.EventsEmit(s.appCtx, "packetCaptureEvent", cp)
			case <-timer.C:
				log.Println("Packet capture duration elapsed.")
				return
			case <-captureCtx.Done():
				log.Println("Packet capture cancelled by user or app shutdown.")
				return
			}
		}
	}()

	return nil
}

// StopPacketCapture manually stops an ongoing packet capture.
func (s *AdvancedNetworkToolsService) StopPacketCapture() {
	s.captureMutex.Lock()
	defer s.captureMutex.Unlock()

	if s.stopCapture != nil {
		s.stopCapture() // Signal the goroutine to stop
	}
}

// processPacket converts a gopacket.Packet into our custom struct.
func (s *AdvancedNetworkToolsService) processPacket(packet gopacket.Packet) anynetwork.CapturedPacket {
	cp := anynetwork.CapturedPacket{
		Timestamp: packet.Metadata().Timestamp.Format(time.RFC3339Nano),
		Length:    packet.Metadata().Length,
	}

	summaryParts := []string{}

	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		summaryParts = append(summaryParts, fmt.Sprintf("Eth %s->%s", ethLayer.(*layers.Ethernet).SrcMAC, ethLayer.(*layers.Ethernet).DstMAC))
	}

	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip4 := ipLayer.(*layers.IPv4)
		cp.Source = ip4.SrcIP.String()
		cp.Destination = ip4.DstIP.String()
		summaryParts = append(summaryParts, fmt.Sprintf("IPv4 %s->%s Proto:%s", ip4.SrcIP, ip4.DstIP, ip4.Protocol))
	} else if ipLayer := packet.Layer(layers.LayerTypeIPv6); ipLayer != nil {
		ip6 := ipLayer.(*layers.IPv6)
		cp.Source = ip6.SrcIP.String()
		cp.Destination = ip6.DstIP.String()
		summaryParts = append(summaryParts, fmt.Sprintf("IPv6 %s->%s Proto:%s", ip6.SrcIP, ip6.DstIP, ip6.NextHeader))
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		var tcpFlags []string
		if tcp.SYN { tcpFlags = append(tcpFlags, "SYN") }
		if tcp.ACK { tcpFlags = append(tcpFlags, "ACK") }
		if tcp.FIN { tcpFlags = append(tcpFlags, "FIN") }
		if tcp.RST { tcpFlags = append(tcpFlags, "RST") }
		if tcp.PSH { tcpFlags = append(tcpFlags, "PSH") }
		if tcp.URG { tcpFlags = append(tcpFlags, "URG") }
		if tcp.ECE { tcpFlags = append(tcpFlags, "ECE") }
		if tcp.CWR { tcpFlags = append(tcpFlags, "CWR") }
		cp.Protocol = "TCP"
		summaryParts = append(summaryParts, fmt.Sprintf("TCP %d->%d Flags:[%s]", tcp.SrcPort, tcp.DstPort, strings.Join(tcpFlags, ",")))
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		cp.Protocol = "UDP"
		summaryParts = append(summaryParts, fmt.Sprintf("UDP %d->%d", udp.SrcPort, udp.DstPort))
	} else if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
		icmp := icmpLayer.(*layers.ICMPv4)
		cp.Protocol = "ICMPv4"
		summaryParts = append(summaryParts, fmt.Sprintf("ICMPv4 Type:%d Code:%d", icmp.TypeCode.Type(), icmp.TypeCode.Code()))
	} else if icmp6Layer := packet.Layer(layers.LayerTypeICMPv6); icmp6Layer != nil {
		icmp6 := icmp6Layer.(*layers.ICMPv6)
		cp.Protocol = "ICMPv6"
		summaryParts = append(summaryParts, fmt.Sprintf("ICMPv6 Type:%d Code:%d", icmp6.TypeCode.Type(), icmp6.TypeCode.Code()))
	}

	cp.Summary = strings.Join(summaryParts, " ")

	return cp
}

// GetCaptureTemplates returns a list of pre-configured and user-defined packet capture templates.
func (s *AdvancedNetworkToolsService) GetCaptureTemplates() []anynetwork.CaptureTemplate {
	// Pre-configured templates
	predefinedTemplates := []anynetwork.CaptureTemplate{
		{
			Name:        "HTTP/HTTPS Traffic",
			Description: "Captures HTTP (port 80) and HTTPS (port 443) traffic.",
			BPFFilter:   "tcp port 80 or tcp port 443",
		},
		{
			Name:        "DNS Queries",
			Description: "Captures DNS (port 53) queries and responses.",
			BPFFilter:   "udp port 53",
		},
		{
			Name:        "ARP Traffic",
			Description: "Captures Address Resolution Protocol (ARP) requests and replies.",
			BPFFilter:   "arp",
		},
		{
			Name:        "ICMP (Ping) Traffic",
			Description: "Captures Internet Control Message Protocol (ICMP) traffic, commonly used by ping.",
			BPFFilter:   "icmp",
		},
		{
			Name:        "SSH Traffic",
			Description: "Captures Secure Shell (SSH) traffic (port 22).",
			BPFFilter:   "tcp port 22",
		},
		{
			Name:        "RDP Traffic",
			Description: "Captures Remote Desktop Protocol (RDP) traffic (port 3389).",
			BPFFilter:   "tcp port 3389",
		},
		{
			Name:        "All IPv4 Traffic",
			Description: "Captures all IPv4 traffic.",
			BPFFilter:   "ip",
		},
		{
			Name:        "All IPv6 Traffic",
			Description: "Captures all IPv6 traffic.",
			BPFFilter:   "ip6",
		},
		{
			Name:        "Broadcast Traffic",
			Description: "Captures broadcast packets.",
			BPFFilter:   "broadcast",
		},
		{
			Name:        "Multicast Traffic",
			Description: "Captures multicast packets.",
			BPFFilter:   "multicast",
		},
	}

	userTemplates, err := s.loadUserTemplates()
	if err != nil {
		log.Printf("ERROR: Failed to load user templates: %v", err)
		// Return only predefined templates if user templates can't be loaded
		return predefinedTemplates
	}

	// Combine predefined and user-defined templates
	return append(predefinedTemplates, userTemplates...)
}

// SaveCaptureTemplate saves a new user-defined packet capture template.
func (s *AdvancedNetworkToolsService) SaveCaptureTemplate(template anynetwork.CaptureTemplate) error {
	// Load existing user templates
	userTemplates, err := s.loadUserTemplates()
	if err != nil {
		return fmt.Errorf("failed to load existing user templates: %w", err)
	}

	// Check for duplicate name (case-insensitive)
	for _, t := range userTemplates {
		if strings.EqualFold(t.Name, template.Name) {
			return fmt.Errorf("template with name '%s' already exists", template.Name)
		}
	}

	// Add the new template
	userTemplates = append(userTemplates, template)

	// Save all templates back to file
	return s.saveUserTemplates(userTemplates)
}

// ... other methods like GetARPCache, GetActiveConnections would be here ...