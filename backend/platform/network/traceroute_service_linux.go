//go:build linux

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

// LinuxTracerouteService ist die Linux-spezifische Implementierung von TracerouteService.
type LinuxTracerouteService struct{}

// Traceroute führt den traceroute-Befehl auf Linux aus und parst die Ausgabe.
func (s *LinuxTracerouteService) Traceroute(host string) ([]tools.TracerouteHop, error) {
	cmd := exec.Command("traceroute", host)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	fmt.Printf("Executing traceroute command: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))

	if err := cmd.Run(); err != nil {
		fmt.Printf("Traceroute command failed. Stderr: %s\n", stderr.String())
		return nil, fmt.Errorf("Traceroute-Befehl fehlgeschlagen: %w — %s", err, stderr.String())
	}

	fmt.Printf("Traceroute command stdout: %s\n", out.String())
	return parseLinuxTracerouteOutput(out.String()), nil
}

// parseLinuxTracerouteOutput parst die Ausgabe und extrahiert Hop‑Nummer, Host, IP und RTT.
func parseLinuxTracerouteOutput(output string) []tools.TracerouteHop {
	hops := []tools.TracerouteHop{}
	lines := strings.Split(output, "\n")

	// Regex für typische Zeile: Hop‑Nr, Host, (IP), RTT oder *
	re := regexp.MustCompile(`^\s*(\d+)\s+([^\s]+)\s+\(([^\)]+)\)\s+([\d\.]+(?:\s*ms)?|\*).*`)

	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); len(matches) == 5 {
			n, _ := strconv.Atoi(matches[1])
			hostname := matches[2]
			addr := matches[3]
			rtt := matches[4]

			hops = append(hops, tools.TracerouteHop{
				N:       n,
				Host:    hostname,
				Address: addr,
				RTT:     rtt,
			})
		}
	}
	return hops
}

// NewTracerouteService erstellt eine neue Instanz des LinuxTracerouteService.
func NewTracerouteService() tools.TracerouteService {
	return &LinuxTracerouteService{}
}

