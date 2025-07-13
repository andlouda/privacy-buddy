//go:build linux

package platform

import "os/user"

func GetUsername() string {
	u, err := user.Current()
	if err != nil {
		return "unbekannt"
	}
	return u.Username
}
