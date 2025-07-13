package network

import (
	"context"
	"time"
)

// PacketCaptureService defines the interface for network packet capturing.
type PacketCaptureService interface {
	StartCapture(ctx context.Context, interfaceName string, bpfFilter string, duration time.Duration) (<-chan CapturedPacket, error)
	StopCapture() error
}

// NewPacketCaptureService creates a new instance of PacketCaptureService based on the operating system.
// This function is implemented in platform-specific files (e.g., packet_capture_service_darwin.go, etc.).
func NewPacketCaptureService() PacketCaptureService {
	// This function will be overridden by platform-specific implementations
	// using build tags. For example, on Linux, packet_capture_service_linux.go
	// will provide the actual implementation of NewPacketCaptureService.
	return nil // Should never be reached
}
