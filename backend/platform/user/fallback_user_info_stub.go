//go:build !windows && !linux && !darwin

package platform

func GetUsername() string {
	return "unbekannt"
}
