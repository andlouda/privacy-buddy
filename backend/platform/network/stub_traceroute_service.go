//go:build !linux && !windows

package network

import (
	"fmt"

	"privacy-buddy/backend/network/tools"
)

// StubTracerouteService ist eine Platzhalter-Implementierung für nicht unterstützte Plattformen.
type StubTracerouteService struct{}

// Traceroute gibt einen Fehler zurück, da Traceroute auf dieser Plattform nicht unterstützt wird.
func (s *StubTracerouteService) Traceroute(host string) ([]tools.TracerouteHop, error) {
	return nil, fmt.Errorf("Traceroute wird auf dieser Plattform nicht unterstützt: %s", host)
}

// NewTracerouteService erstellt eine neue Instanz des StubTracerouteService.
func NewTracerouteService() tools.TracerouteService {
	return &StubTracerouteService{}
}
