package ui

import (
	"fmt"
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
		return (octave+1)*12 + noteIndex
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

		// Dynamically filter options based on input
		if s == "" {
			refNoteEntry.SetOptions(refNoteOptions)
		} else {
			var filtered []string
			for _, opt := range refNoteOptions {
				if len(opt) >= len(s) && opt[:len(s)] == s {
					filtered = append(filtered, opt)
				}
			}
			if len(filtered) > 0 {
				refNoteEntry.SetOptions(filtered)
			}
		}
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

	// Table displaying all MIDI notes
	table = widget.NewTableWithHeaders(
		func() (int, int) { return 128, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.Alignment = fyne.TextAlignLeading

			// Get reference note and frequency
			refNoteVal, _ := refNote.Get()
			refFreqVal, _ := refFreq.Get()
			refMidi := getMidiNote(refNoteVal)
			refHz := logic.ParseFloat(refFreqVal)

			// Calculate frequency: freq = refFreq * 2^((thisMidi - refMidi)/12)
			// Using manual calculation since we need semitone steps
			semitones := float64(id.Row - refMidi)
			ratio := 1.0
			semitoneRatio := 1.0594630943592953 // 2^(1/12)

			if semitones > 0 {
				for i := 0; i < int(semitones); i++ {
					ratio *= semitoneRatio
				}
			} else if semitones < 0 {
				for i := 0; i < int(-semitones); i++ {
					ratio /= semitoneRatio
				}
			}
			hz := refHz * ratio

			switch id.Col {
			case 0:
				l.SetText(fmt.Sprintf("%s %d", noteNames[id.Row%12], (id.Row/12)-1))
			case 1:
				l.SetText(fmt.Sprintf("%.2f Hz", hz))
			case 2:
				l.SetText("+0.00")
			case 3:
				l.SetText(fmt.Sprintf("%d", id.Row))
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
	refFreq.AddListener(binding.NewDataListener(func() { table.Refresh() }))
	refNote.AddListener(binding.NewDataListener(func() { table.Refresh() }))
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
