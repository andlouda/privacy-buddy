//go:build darwin

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	anynetwork "privacy-buddy/backend/network"
)

// DarwinPacketCaptureService provides macOS-specific implementation for packet capturing.
type DarwinPacketCaptureService struct {
	handle *pcap.Handle
}

// StartCapture starts capturing packets on the specified interface with an optional BPF filter and duration.
func (s *DarwinPacketCaptureService) StartCapture(ctx context.Context, interfaceName string, bpfFilter string, duration time.Duration) (<-chan anynetwork.CapturedPacket, error) {

	if s.handle != nil {
		return nil, fmt.Errorf("capture already in progress")
	}

	// Open the device for capturing
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.NextPacket)
	if err != nil {
		return nil, fmt.Errorf("failed to open device %s: %w", interfaceName, err)
	}
	s.handle = handle

	// Set BPF filter if provided
	if bpfFilter != "" {
		if err := s.handle.SetBPFFilter(bpfFilter);
		err != nil {
			s.handle.Close()
			s.handle = nil
			return nil, fmt.Errorf("failed to set BPF filter: %w", err)
		}
	}

	packetSource := gopacket.NewPacketSource(s.handle, s.handle.LinkType())
	packetChannel := make(chan CapturedPacket)

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
				// This is a simplified example, real implementation would parse layers
				packetChannel <- anynetwork.CapturedPacket{
					Timestamp: packet.Metadata().Timestamp.Format(time.RFC3339Nano),
					Length:    packet.Metadata().Length,
				}
			}
		}
	}()

	return packetChannel, nil
}

// StopCapture stops the ongoing packet capture.
func (s *DarwinPacketCaptureService) StopCapture() error {
	if s.handle != nil {
		s.handle.Close()
		s.handle = nil
		return nil
	}
	return fmt.Errorf("no active capture to stop")
}

