//go:build darwin

package network

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func getPublicIPInfo() (*PublicIPInfo, error) {
	defaultInterface, err := getDefaultInterfaceNameDarwin()
	if err != nil {
		// Fallback to just getting the IP if interface lookup fails
		ip, ipErr := fetchPublicIP()
		return &PublicIPInfo{IPAddress: ip, InterfaceName: ""}, ipErr
	}

	ip, err := fetchPublicIP()
	if err != nil {
		return nil, err
	}

	return &PublicIPInfo{IPAddress: ip, InterfaceName: defaultInterface}, nil
}

func getDefaultInterfaceNameDarwin() (string, error) {
	cmd := exec.Command("route", "-n", "get", "default")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "interface:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1], nil
			}
		}
	}

	return "", fmt.Errorf("default network interface not found")
}

func fetchPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}
