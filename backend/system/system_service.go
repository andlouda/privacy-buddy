package system

import (
	user "privacy-buddy/backend/platform/user"
	"os"
	"runtime"
)

type SystemService struct{}

type SystemInfo struct {
	Username   string `json:"username"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	Hostname   string `json:"hostname"`
	WorkingDir string `json:"workingDir"`
}

func (s *SystemService) GetSystemInfo() (*SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unbekannt"
	}

	wd, err := os.Getwd()
	if err != nil {
		wd = "unbekannt"
	}

	info := &SystemInfo{
		Username:   user.GetUsername(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Hostname:   hostname,
		WorkingDir: wd,
	}
	return info, nil
}
