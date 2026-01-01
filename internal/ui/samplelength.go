package ui

import (
	"fmt"
	"musicalc/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewSampleLengthTab() fyne.CanvasObject {
	// Sample rate binding
	sampleRate := binding.NewString()
	_ = sampleRate.Set("44100")

	// Beats binding
	beats := binding.NewString()
	_ = beats.Set("4")

	// Common sample rates (including lo-fi)
	sampleRates := []string{
		"8000",   // Lo-fi, telephone quality
		"11025",  // Lo-fi, quarter CD quality
		"16000",  // Wideband audio
		"22050",  // Half CD quality
		"32000",  // MiniDv, video
		"44100",  // CD quality (standard)
		"48000",  // DVD, professional audio
		"88200",  // High-res audio
		"96000",  // High-res audio, Blu-ray
		"192000", // Ultra high-res
	}

	// Sample rate selector (callback will be set after calcFromTempo is defined)
	sampleRateSelect := widget.NewSelect(sampleRates, nil)
	sampleRateSelect.SetSelected("44100")

	// Beats input
	beatsInput := widget.NewEntry()

	// Bidirectional input/output fields
	samplesEntry := widget.NewEntry()
	msEntry := widget.NewEntry()
	tempoEntry := widget.NewEntry()

	// Flag to prevent circular updates
	updating := false

	// Declare calculation functions as variables first
	var calcFromTempo func()
	var calcFromSamples func()
	var calcFromMS func()

	// Reset function
	resetToDefaults := func() {
		updating = true
		_ = sampleRate.Set("44100")
		_ = beats.Set("4")
		tempoEntry.SetText("120.00")
		updating = false
		calcFromTempo()
	}

	// Calculate from Tempo (when Tempo is input)
	calcFromTempo = func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		srVal, _ := sampleRate.Get()
		sr := logic.ParseFloat(srVal)
		bt := logic.ParseFloat(beatsInput.Text)
		bpm := logic.ParseFloat(tempoEntry.Text)

		if bpm > 0 && sr > 0 {
			res := logic.GetSampleData(sr, bpm, bt)
			samplesEntry.SetText(fmt.Sprintf("%d", res.Samples))
			msEntry.SetText(fmt.Sprintf("%.2f", res.MS))
		}
	}

	// Calculate from Samples (when Samples is input)
	calcFromSamples = func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		srVal, _ := sampleRate.Get()
		sr := logic.ParseFloat(srVal)
		bt := logic.ParseFloat(beatsInput.Text)
		samples := logic.ParseFloat(samplesEntry.Text)

		if samples > 0 && sr > 0 && bt > 0 {
			ms := (samples / sr) * 1000.0
			bpm := (60.0 / (ms / 1000.0)) * bt
			msEntry.SetText(fmt.Sprintf("%.2f", ms))
			tempoEntry.SetText(fmt.Sprintf("%.2f", bpm))
		}
	}

	// Calculate from MS (when MS is input)
	calcFromMS = func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		srVal, _ := sampleRate.Get()
		sr := logic.ParseFloat(srVal)
		bt := logic.ParseFloat(beatsInput.Text)
		ms := logic.ParseFloat(msEntry.Text)

		if ms > 0 && sr > 0 && bt > 0 {
			samples := int((ms / 1000.0) * sr)
			bpm := (60.0 / (ms / 1000.0)) * bt
			samplesEntry.SetText(fmt.Sprintf("%d", samples))
			tempoEntry.SetText(fmt.Sprintf("%.2f", bpm))
		}
	}

	// Wire up change handlers
	tempoEntry.OnChanged = func(s string) { calcFromTempo() }
	samplesEntry.OnChanged = func(s string) { calcFromSamples() }
	msEntry.OnChanged = func(s string) { calcFromMS() }

	// Set sample rate and beats callbacks now that functions are defined
	sampleRateSelect.OnChanged = func(s string) {
		_ = sampleRate.Set(s)
		calcFromTempo()
	}
	beatsInput.SetText("4")
	beatsInput.OnChanged = func(s string) {
		_ = beats.Set(s)
		calcFromTempo()
	}

	// Initialize values on startup
	tempoEntry.SetText("120.00")
	calcFromTempo()

	// Reset button
	resetBtn := widget.NewButton("↻", func() {
		resetToDefaults()
		sampleRateSelect.SetSelected("44100")
		beatsInput.SetText("4")
	})

	return container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Sample Rate:"),
			sampleRateSelect,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Beats:"),
			beatsInput,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Reset:"),
			resetBtn,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Length in samples:"),
			samplesEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Length in ms:"),
			msEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Tempo:"),
			tempoEntry,
		),
	)
}
