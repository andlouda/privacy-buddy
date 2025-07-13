//go:build windows

package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RegisterDesktopEntry(appName, comment, iconPath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	shortcutPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs", fmt.Sprintf("%s.lnk", appName))

	ps := fmt.Sprintf(`
$WScriptShell = New-Object -ComObject WScript.Shell
$Shortcut = $WScriptShell.CreateShortcut("%s")
$Shortcut.TargetPath = "%s"
$Shortcut.IconLocation = "%s"
$Shortcut.Save()
`, shortcutPath, execPath, iconPath)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	return cmd.Run()
}
