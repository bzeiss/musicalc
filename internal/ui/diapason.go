package ui

import (
	"fmt"
	"math"
	"musicalc/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewDiapasonTab() fyne.CanvasObject {
	// Reference frequency binding
	refFreq := binding.NewString()
	_ = refFreq.Set("440.00")

	// Reference note binding (default A3 = MIDI 57)
	refNote := binding.NewString()
	_ = refNote.Set("A3")

	// Tuning system binding
	tuning := binding.NewString()
	_ = tuning.Set("Equal Temperament")

	// Build reference note options (C0 to B9)
	noteNames := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	var refNoteOptions []string
	for octave := 0; octave <= 9; octave++ {
		for _, note := range noteNames {
			refNoteOptions = append(refNoteOptions, fmt.Sprintf("%s%d", note, octave))
		}
	}

	// Calculate MIDI note number from note name
	getMidiNote := func(noteName string) int {
		if len(noteName) < 2 {
			return 69 // Default to A4
		}

		// Parse note name (e.g., "A3", "C#4")
		var note string
		var octave int

		// Handle sharp notes (e.g., C#, D#)
		if len(noteName) > 2 && noteName[1] == '#' {
			note = noteName[:2]
			fmt.Sscanf(noteName[2:], "%d", &octave)
		} else {
			note = noteName[:1]
			fmt.Sscanf(noteName[1:], "%d", &octave)
		}

		noteIndex := 0
		for i, n := range noteNames {
			if n == note {
				noteIndex = i
				break
			}
		}
		// Our table starts at C-2 = MIDI 0 (not C-1 = MIDI 0 as in standard)
		return (octave+2)*12 + noteIndex
	}

	// Reset function
	resetToDefaults := func() {
		_ = refNote.Set("A3")
		_ = refFreq.Set("440.00")
		_ = tuning.Set("Equal Temperament")
	}

	// Reference note selector with autocomplete (SelectEntry)
	refNoteEntry := widget.NewSelectEntry(refNoteOptions)
	refNoteEntry.SetText("A3")
	refNoteEntry.OnChanged = func(s string) {
		_ = refNote.Set(s)
	}

	// Frequency input
	freqInput := widget.NewEntry()
	freqInput.SetText("440.00")
	freqInput.OnChanged = func(s string) {
		_ = refFreq.Set(s)
	}

	// Tuning selector
	tuningSelect := widget.NewSelect([]string{"Equal Temperament"}, func(selected string) {
		_ = tuning.Set(selected)
	})
	tuningSelect.SetSelected("Equal Temperament")

	// Declare table variable first for use in reset button
	var table *widget.Table

	// Cache for performance - avoid binding.Get() on every cell render
	var cachedRefMidi int
	var cachedRefHz float64

	// Update cache function
	updateCache := func() {
		refNoteVal, _ := refNote.Get()
		refFreqVal, _ := refFreq.Get()
		cachedRefMidi = getMidiNote(refNoteVal)
		cachedRefHz = logic.ParseFloat(refFreqVal)
	}

	// Initialize cache
	updateCache()

	// Reset button
	resetBtn := widget.NewButton("↻", func() {
		resetToDefaults()
		refNoteEntry.SetText("A3")
		freqInput.SetText("440.00")
		tuningSelect.SetSelected("Equal Temperament")
		if table != nil {
			table.Refresh()
		}
	})

	// Table displaying all MIDI notes (C-2 to B8, MIDI 0-131)
	table = widget.NewTableWithHeaders(
		func() (int, int) { return 132, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.Alignment = fyne.TextAlignLeading

			// MIDI note starts at 0 (C-2)
			midiNote := id.Row

			// Use cached values for performance (updated only when bindings change)
			// Calculate frequency: freq = refFreq * 2^((thisMidi - refMidi)/12)
			semitones := float64(midiNote - cachedRefMidi)
			hz := cachedRefHz * math.Pow(2.0, semitones/12.0)

			switch id.Col {
			case 0:
				l.SetText(fmt.Sprintf("%s%d", noteNames[midiNote%12], (midiNote/12)-2))
			case 1:
				l.SetText(fmt.Sprintf("%.2f Hz", hz))
			case 2:
				l.SetText("+0.00")
			case 3:
				// Show both MIDI conventions: standard (C4=60) / alternative (C3=60)
				// MIDI values only go 0-127, show "-" for invalid values
				midiStandard := midiNote
				midiAlternative := midiNote + 12

				standardStr := fmt.Sprintf("%d", midiStandard)
				if midiStandard > 127 {
					standardStr = "-"
				}

				alternativeStr := fmt.Sprintf("%d", midiAlternative)
				if midiAlternative > 127 {
					alternativeStr = "-"
				}

				l.SetText(fmt.Sprintf("%s / %s", standardStr, alternativeStr))
			}
		},
	)

	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		l := o.(*widget.Label)

		// Hide row headers (left column)
		if id.Col == -1 {
			l.SetText("")
			return
		}

		l.TextStyle = fyne.TextStyle{Bold: true}
		l.Alignment = fyne.TextAlignLeading

		switch id.Col {
		case 0:
			l.SetText("Note")
		case 1:
			l.SetText("Frequency")
		case 2:
			l.SetText("Cents")
		case 3:
			l.SetText("MIDI")
		}
	}

	// Hide row header column
	table.ShowHeaderColumn = false

	table.SetColumnWidth(0, 80)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 80)
	table.SetColumnWidth(3, 80)

	// Add listeners to refresh table when inputs change
	refFreq.AddListener(binding.NewDataListener(func() {
		updateCache()
		table.Refresh()
	}))
	refNote.AddListener(binding.NewDataListener(func() {
		updateCache()
		table.Refresh()
	}))
	tuning.AddListener(binding.NewDataListener(func() { table.Refresh() }))

	// Build UI layout
	return container.NewBorder(
		container.NewVBox(
			container.NewGridWithColumns(2,
				widget.NewLabel("Reference:"),
				refNoteEntry,
			),
			container.NewGridWithColumns(2,
				widget.NewLabel("Frequency:"),
				freqInput,
			),
			container.NewGridWithColumns(2,
				widget.NewLabel("Tuning:"),
				tuningSelect,
			),
			container.NewGridWithColumns(2,
				widget.NewLabel("Reset:"),
				resetBtn,
			),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		table,
	)
}
