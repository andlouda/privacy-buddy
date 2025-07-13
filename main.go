package main

import (
	"context"
	"embed"

	"privacy-buddy/backend"
	appsvc "privacy-buddy/backend/appsvc"
	anynetwork "privacy-buddy/backend/network"
	anynettools "privacy-buddy/backend/network/tools"
	platform_network "privacy-buddy/backend/platform/network"
	"privacy-buddy/backend/report"
	"privacy-buddy/backend/system"

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
	networkToolsSvc := anynettools.NewNetworkToolsService(tracerouteSvc)
	advancedNetworkToolsSvc := anynettools.GetAdvancedNetworkToolsService() // ✅ holt Singleton

	// ✅ Korrekte Initialisierung über Konstruktor

	err := wails.Run(&options.App{
		Title:  "privacy-buddy",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			// Dein ursprünglicher Startup-Call
			appsvcInstance.Startup(ctx)

			// 👇 Wichtig: Singleton bekommt seinen Context
			advancedNetworkToolsSvc.WailsInit(ctx)
		},
		Bind: []interface{}{
			appsvcInstance,
			systemSvc,
			&backend.SetupService{},
			networkSvc,
			&anynetwork.NetInfoService{},
			&anynetwork.PublicIPService{},
			reportSvc,
			networkToolsSvc,
			advancedNetworkToolsSvc, // ✅ Jetzt korrekt initialisiert
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
