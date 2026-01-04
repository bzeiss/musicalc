package main

import (
	_ "embed"
	"musicalc/internal/ui"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

//go:embed VERSION
var version string

func main() {
	myApp := app.NewWithID("com.musicalc")

	// Apply custom theme only on desktop, use default theme on mobile
	if !fyne.CurrentDevice().IsMobile() {
		myApp.Settings().SetTheme(ui.NewCustomTheme())
	}

	// Read version and create window title
	ver := strings.TrimSpace(version)
	windowTitle := "MusiCalc"
	if ver != "" {
		windowTitle = "MusiCalc v" + ver
	}
	window := myApp.NewWindow(windowTitle)
	window.Resize(fyne.NewSize(450, 650))
	window.SetFixedSize(true)

	// Load and set application icon (if available)
	// Icon should be at: icons/appicon.png
	if icon, err := fyne.LoadResourceFromPath("icons/appicon.png"); err == nil {
		window.SetIcon(icon)
	}

	// Create icon-only tabs using bundled resources (no text to save space on mobile)
	timecodeTab := container.NewTabItem("", ui.NewTimecodeTab())
	timecodeTab.Icon = ui.ResourceTimecodeSvg

	tempoTab := container.NewTabItem("", ui.NewTempoTab())
	tempoTab.Icon = ui.ResourceDelaySvg

	note2freqTab := container.NewTabItem("", ui.NewDiapasonTab())
	note2freqTab.Icon = ui.ResourceNote2freqSvg

	freq2noteTab := container.NewTabItem("", ui.NewFrequencyToNoteTab())
	freq2noteTab.Icon = ui.ResourceFreq2noteSvg

	sampleLengthTab := container.NewTabItem("", ui.NewSampleLengthTab())
	sampleLengthTab.Icon = ui.ResourceSamplelengthSvg

	tempoChangeTab := container.NewTabItem("", ui.NewTempoChangeTab())
	tempoChangeTab.Icon = ui.ResourceTempochangeSvg

	tabs := container.NewAppTabs(
		timecodeTab,
		tempoTab,
		note2freqTab,
		freq2noteTab,
		sampleLengthTab,
		tempoChangeTab,
	)

	tabs.SetTabLocation(container.TabLocationBottom) // Mobile ergonomic standard

	window.SetContent(tabs)
	window.ShowAndRun()
}
