package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"musicalc/internal/logic"
)

// NewAlignmentDelayTab creates the Multi-Mic Alignment Delay calculator
func NewAlignmentDelayTab() fyne.CanvasObject {
	// Temperature input
	tempEntry := widget.NewEntry()
	//tempEntry.SetText("22")
	tempEntry.SetPlaceHolder("Room Temp.")

	tempUnitSelect := widget.NewSelect([]string{"C", "F"}, nil)
	tempUnitSelect.SetSelected("C")

	// Reference distance input
	refDistEntry := widget.NewEntry()
	//refDistEntry.SetText("5")
	refDistEntry.SetPlaceHolder("Ref. Dist.")

	refUnitSelect := widget.NewSelect([]string{"m", "ft"}, nil)
	refUnitSelect.SetSelected("m")

	// Sample rate dropdown
	sampleRateSelect := widget.NewSelect([]string{
		"44100",
		"48000",
		"88200",
		"96000",
		"192000",
	}, nil)
	sampleRateSelect.SetSelected("48000")

	// Common microphone positions for dropdown
	commonMicPositions := []string{
		"Kick In", "Kick Out", "Snare Top", "Snare Bottom",
		"Hi-Hat", "Tom 1", "Tom 2", "Floor Tom",
		"Overhead L", "Overhead R", "Room L", "Room R",
		"Vocal", "Guitar Amp", "Bass Amp",
		"Piano L", "Piano R", "Strings", "Brass", "Ambience",
	}

	// Mic name as Select (dropdown only for consistent height)
	micNameSelect := widget.NewSelectEntry(commonMicPositions)
	micNameSelect.SetPlaceHolder("Name")

	// Target distance input
	targetDistEntry := widget.NewEntry()
	targetDistEntry.SetPlaceHolder("Distance to Ref.")

	targetUnitSelect := widget.NewSelect([]string{"m", "ft"}, nil)
	targetUnitSelect.SetSelected("m")
	targetUnitSelect.Resize(fyne.NewSize(40, targetUnitSelect.MinSize().Height))

	calc := logic.NewAlignmentDelayCalculator()

	formatTrimFloat := func(v float64, decimals int) string {
		s := fmt.Sprintf("%.*f", decimals, v)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		if s == "" || s == "-0" {
			return "0"
		}
		return s
	}

	// Table to display microphones with sticky header
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(calc.Mics), 4
		},
		func() fyne.CanvasObject {
			l := widget.NewLabel("")
			l.Truncation = fyne.TextTruncateClip
			return l
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			if id.Row >= len(calc.Mics) {
				return
			}
			label := cell.(*widget.Label)
			mic := calc.Mics[id.Row]
			switch id.Col {
			case 0:
				label.Alignment = fyne.TextAlignLeading
				label.SetText(mic.Name)
			case 1:
				label.Alignment = fyne.TextAlignLeading
				unit := targetUnitSelect.Selected
				label.SetText(formatTrimFloat(logic.FromMeters(mic.DistanceMeters, unit), 3) + unit)
			case 2:
				label.Alignment = fyne.TextAlignLeading
				if mic.IsBeyondReference {
					label.SetText("N/A")
					break
				}
				label.SetText(formatTrimFloat(mic.DelayMS, 2) + fmt.Sprintf("ms / %d smp", mic.DelaySamples))
			case 3:
				label.Alignment = fyne.TextAlignCenter
				label.SetText("🗑️")
			}
			label.TextStyle = fyne.TextStyle{Bold: false}
		},
	)

	// Configure sticky header
	table.CreateHeader = func() fyne.CanvasObject {
		l := widget.NewLabel("")
		l.Truncation = fyne.TextTruncateClip
		return l
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		label := o.(*widget.Label)
		if id.Col == -1 {
			label.SetText("")
			return
		}
		label.TextStyle = fyne.TextStyle{Bold: true}
		if id.Col == 3 {
			label.Alignment = fyne.TextAlignCenter
			label.SetText("🛠️")
			return
		}
		label.Alignment = fyne.TextAlignLeading
		headers := []string{"Name", "Distance", "Delay", "Actions"}
		label.SetText(headers[id.Col])
	}

	table.SetColumnWidth(0, 120)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 60)

	// Refresh table data
	refreshTable := func() {
		calc.SetTemperature(logic.ParseFloat(tempEntry.Text), tempUnitSelect.Selected)
		calc.SetSampleRateLabel(sampleRateSelect.Selected)
		calc.SetReferenceDistance(logic.ParseFloat(refDistEntry.Text), refUnitSelect.Selected)
		calc.Recalculate()
		table.Refresh()
	}

	// Add microphone button with emphasis
	addButton := widget.NewButton("+", func() {
		// Button handler
		name := strings.TrimSpace(micNameSelect.Text)
		if name == "" {
			name = fmt.Sprintf("Mic %d", len(calc.Mics)+1)
		}

		dist := logic.ParseFloat(targetDistEntry.Text)
		if dist == 0 {
			return
		}
		calc.AddMic(name, dist, targetUnitSelect.Selected)

		micNameSelect.SetText("")
		targetDistEntry.SetText("")
		refreshTable()
	})

	// Handle table clicks for remove
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Col == 3 && id.Row < len(calc.Mics) {
			// Remove microphone
			calc.RemoveMicAt(id.Row)
			refreshTable()
		}
		table.UnselectAll()
	}

	// Recalculate when temperature or reference changes
	tempEntry.OnChanged = func(string) { refreshTable() }
	tempUnitSelect.OnChanged = func(string) { refreshTable() }
	refDistEntry.OnChanged = func(string) { refreshTable() }
	refUnitSelect.OnChanged = func(string) { refreshTable() }
	sampleRateSelect.OnChanged = func(string) { refreshTable() }
	targetUnitSelect.OnChanged = func(string) { refreshTable() }

	// Create responsive table wrapper for proper column sizing
	responsiveTableWidget := NewResponsiveTable(
		table,
		[]float32{0.30, 0.20, 0.40, 0.10}, // Column proportions: Name, Distance, Delay, Remove
		100,                               // min width
		60,                                // padding
	)

	// Compact layout with minimal widths
	tempEntry.Resize(fyne.NewSize(80, tempEntry.MinSize().Height))
	refDistEntry.Resize(fyne.NewSize(80, refDistEntry.MinSize().Height))
	targetDistEntry.Resize(fyne.NewSize(80, targetDistEntry.MinSize().Height))

	fixedWrap := func(obj fyne.CanvasObject, w float32) fyne.CanvasObject {
		return container.NewGridWrap(fyne.NewSize(w, obj.MinSize().Height), obj)
	}

	// Top section: Temperature, Reference, and Sample Rate on same row

	tempGroup := container.NewBorder(
		nil, nil, nil,
		fixedWrap(tempUnitSelect, 60),
		tempEntry,
	)

	refGroup := container.NewBorder(
		nil, nil, nil,
		fixedWrap(refUnitSelect, 70),
		refDistEntry,
	)

	tempRefRow := container.NewGridWithColumns(2,
		tempGroup,
		refGroup,
	)

	topRow := container.NewBorder(
		nil, nil,
		fixedWrap(sampleRateSelect, 100), nil,
		tempRefRow,
	)

	// Emphasize the add button and make it wider
	addButton.Importance = widget.HighImportance

	fixedAddButton := container.NewGridWrap(
		fyne.NewSize(50, micNameSelect.MinSize().Height),
		addButton,
	)

	micRow := container.NewBorder(
		nil, nil, nil,
		fixedAddButton,
		container.NewBorder(nil, nil, nil, fixedWrap(targetUnitSelect, 70),
			container.NewGridWithColumns(2,
				micNameSelect,
				targetDistEntry,
			),
		),
	)

	// Use Border layout to make table stretch vertically
	content := container.NewBorder(
		container.NewVBox(
			topRow,
			widget.NewSeparator(),
			micRow,
			// widget.NewSeparator(),
		),
		nil, nil, nil,
		responsiveTableWidget,
	)

	return content
}
