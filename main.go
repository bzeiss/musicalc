package main

import (
	"musicalc/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	myApp := app.NewWithID("com.musicalc")

	window := myApp.NewWindow("MusiCalc")
	window.Resize(fyne.NewSize(620, 750))

	// Load and set application icon (if available)
	// Icon should be at: icons/appicon.png
	if icon, err := fyne.LoadResourceFromPath("icons/appicon.png"); err == nil {
		window.SetIcon(icon)
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("Tempo to Delay", ui.NewTempoTab()),
		container.NewTabItem("Note to Frequency", ui.NewDiapasonTab()),
		container.NewTabItem("Sample Length", ui.NewSampleLengthTab()),
		container.NewTabItem("Tempo Change", ui.NewTempoChangeTab()),
	)

	tabs.SetTabLocation(container.TabLocationBottom) // Mobile ergonomic standard

	window.SetContent(tabs)
	window.ShowAndRun()
}
