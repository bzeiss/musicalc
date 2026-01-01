package ui

import (
	"fmt"
	"musicalc/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewSamplerTab() fyne.CanvasObject {
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

	// Sample rate selector
	sampleRateSelect := widget.NewSelect(sampleRates, func(selected string) {
		_ = sampleRate.Set(selected)
	})
	sampleRateSelect.SetSelected("44100")

	// Beats input
	beatsInput := widget.NewEntry()
	beatsInput.SetText("4")
	beatsInput.OnChanged = func(s string) {
		_ = beats.Set(s)
	}

	// Tempo values to display in table
	tempos := []float64{
		60, 70, 80, 90, 100, 110, 120, 130, 140, 150,
		160, 170, 180, 190, 200, 210, 220, 230, 240,
	}

	// Table displaying tempo calculations
	table := widget.NewTableWithHeaders(
		func() (int, int) { return len(tempos), 3 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.Alignment = fyne.TextAlignLeading

			srVal, _ := sampleRate.Get()
			beatsVal, _ := beats.Get()
			sr := logic.ParseFloat(srVal)
			bt := logic.ParseFloat(beatsVal)
			bpm := tempos[id.Row]

			res := logic.GetSampleData(sr, bpm, bt)

			switch id.Col {
			case 0:
				l.SetText(fmt.Sprintf("%.0f BPM", bpm))
			case 1:
				l.SetText(fmt.Sprintf("%d", res.Samples))
			case 2:
				l.SetText(fmt.Sprintf("%.2f ms", res.MS))
			}
		},
	)

	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		l := o.(*widget.Label)

		if id.Col == -1 {
			l.SetText("")
			return
		}

		l.TextStyle = fyne.TextStyle{Bold: true}
		l.Alignment = fyne.TextAlignLeading

		switch id.Col {
		case 0:
			l.SetText("Tempo")
		case 1:
			l.SetText("Samples")
		case 2:
			l.SetText("Length")
		}
	}

	table.ShowHeaderColumn = false

	table.SetColumnWidth(0, 100)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 120)

	sampleRate.AddListener(binding.NewDataListener(func() { table.Refresh() }))
	beats.AddListener(binding.NewDataListener(func() { table.Refresh() }))

	return container.NewBorder(
		container.NewVBox(
			container.NewGridWithColumns(2,
				widget.NewLabel("Sample Rate:"),
				sampleRateSelect,
			),
			container.NewGridWithColumns(2,
				widget.NewLabel("Beats:"),
				beatsInput,
			),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		table,
	)
}
