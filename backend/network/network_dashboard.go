package network

import (
	"context"
	"fmt"
	"time"
)

type NetworkDashboardService struct {
	publicIPService PublicIPService
	localIPService  NetInfoService
	packetCaptureService PacketCaptureService
	networkInterfaceService NetworkInterfaceService
	arpCacheService ARPCacheService
	networkConnectionService NetworkConnectionService
}

// Wails ruft diese Methode nach dem Binden automatisch auf
func (s *NetworkDashboardService) Init() {
	s.publicIPService = PublicIPService{}
	s.localIPService = NetInfoService{}
	s.packetCaptureService = NewPacketCaptureService()
	s.networkInterfaceService = NewNetworkInterfaceService()
	s.arpCacheService = NewARPCacheService()
	s.networkConnectionService = NewNetworkConnectionService()
}

func (s *NetworkDashboardService) GetPublicIP() string {
	ip, err := s.publicIPService.GetPublicIP()
	if err != nil {
		return "Fehler beim Abrufen"
	}
	return ip
}

func (s *NetworkDashboardService) GetLocalIP() string {
	ip, err := s.localIPService.GetLocalIP()
	if err != nil {
		return "unbekannt"
	}
	return ip
}

// StartPacketCapture starts capturing packets on the specified interface.
func (s *NetworkDashboardService) StartPacketCapture(interfaceName string, bpfFilter string, durationStr string) (bool, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return false, fmt.Errorf("invalid duration format: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	// Store cancel function to stop capture later if needed
	// For simplicity, not storing it here, but in a real app, you'd manage this.

	packetChan, err := s.packetCaptureService.StartCapture(ctx, interfaceName, bpfFilter, duration)
	if err != nil {
		cancel()
		return false, fmt.Errorf("failed to start packet capture: %w", err)
	}

	// In a real Wails app, you'd send packets to the frontend via events.
	// For now, just consume them to prevent goroutine leak.
	go func() {
		for packet := range packetChan {
			// Process packet, e.g., send to frontend via Wails runtime.EventsEmit
			_ = packet // Suppress unused warning
		}
		cancel() // Ensure context is cancelled when channel closes
	}()

	return true, nil
}

// StopPacketCapture stops the ongoing packet capture.
func (s *NetworkDashboardService) StopPacketCapture() error {
	return s.packetCaptureService.StopCapture()
}

// GetNetworkInterfaces lists all network interfaces.
func (s *NetworkDashboardService) GetNetworkInterfaces() ([]NetworkInterface, error) {
	return s.networkInterfaceService.ListInterfaces()
}

// GetARPEntries retrieves all ARP cache entries.
func (s *NetworkDashboardService) GetARPEntries() ([]ARPEntry, error) {
	return s.arpCacheService.GetARPEntries()
}

// GetNetworkConnections retrieves all active network connections.
func (s *NetworkDashboardService) GetNetworkConnections() ([]NetworkConnection, error) {
	return s.networkConnectionService.GetConnections()
}