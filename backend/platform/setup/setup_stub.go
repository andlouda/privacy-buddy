//go:build !linux && !windows

package setup

func RegisterDesktopEntry(appName, comment, iconPath string) error {
	return nil // no-op for unsupported platforms
}
