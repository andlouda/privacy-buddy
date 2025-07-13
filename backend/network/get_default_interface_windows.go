//go:build windows

package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"syscall"
	"unsafe"
)

// MIB_IPFORWARDROW defines an entry in the IP forwarding table.
// https://docs.microsoft.com/en-us/windows/win32/api/ipmib/ns-ipmib-mib_ipforwardrow
type MIB_IPFORWARDROW struct {
	ForwardDest      [4]byte
	ForwardMask      [4]byte
	ForwardPolicy    uint32
	ForwardNextHop   [4]byte
	ForwardIfIndex   uint32
	ForwardType      uint32
	ForwardProto     uint32
	ForwardAge       uint32
	ForwardNextHopAS uint32
	ForwardMetric1   uint32
	ForwardMetric2   uint32
	ForwardMetric3   uint32
	ForwardMetric4   uint32
	ForwardMetric5   uint32
}

func getPublicIPInfo() (*PublicIPInfo, error) {
	defaultInterface, err := getDefaultInterfaceNameWindows()
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

func getDefaultInterfaceNameWindows() (string, error) {
	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	getIpForwardTable := iphlpapi.NewProc("GetIpForwardTable")

	var buffer []byte
	bufferSize := uint32(0)

	// Get buffer size
	ret, _, _ := getIpForwardTable.Call(uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(&bufferSize)), 0)
	if syscall.Errno(ret) != syscall.ERROR_INSUFFICIENT_BUFFER {
		return "", fmt.Errorf("GetIpForwardTable failed to get buffer size: %v", syscall.Errno(ret))
	}

	buffer = make([]byte, bufferSize)

	// Get the table
	ret, _, _ = getIpForwardTable.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&bufferSize)), 0)
	if syscall.Errno(ret) != 0 { // 0 means NO_ERROR
		return "", fmt.Errorf("GetIpForwardTable failed: %v", syscall.Errno(ret))
	}

	// First 4 bytes are the number of entries
	numEntries := *(*uint32)(unsafe.Pointer(&buffer[0]))
	rows := (*[1 << 20]MIB_IPFORWARDROW)(unsafe.Pointer(&buffer[4]))[:numEntries]

	for _, row := range rows {
		// Default route has a destination of 0.0.0.0
		if row.ForwardDest == [4]byte{0, 0, 0, 0} {
			iface, err := net.InterfaceByIndex(int(row.ForwardIfIndex))
			if err != nil {
				continue // Try next default route if this one fails
			}
			return iface.Name, nil
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
