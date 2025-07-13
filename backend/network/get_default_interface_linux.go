//go:build linux

package network

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

func getPublicIPInfo() (*PublicIPInfo, error) {
	defaultInterface, err := getDefaultInterfaceNameLinux()
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

func getDefaultInterfaceNameLinux() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		// The default route is the one with destination 00000000
		if len(fields) >= 2 && fields[1] == "00000000" {
			return fields[0], nil
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
