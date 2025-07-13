package backend

import (
	setup "privacy-buddy/backend/platform/setup"
	"fmt"
	"os"
	"path/filepath"
)

type SetupService struct{}

func (s *SetupService) RegisterDesktopEntry() error {
	if isDevMode() {
		fmt.Println("🟡 Dev mode: Skipping desktop registration.")
		return nil
	}

	execDir, err := os.Executable()
	if err != nil {
		fmt.Println("❌ Could not resolve executable path:", err)
		return err
	}

	iconPath := filepath.Join(filepath.Dir(execDir), "logo-privacy-buddy.png")

	fmt.Println("🟢 Build mode: Registering desktop entry.")
	return setup.RegisterDesktopEntry("Privacy Buddy", "Cross-platform tool", iconPath)
}
