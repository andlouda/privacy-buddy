package main

import (
	"privacy-buddy/backend"
	appsvc "privacy-buddy/backend/appsvc"
	anynetwork "privacy-buddy/backend/network"
	"privacy-buddy/backend/network/tools"
	platform_network "privacy-buddy/backend/platform/network"
	"privacy-buddy/backend/report"
	"privacy-buddy/backend/system"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	appsvcInstance := appsvc.NewApp()
	systemSvc := &system.SystemService{}
	networkSvc := &anynetwork.NetworkDashboardService{}
	reportSvc := report.NewReportService(systemSvc, networkSvc)

	tracerouteSvc := platform_network.NewTracerouteService()
	networkToolsSvc := tools.NewNetworkToolsService(tracerouteSvc)
	advancedNetworkToolsSvc := tools.NewAdvancedNetworkToolsService()

	err := wails.Run(&options.App{
		Title:  "privacy-buddy",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        appsvcInstance.Startup,
		Bind: []interface{}{
			appsvcInstance,
			systemSvc,
			&backend.SetupService{},
			networkSvc,
			&anynetwork.NetInfoService{},
			&anynetwork.PublicIPService{},
			reportSvc,
			networkToolsSvc,
			advancedNetworkToolsSvc,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
