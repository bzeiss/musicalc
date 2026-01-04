package ui

import (
	"fmt"
	"musicalc/internal/audio"
	"musicalc/internal/logic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewFrequencyToNoteTab() fyne.CanvasObject {
	// Input field
	frequencyEntry := widget.NewEntry()
	frequencyEntry.SetPlaceHolder("Frequency")
	frequencyEntry.SetText("440")

	// Middle C convention selector (radio buttons)
	middleCRadio := widget.NewRadioGroup([]string{"C3", "C4"}, nil)
	middleCRadio.SetSelected("C3")
	middleCRadio.Horizontal = true
	middleCRadio.Required = true

	// Output labels
	note100Label := widget.NewLabel("A4")
	cents100Label := widget.NewLabel("0")
	note50Label := widget.NewLabel("A4")
	cents50Label := widget.NewLabel("0")

	// Flag to prevent circular updates
	updating := false

	// Calculate and update all fields
	calculateFromFrequency := func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		freq := logic.ParseFloat(frequencyEntry.Text)
		if freq > 0 {
			// Determine octave offset based on Middle C setting
			octaveOffset := 1 // C4 convention (default)
			if middleCRadio.Selected == "C3" {
				octaveOffset = 2 // C3 convention
			}

			result := logic.FrequencyToNote(freq, octaveOffset)

			note100Label.SetText(result.Note100)

			// Format cents with sign
			cents100Sign := ""
			if result.Cents100 > 0 {
				cents100Sign = "+"
			}
			cents100Label.SetText(fmt.Sprintf("%s%d", cents100Sign, result.Cents100))

			note50Label.SetText(result.Note50)

			cents50Sign := ""
			if result.Cents50 > 0 {
				cents50Sign = "+"
			}
			cents50Label.SetText(fmt.Sprintf("%s%d", cents50Sign, result.Cents50))
		}
	}

	// Wire up change handlers
	frequencyEntry.OnChanged = func(s string) {
		calculateFromFrequency()
	}

	middleCRadio.OnChanged = func(s string) {
		calculateFromFrequency()
	}

	// Quick-select buttons
	c3Button := widget.NewButton("C3", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCRadio.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("C", 3, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	a3Button := widget.NewButton("A3", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCRadio.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("A", 3, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	c4Button := widget.NewButton("C4", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCRadio.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("C", 4, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	a4Button := widget.NewButton("A4", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCRadio.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("A", 4, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	// Play button - generates sine wave tone
	playButton := widget.NewButton("▶ Play", func() {
		freq := logic.ParseFloat(frequencyEntry.Text)
		if freq > 0 && freq < 20000 {
			go audio.PlayFrequency(freq, 3*time.Second)
		}
	})

	// Reset button
	resetButton := widget.NewButton("↻ Reset", func() {
		frequencyEntry.SetText("440")
		middleCRadio.SetSelected("C3")
	})

	// Initialize with default frequency (A4 = 440.00 Hz)
	calculateFromFrequency()

	return container.NewVBox(
		widget.NewLabelWithStyle("Frequency to Note", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Frequency (Hz)"),
			frequencyEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Quickselect"),
			container.NewHBox(c3Button, a3Button, c4Button, a4Button),
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Middle C"),
			middleCRadio,
		),
		container.NewGridWithColumns(2,
			playButton,
			resetButton,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Note"),
			note100Label,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Cents"),
			cents100Label,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Note (50 cents)"),
			note50Label,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Cents (50 cents)"),
			cents50Label,
		),
	)
}
