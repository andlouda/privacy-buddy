package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	anynetwork "privacy-buddy/backend/network"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	templatesFileName = "capture_templates.json"
)

// Singleton instance (thread-safe)
var (
	instance *AdvancedNetworkToolsService
	once     sync.Once
)

// GetAdvancedNetworkToolsService returns the singleton instance (instead von NewXYZ)
func GetAdvancedNetworkToolsService() *AdvancedNetworkToolsService {
	once.Do(func() {
		instance = &AdvancedNetworkToolsService{}
		log.Printf("✅ Singleton AdvancedNetworkToolsService created at %p", instance)
	})
	log.Printf("➡️ Returning Singleton instance: %p", instance)
	return instance
}

func (s *AdvancedNetworkToolsService) WailsInit(ctx context.Context) {
	log.Printf("✅ WailsInit called on instance %p with context %p", s, ctx)
	s.appCtx = ctx
}

// AdvancedNetworkToolsService provides advanced network diagnostic functionalities.
type AdvancedNetworkToolsService struct {
	appCtx context.Context

	// Packet capture state
	captureMutex sync.Mutex
	isCapturing  bool
	stopCapture  context.CancelFunc
}

// getTemplatesFilePath returns the full path to the templates file.
func (s *AdvancedNetworkToolsService) getTemplatesFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	appConfigDir := filepath.Join(configDir, "PrivacyBuddy")
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
			return []anynetwork.CaptureTemplate{}, nil
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
func (s *AdvancedNetworkToolsService) StartPacketCapture(iface string, bpfFilter string, durationSeconds int) error {
	log.Printf("DEBUG: StartPacketCapture called for instance %p. Current s.appCtx: %p", s, s.appCtx)

	if s.appCtx == nil {
		log.Println("CRITICAL ERROR: s.appCtx is nil. This indicates WailsInit was not called properly.")
		return fmt.Errorf("internal error: backend not initialized correctly (missing context)")
	}

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

	var captureCtx context.Context
	captureCtx, s.stopCapture = context.WithCancel(s.appCtx)
	s.isCapturing = true

	go func() {
		defer handle.Close()
		defer func() {
			s.captureMutex.Lock()
			s.isCapturing = false
			s.stopCapture = nil
			runtime.EventsEmit(s.appCtx, "packetCaptureStopped", "Capture finished or was stopped.")
			s.captureMutex.Unlock()
		}()

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		timer := time.NewTimer(time.Duration(durationSeconds) * time.Second)

		for {
			select {
			case packet := <-packetSource.Packets():
				cp := s.processPacket(packet)
				runtime.EventsEmit(s.appCtx, "packetCaptureEvent", cp)
			case <-timer.C:
				log.Println("Packet capture duration elapsed.")
				return
			case <-captureCtx.Done():
				log.Println("Packet capture cancelled.")
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
		s.stopCapture()
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
	} else if ip6Layer := packet.Layer(layers.LayerTypeIPv6); ip6Layer != nil {
		ipv6 := ip6Layer.(*layers.IPv6)
		cp.Source = ipv6.SrcIP.String()
		cp.Destination = ipv6.DstIP.String()
		summaryParts = append(summaryParts, fmt.Sprintf("IPv6 %s->%s Proto:%s", ipv6.SrcIP, ipv6.DstIP, ipv6.NextHeader))
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		var flags []string
		if tcp.SYN {
			flags = append(flags, "SYN")
		}
		if tcp.ACK {
			flags = append(flags, "ACK")
		}
		if tcp.FIN {
			flags = append(flags, "FIN")
		}
		if tcp.RST {
			flags = append(flags, "RST")
		}
		if tcp.PSH {
			flags = append(flags, "PSH")
		}
		if tcp.URG {
			flags = append(flags, "URG")
		}
		if tcp.ECE {
			flags = append(flags, "ECE")
		}
		if tcp.CWR {
			flags = append(flags, "CWR")
		}
		cp.Protocol = "TCP"
		summaryParts = append(summaryParts, fmt.Sprintf("TCP %d->%d Flags:[%s]", tcp.SrcPort, tcp.DstPort, strings.Join(flags, ",")))
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		cp.Protocol = "UDP"
		summaryParts = append(summaryParts, fmt.Sprintf("UDP %d->%d", udp.SrcPort, udp.DstPort))
	} else if icmp := packet.Layer(layers.LayerTypeICMPv4); icmp != nil {
		icmpv4 := icmp.(*layers.ICMPv4)
		cp.Protocol = "ICMPv4"
		summaryParts = append(summaryParts, fmt.Sprintf("ICMPv4 Type:%d Code:%d", icmpv4.TypeCode.Type(), icmpv4.TypeCode.Code()))
	} else if icmp6 := packet.Layer(layers.LayerTypeICMPv6); icmp6 != nil {
		icmpv6 := icmp6.(*layers.ICMPv6)
		cp.Protocol = "ICMPv6"
		summaryParts = append(summaryParts, fmt.Sprintf("ICMPv6 Type:%d Code:%d", icmpv6.TypeCode.Type(), icmpv6.TypeCode.Code()))
	}

	cp.Summary = strings.Join(summaryParts, " ")
	return cp
}

// GetCaptureTemplates returns all available templates.
func (s *AdvancedNetworkToolsService) GetCaptureTemplates() []anynetwork.CaptureTemplate {
	predefined := []anynetwork.CaptureTemplate{
		{Name: "HTTP/HTTPS", Description: "HTTP & HTTPS traffic", BPFFilter: "tcp port 80 or tcp port 443"},
		{Name: "DNS", Description: "DNS queries", BPFFilter: "udp port 53"},
		{Name: "ARP", Description: "Address resolution", BPFFilter: "arp"},
		{Name: "ICMP", Description: "Ping traffic", BPFFilter: "icmp"},
		{Name: "IPv4", Description: "All IPv4", BPFFilter: "ip"},
		{Name: "IPv6", Description: "All IPv6", BPFFilter: "ip6"},
		{Name: "SSH", Description: "SSH access", BPFFilter: "tcp port 22"},
		{Name: "RDP", Description: "Remote desktop", BPFFilter: "tcp port 3389"},
	}

	userTemplates, err := s.loadUserTemplates()
	if err != nil {
		log.Printf("WARN: Could not load user templates: %v", err)
		return predefined
	}
	return append(predefined, userTemplates...)
}

// SaveCaptureTemplate adds a new user-defined capture template.
func (s *AdvancedNetworkToolsService) SaveCaptureTemplate(tpl anynetwork.CaptureTemplate) error {
	userTemplates, err := s.loadUserTemplates()
	if err != nil {
		return fmt.Errorf("could not load user templates: %w", err)
	}

	for _, t := range userTemplates {
		if strings.EqualFold(t.Name, tpl.Name) {
			return fmt.Errorf("template with name '%s' already exists", tpl.Name)
		}
	}

	userTemplates = append(userTemplates, tpl)
	return s.saveUserTemplates(userTemplates)
}
