//go:build windows

package network

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"privacy-buddy/backend/network/tools"
)

// WindowsTracerouteService ist die Windows-spezifische Implementierung von TracerouteService.
type WindowsTracerouteService struct{}

// Traceroute führt einen tracert-Befehl auf Windows aus und parst die Ausgabe.
func (s *WindowsTracerouteService) Traceroute(host string) ([]tools.TracerouteHop, error) {
	cmd := exec.Command("tracert", host)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	fmt.Printf("Executing tracert command: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Tracert command failed. Stderr: %s\n", stderr.String())
		return nil, fmt.Errorf("Tracert-Befehl fehlgeschlagen: %w, %s", err, stderr.String())
	}

	fmt.Printf("Tracert command stdout: %s\n", out.String())
	return parseWindowsTracerouteOutput(out.String()), nil
}

// parseWindowsTracerouteOutput parst die Ausgabe des Windows-tracert-Befehls.
func parseWindowsTracerouteOutput(output string) []tools.TracerouteHop {
	hops := []tools.TracerouteHop{}
	lines := strings.Split(output, "\n")

	re := regexp.MustCompile(`^\s*(\d+)\s+([\d\s\*ms]+)\s+([^\s]+)(?:\s+\[([\d\.]+)\])?.*`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			n, _ := strconv.Atoi(matches[1])
			host := matches[3]
			address := matches[4] // Kann leer sein
			rtt := strings.TrimSpace(matches[2])

			if address == "" {
				address = host // Wenn keine IP in Klammern, ist der Host die IP
			}

			hops = append(hops, tools.TracerouteHop{
				N:       n,
				Host:    host,
				Address: address,
				RTT:     rtt,
			})
		}
	}
	return hops
}

// NewTracerouteService erstellt eine neue Instanz des WindowsTracerouteService.
func NewTracerouteService() tools.TracerouteService {
	return &WindowsTracerouteService{}
}

