package network

import (
	"io/ioutil"
	"net/http"
)

// PublicIPInfo contains the public IP and the interface used to determine it.
type PublicIPInfo struct {
	IPAddress     string `json:"ipAddress"`
	InterfaceName string `json:"interfaceName"`
}

type PublicIPService struct{}

// GetPublicIP retrieves only the public IP address as a string.
func (i *PublicIPService) GetPublicIP() (string, error) {
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

// GetPublicIPInfo retrieves the public IP and the name of the default network interface.
// This method relies on platform-specific implementations.
func (i *PublicIPService) GetPublicIPInfo() (*PublicIPInfo, error) {
	// The platform-specific implementation will be in get_default_interface_*.go files
	return getPublicIPInfo()
}
