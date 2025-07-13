//go:build windows

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/wailsapp/wails/v2/pkg/runtime" // Import runtime

	anynetwork "privacy-buddy/backend/network"
)

// WindowsPacketCaptureService provides Windows-specific implementation for packet capturing.
type WindowsPacketCaptureService struct {
	handle *pcap.Handle
	ctx    context.Context // Add context to the struct
}

// WailsInit is called at application startup
func (s *WindowsPacketCaptureService) WailsInit(ctx context.Context) {
	s.ctx = ctx
}

// StartCapture starts capturing packets on the specified interface with an optional BPF filter and duration.
func (s *WindowsPacketCaptureService) StartCapture(ctx context.Context, interfaceName string, bpfFilter string, duration time.Duration) (<-chan anynetwork.CapturedPacket, error) {

	if s.handle != nil {
		return nil, fmt.Errorf("capture already in progress")
	}

	// Open the device for capturing
	handle, err := pcap.OpenLive(interfaceName, 1600, true, time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to open device %s: %w", interfaceName, err)
	}
	s.handle = handle

	// Set BPF filter if provided
	if bpfFilter != "" {
		if err := s.handle.SetBPFFilter(bpfFilter); err != nil {
			s.handle.Close()
			s.handle = nil
			return nil, fmt.Errorf("failed to set BPF filter: %w", err)
		}
	}

	packetSource := gopacket.NewPacketSource(s.handle, s.handle.LinkType())
	packetChannel := make(chan anynetwork.CapturedPacket)

	go func() {
		defer close(packetChannel)
		defer s.StopCapture()

		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Packet capture stopped due to context cancellation.")
				return
			case <-ticker.C:
				log.Printf("Packet capture stopped due to duration expiry.")
				return
			case packet := <-packetSource.Packets():
				if packet == nil {
					log.Printf("Packet source closed.")
					return
				}
				// Extract relevant information and send to channel
				srcIP := "Unknown"
				dstIP := "Unknown"
				protocol := "Unknown"

				networkLayer := packet.NetworkLayer()
				if networkLayer != nil {
					if ipv4, ok := networkLayer.(*layers.IPv4); ok {
						srcIP = ipv4.SrcIP.String()
						dstIP = ipv4.DstIP.String()
						protocol = ipv4.Protocol.String()
					} else if ipv6, ok := networkLayer.(*layers.IPv6); ok {
						srcIP = ipv6.SrcIP.String()
						dstIP = ipv6.DstIP.String()
						protocol = ipv6.NextHeader.String()
					}
				}

				transportLayer := packet.TransportLayer()
				if transportLayer != nil {
					if tcp, ok := transportLayer.(*layers.TCP); ok {
						protocol = "TCP"
						srcIP = fmt.Sprintf("%s:%d", srcIP, tcp.SrcPort)
						dstIP = fmt.Sprintf("%s:%d", dstIP, tcp.DstPort)
					} else if udp, ok := transportLayer.(*layers.UDP); ok {
						protocol = "UDP"
						srcIP = fmt.Sprintf("%s:%d", srcIP, udp.SrcPort)
						dstIP = fmt.Sprintf("%s:%d", dstIP, udp.DstPort)
					}
				}

				capturedPacket := anynetwork.CapturedPacket{
					Timestamp:   packet.Metadata().Timestamp.Format(time.RFC3339Nano),
					Source:      srcIP,
					Destination: dstIP,
					Protocol:    protocol,
					Length:      packet.Metadata().Length,
				}
				packetChannel <- capturedPacket
				runtime.EventsEmit(s.ctx, "packetCaptureEvent", capturedPacket) // Emit event to frontend
			}
		}
	}()

	return packetChannel, nil
}

// StopCapture stops the ongoing packet capture.
func (s *WindowsPacketCaptureService) StopCapture() error {
	if s.handle != nil {
		s.handle.Close()
		s.handle = nil
		return nil
	}
	return fmt.Errorf("no active capture to stop")
}
