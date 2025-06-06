package main

import (
	"embed"
	_ "embed"
	"github.com/geff0319/galaxy3/bridge"
	"log"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	bridge.InitBridge(assets)
	bridge.InitMianWin()
	//bridge.InitWidgetsWin()
	//bridge.InitSystray()

	// Run the application. This blocks until the application has been exited.
	//bridge.InitScheduledTasks()
	//bridge.MqNotifyConsumer()
	//bridge.CreateHook()
	err := bridge.MainApp.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
