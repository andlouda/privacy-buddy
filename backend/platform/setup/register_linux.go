//go:build linux

package setup

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func RegisterDesktopEntry(appName, comment, iconPath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	desktopDir := filepath.Join(usr.HomeDir, ".local/share/applications")
	err = os.MkdirAll(desktopDir, 0755)
	if err != nil {
		return err
	}

	desktopFile := filepath.Join(desktopDir, fmt.Sprintf("%s.desktop", appName))
	content := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Comment=%s
Exec=%s
Icon=%s
Terminal=false
Categories=Utility;Network;
StartupNotify=true
`, appName, comment, execPath, iconPath)

	return os.WriteFile(desktopFile, []byte(content), 0644)
}
