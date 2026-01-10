package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"musicalc/internal/logic"
)

type compactIconButton struct {
	widget.BaseWidget
	icon     fyne.Resource
	onTapped func()
}

type tightVBoxLayout struct{}

type cardLine1Layout struct{}

func (l *tightVBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	y := float32(0)
	for _, o := range objects {
		if o == nil || !o.Visible() {
			continue
		}
		h := o.MinSize().Height
		o.Move(fyne.NewPos(0, y))
		o.Resize(fyne.NewSize(size.Width, h))
		y += h
	}
}

func (l *tightVBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w := float32(0)
	h := float32(0)
	for _, o := range objects {
		if o == nil || !o.Visible() {
			continue
		}
		ms := o.MinSize()
		if ms.Width > w {
			w = ms.Width
		}
		h += ms.Height
	}
	return fyne.NewSize(w, h)
}

func (l *cardLine1Layout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	left := objects[0]
	right := objects[1]

	rightMS := right.MinSize()
	leftMS := left.MinSize()

	leftW := size.Width - rightMS.Width
	if leftW < 0 {
		leftW = 0
	}

	leftY := (size.Height - leftMS.Height) / 2
	if leftY < 0 {
		leftY = 0
	}
	left.Move(fyne.NewPos(0, leftY))
	left.Resize(fyne.NewSize(leftW, leftMS.Height))

	rightY := (size.Height - rightMS.Height) / 2
	if rightY < 0 {
		rightY = 0
	}
	right.Move(fyne.NewPos(size.Width-rightMS.Width, rightY))
	right.Resize(rightMS)
}

func (l *cardLine1Layout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(0, 0)
	}
	leftMS := objects[0].MinSize()
	rightMS := objects[1].MinSize()
	h := leftMS.Height
	if rightMS.Height > h {
		h = rightMS.Height
	}
	return fyne.NewSize(leftMS.Width+rightMS.Width, h)
}

func newCompactIconButton(icon fyne.Resource, onTapped func()) *compactIconButton {
	b := &compactIconButton{icon: icon, onTapped: onTapped}
	b.ExtendBaseWidget(b)
	return b
}

func (b *compactIconButton) MinSize() fyne.Size {
	return fyne.NewSize(16, 16)
}

func (b *compactIconButton) Tapped(*fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *compactIconButton) CreateRenderer() fyne.WidgetRenderer {
	img := canvas.NewImageFromResource(b.icon)
	img.FillMode = canvas.ImageFillContain
	return widget.NewSimpleRenderer(img)
}

// NewAlignmentDelayTab creates the Multi-Mic Alignment Delay calculator
func NewAlignmentDelayTab() fyne.CanvasObject {
	content, _, _ := NewAlignmentDelayTabWithExport()
	return content
}

func NewAlignmentDelayTabWithExport() (fyne.CanvasObject, func() (string, error), func(string) error) {
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
		"Hi-Hat", "Tom 1", "Tom 2", "Tom 3", "Floor Tom",
		"Overhead 1", "Overhead 2", "Room 1", "Room 2",
		"Vocal", "Guitar Amp", "Guitar DI", "Bass Amp", "Bass DI",
		"Piano 1", "Piano 2", "Strings", "Brass", "Ambience",
	}

	// Mic name as Select (dropdown only for consistent height)
	micNameSelect := widget.NewSelectEntry(commonMicPositions)
	micNameSelect.SetPlaceHolder("Name")

	// Target distance input
	targetDistEntry := widget.NewEntry()
	targetDistEntry.SetPlaceHolder("Dist. to Ref.")

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

	var refreshTable func()

	// Mobile-friendly card list (2 lines per mic)
	cardList := widget.NewList(
		func() int {
			return len(calc.Mics)
		},
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			nameLabel.Truncation = fyne.TextTruncateClip
			removeIcon := newCompactIconButton(theme.DeleteIcon(), nil)
			removeIconSlot := container.NewGridWrap(fyne.NewSize(22, 16), container.NewCenter(removeIcon))
			line1 := container.New(&cardLine1Layout{}, nameLabel, removeIconSlot)

			distLabel := widget.NewLabel("")
			distLabel.Truncation = fyne.TextTruncateClip
			delayLabel := widget.NewLabel("")
			delayLabel.Alignment = fyne.TextAlignTrailing
			delayLabel.Truncation = fyne.TextTruncateClip
			line2 := container.NewHBox(distLabel, layout.NewSpacer(), delayLabel)

			content := container.New(&tightVBoxLayout{}, line1, line2)
			padded := container.New(layout.NewCustomPaddedLayout(2, 2, 1, 1), content)

			bg := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
			bg.StrokeColor = theme.Color(theme.ColorNameSeparator)
			bg.StrokeWidth = 1
			bg.CornerRadius = 6

			return container.NewStack(bg, padded)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			if id < 0 || id >= len(calc.Mics) {
				return
			}
			mic := calc.Mics[id]
			outer := o.(*fyne.Container)
			padded := outer.Objects[1].(*fyne.Container)
			content := padded.Objects[0].(*fyne.Container)

			line1 := content.Objects[0].(*fyne.Container)
			nameLabel := line1.Objects[0].(*widget.Label)
			removeIconSlot := line1.Objects[1].(*fyne.Container)
			removeIconCenter := removeIconSlot.Objects[0].(*fyne.Container)
			removeIcon := removeIconCenter.Objects[0].(*compactIconButton)

			line2 := content.Objects[1].(*fyne.Container)
			distLabel := line2.Objects[0].(*widget.Label)
			delayLabel := line2.Objects[2].(*widget.Label)

			nameLabel.SetText(mic.Name)
			unit := targetUnitSelect.Selected
			distLabel.SetText(formatTrimFloat(logic.FromMeters(mic.DistanceMeters, unit), 3) + unit)
			if mic.IsBeyondReference {
				delayLabel.SetText("N/A")
			} else {
				delayLabel.SetText(formatTrimFloat(mic.DelayMS, 2) + fmt.Sprintf("ms / %d smp", mic.DelaySamples))
			}

			rowID := id
			removeIcon.onTapped = func() {
				if rowID >= 0 && rowID < len(calc.Mics) {
					calc.RemoveMicAt(rowID)
					if refreshTable != nil {
						refreshTable()
					}
				}
			}
		},
	)

	emptyMicsLabel := widget.NewLabelWithStyle("No microphones added yet", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	// Refresh data
	refreshTable = func() {
		calc.SetTemperature(logic.ParseFloat(tempEntry.Text), tempUnitSelect.Selected)
		calc.SetSampleRateLabel(sampleRateSelect.Selected)
		calc.SetReferenceDistance(logic.ParseFloat(refDistEntry.Text), refUnitSelect.Selected)
		calc.Recalculate()
		if len(calc.Mics) == 0 {
			emptyMicsLabel.Show()
			cardList.Hide()
		} else {
			emptyMicsLabel.Hide()
			cardList.Show()
		}
		cardList.Refresh()
	}

	exportCSV := func() (string, error) {
		if refreshTable != nil {
			refreshTable()
		}

		unit := targetUnitSelect.Selected
		roomTemp := logic.ParseFloat(tempEntry.Text)
		roomTempUnit := tempUnitSelect.Selected
		refDist := logic.ParseFloat(refDistEntry.Text)
		refDistUnit := refUnitSelect.Selected

		return logic.AlignmentDelayExportCSV(calc, unit, roomTemp, roomTempUnit, refDist, refDistUnit)
	}

	importCSV := func(csvData string) error {
		res, err := logic.AlignmentDelayImportCSV(csvData)
		if err != nil {
			return err
		}

		sampleRateLabel := fmt.Sprintf("%d", res.SampleRate)
		foundRate := false
		for _, opt := range sampleRateSelect.Options {
			if opt == sampleRateLabel {
				foundRate = true
				break
			}
		}
		if foundRate {
			sampleRateSelect.SetSelected(sampleRateLabel)
		} else {
			sampleRateSelect.SetSelected("48000")
		}

		tempUnit := strings.ToUpper(strings.TrimSpace(res.RoomTempUnit))
		if tempUnit == "F" {
			tempUnitSelect.SetSelected("F")
		} else {
			tempUnitSelect.SetSelected("C")
		}
		tempEntry.SetText(formatTrimFloat(res.RoomTemp, 2))

		refUnit := strings.ToLower(strings.TrimSpace(res.RefDistUnit))
		if refUnit == "ft" {
			refUnitSelect.SetSelected("ft")
		} else {
			refUnitSelect.SetSelected("m")
		}
		refDistEntry.SetText(formatTrimFloat(res.RefDist, 3))

		distUnit := strings.ToLower(strings.TrimSpace(res.DistUnit))
		if distUnit == "ft" {
			targetUnitSelect.SetSelected("ft")
		} else {
			targetUnitSelect.SetSelected("m")
		}

		calc.Mics = nil
		for _, mic := range res.Mics {
			u := strings.ToLower(strings.TrimSpace(mic.DistUnit))
			if u != "ft" {
				u = "m"
			}
			calc.AddMic(mic.Name, mic.Dist, u)
		}

		micNameSelect.SetText("")
		targetDistEntry.SetText("")
		if refreshTable != nil {
			refreshTable()
		}
		return nil
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

	// Recalculate when temperature or reference changes
	tempEntry.OnChanged = func(string) { refreshTable() }
	tempUnitSelect.OnChanged = func(string) { refreshTable() }
	refDistEntry.OnChanged = func(string) { refreshTable() }
	refUnitSelect.OnChanged = func(string) { refreshTable() }
	sampleRateSelect.OnChanged = func(string) { refreshTable() }
	targetUnitSelect.OnChanged = func(string) { refreshTable() }

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
		fixedWrap(tempUnitSelect, 90),
		tempEntry,
	)

	refGroup := container.NewBorder(
		nil, nil, nil,
		fixedWrap(refUnitSelect, 90),
		refDistEntry,
	)

	topRow := container.NewVBox(
		sampleRateSelect,
		tempGroup,
		refGroup,
	)

	// Emphasize the add button and make it wider
	addButton.Importance = widget.HighImportance

	fixedAddButton := container.NewGridWrap(
		fyne.NewSize(50, micNameSelect.MinSize().Height),
		addButton,
	)

	unitWidth := float32(70)
	if fyne.CurrentDevice().IsMobile() {
		unitWidth = 90
	}

	micNameRow := container.NewBorder(nil, nil, nil, nil, micNameSelect)

	micDistRight := container.NewHBox(
		fixedWrap(targetUnitSelect, unitWidth),
		fixedAddButton,
	)

	micDistRow := container.NewBorder(
		nil, nil, nil,
		micDistRight,
		targetDistEntry,
	)

	topControls := container.NewVBox(
		topRow,
		widget.NewSeparator(),
		micNameRow,
		micDistRow,
	)

	micsHeading := widget.NewLabelWithStyle("Microphones", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	micsHeader := container.NewVBox(micsHeading, container.NewCenter(emptyMicsLabel))
	micsSection := container.NewBorder(micsHeader, nil, nil, nil, cardList)

	content := container.NewBorder(
		topControls,
		nil, nil, nil,
		micsSection,
	)
	return content, exportCSV, importCSV
}
