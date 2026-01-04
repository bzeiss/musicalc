package ui

import (
	"fmt"
	"musicalc/internal/logic"
	sclres "musicalc/internal/logic/scl"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewDiapasonTab() fyne.CanvasObject {
	// Reference frequency binding
	refFreq := binding.NewString()
	_ = refFreq.Set("440")

	// Reference note binding (default A3 = MIDI 57)
	refNote := binding.NewString()
	_ = refNote.Set("A3")

	// Tuning system binding
	tuning := binding.NewString()
	_ = tuning.Set(sclres.DefaultScaleName)

	// Build reference note options (C0 to B9)
	noteNames := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	// Build tuning options from available scales (sorted alphabetically with default at top)
	tuningOptions := make([]string, 0, len(sclres.AvailableScales))
	for name := range sclres.AvailableScales {
		tuningOptions = append(tuningOptions, name)
	}
	sort.Strings(tuningOptions)

	// Move default tuning to the top of the list
	defaultIdx := -1
	for i, name := range tuningOptions {
		if name == sclres.DefaultScaleName {
			defaultIdx = i
			break
		}
	}
	if defaultIdx > 0 {
		// Remove default from its position and prepend to the list
		tuningOptions = append([]string{sclres.DefaultScaleName}, append(tuningOptions[:defaultIdx], tuningOptions[defaultIdx+1:]...)...)
	}

	// Reset function
	resetToDefaults := func() {
		_ = refNote.Set("A3")
		_ = refFreq.Set("440")
		_ = tuning.Set(sclres.DefaultScaleName)
	}

	// Reference note selector (simple entry, no longer needed for note selection)
	refNoteEntry := widget.NewEntry()
	refNoteEntry.SetText("A3")
	refNoteEntry.OnChanged = func(s string) {
		_ = refNote.Set(s)
	}

	// Frequency input
	freqInput := widget.NewEntry()
	freqInput.SetText("440")
	freqInput.PlaceHolder = "Frequency"
	freqInput.OnChanged = func(s string) {
		_ = refFreq.Set(s)
	}

	// Middle C convention selector (radio buttons)
	middleCRadio := widget.NewRadioGroup([]string{"C3", "C4"}, nil)
	middleCRadio.SetSelected("C3")
	middleCRadio.Horizontal = true
	middleCRadio.Required = true

	// Track previous Middle C selection to avoid adjusting on unchanged selections
	previousMiddleC := "C3"

	// Tuning selector (use SelectEntry for searchable dropdown)
	tuningSelect := widget.NewSelectEntry(tuningOptions)
	tuningSelect.SetText(sclres.DefaultScaleName)
	tuningSelect.OnChanged = func(selected string) {
		_ = tuning.Set(selected)
	}

	// Declare table variable first for use in reset button
	var table *widget.Table

	// Cache for performance - avoid binding.Get() on every cell render
	var cachedRefHz float64
	var cachedOctaveOffset int
	var cachedTuningName string

	// Update cache function
	updateCache := func() {
		refFreqVal, _ := refFreq.Get()
		tuningVal, _ := tuning.Get()

		// Determine octave offset based on Middle C setting
		if middleCRadio.Selected == "C3" {
			cachedOctaveOffset = 2 // C3 convention
		} else {
			cachedOctaveOffset = 1 // C4 convention
		}

		cachedRefHz = logic.ParseFloat(refFreqVal)
		cachedTuningName = tuningVal
	}

	// Initialize cache
	updateCache()

	// Reset button
	resetBtn := widget.NewButton("↻ Reset", func() {
		resetToDefaults()
		refNoteEntry.SetText("A3")
		freqInput.SetText("440")
		middleCRadio.SetSelected("C3")
		tuningSelect.SetText(sclres.DefaultScaleName)
		updateCache()
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
			// Calculate frequency using selected tuning
			result := logic.GetFrequency(midiNote, cachedRefHz, cachedTuningName)

			switch id.Col {
			case 0:
				l.SetText(fmt.Sprintf("%s%d", noteNames[midiNote%12], (midiNote/12)-cachedOctaveOffset))
			case 1:
				l.SetText(fmt.Sprintf("%.2f Hz", result.Frequency))
			case 2:
				// Display cents deviation from 12-TET
				if result.Cents >= 0 {
					l.SetText(fmt.Sprintf("+%.2f", result.Cents))
				} else {
					l.SetText(fmt.Sprintf("%.2f", result.Cents))
				}
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

	// Wrap table in responsive container with proportional column widths
	// Proportions: Note (23%), Frequency (28%), Cents (23%), MIDI (26%)
	responsiveTableWidget := NewResponsiveTable(table, []float32{0.23, 0.28, 0.23, 0.26}, 400, 20)

	// Add listener for Middle C radio buttons
	middleCRadio.OnChanged = func(s string) {
		// Only adjust reference if Middle C selection actually changed
		if s == previousMiddleC {
			return
		}

		// Adjust reference note octave when convention changes
		currentRef := refNoteEntry.Text
		if len(currentRef) >= 2 {
			// Parse current note
			var noteName string
			var octave int

			if len(currentRef) > 2 && currentRef[1] == '#' {
				noteName = currentRef[:2]
				fmt.Sscanf(currentRef[2:], "%d", &octave)
			} else {
				noteName = currentRef[:1]
				fmt.Sscanf(currentRef[1:], "%d", &octave)
			}

			// Adjust octave based on convention change
			if s == "C4" {
				// Switching from C3 to C4: increase octave by 1
				octave++
			} else {
				// Switching from C4 to C3: decrease octave by 1
				octave--
			}

			// Update reference note
			newRef := fmt.Sprintf("%s%d", noteName, octave)
			refNoteEntry.SetText(newRef)
			_ = refNote.Set(newRef)
		}

		// Update previous selection tracker
		previousMiddleC = s

		updateCache()
		table.Refresh()
	}

	// Add listeners to refresh table when inputs change
	refFreq.AddListener(binding.NewDataListener(func() {
		updateCache()
		table.Refresh()
	}))
	refNote.AddListener(binding.NewDataListener(func() {
		updateCache()
		table.Refresh()
	}))
	tuning.AddListener(binding.NewDataListener(func() {
		updateCache()
		table.Refresh()
	}))

	// Build UI layout
	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Note to Frequency", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			container.NewGridWithColumns(2,
				widget.NewLabel("Reference"),
				refNoteEntry,
			),
			container.NewGridWithColumns(2,
				widget.NewLabel("Frequency (Hz)"),
				freqInput,
			),
			container.NewGridWithColumns(2,
				widget.NewLabel("Middle C"),
				middleCRadio,
			),
			container.NewGridWithColumns(2,
				tuningSelect,
				resetBtn,
			),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		responsiveTableWidget,
	)
}
