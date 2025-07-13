package report

import (
	"privacy-buddy/backend/network"
	"privacy-buddy/backend/system"
	"encoding/json"
	"os"
)

// ReportService sammelt Daten von anderen Diensten, um einen Diagnosebericht zu erstellen.
type ReportService struct {
	systemSvc  *system.SystemService
	networkSvc *network.NetworkDashboardService
}

// NewReportService erstellt eine neue Instanz des ReportService.
func NewReportService(systemSvc *system.SystemService, networkSvc *network.NetworkDashboardService) *ReportService {
	return &ReportService{
		systemSvc:  systemSvc,
		networkSvc: networkSvc,
	}
}

// ReportData strukturiert die gesammelten Diagnoseinformationen.
type ReportData struct {
	SystemInfo *system.SystemInfo `json:"systemInfo"`
	PublicIP   string             `json:"publicIP"`
	LocalIP    string             `json:"localIP"`
}

// GenerateReport sammelt alle relevanten Informationen und gibt sie als JSON-String zurück.
func (s *ReportService) GenerateReport() (string, error) {
	sysInfo, err := s.systemSvc.GetSystemInfo()
	if err != nil {
		// Auch wenn ein Teil fehlschlägt, können wir einen Teilbericht erstellen
		sysInfo = &system.SystemInfo{Username: "Fehler bei Abruf"}
	}

	data := &ReportData{
		SystemInfo: sysInfo,
		PublicIP:   s.networkSvc.GetPublicIP(),
		LocalIP:    s.networkSvc.GetLocalIP(),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// SaveReport speichert den generierten Bericht an dem vom Frontend übergebenen Pfad.
func (s *ReportService) SaveReport(reportData string, filePath string) error {
	return os.WriteFile(filePath, []byte(reportData), 0644)
}
