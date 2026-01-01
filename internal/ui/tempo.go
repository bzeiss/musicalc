package ui

import (
	"fmt"
	"musicalc/internal/logic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewTempoTab() fyne.CanvasObject {
	bpm := binding.NewString()
	_ = bpm.Set("120.00")
	input := widget.NewEntry()
	input.SetText("120.00")
	input.OnChanged = func(s string) {
		_ = bpm.Set(s)
	}

	notes := []struct {
		Name string
		Mult float64
	}{
		{"1/1", 4.0},
		{"1/1D", 6.0},
		{"1/1T", 8.0 / 3.0},
		{"1/2", 2.0},
		{"1/2D", 3.0},
		{"1/2T", 4.0 / 3.0},
		{"1/4", 1.0},
		{"1/4D", 1.5},
		{"1/4T", 2.0 / 3.0},
		{"1/8", 0.5},
		{"1/8D", 0.75},
		{"1/8T", 1.0 / 3.0},
		{"1/16", 0.25},
		{"1/16D", 0.375},
		{"1/16T", 1.0 / 6.0},
		{"1/32", 0.125},
		{"1/32D", 0.1875},
		{"1/32T", 1.0 / 12.0},
	}

	table := widget.NewTableWithHeaders(
		func() (int, int) { return len(notes), 3 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.Alignment = fyne.TextAlignLeading

			val, _ := bpm.Get()
			f := logic.ParseFloat(val)
			res := logic.GetTempoData(f, notes[id.Row].Mult)

			switch id.Col {
			case 0:
				l.SetText(notes[id.Row].Name)
			case 1:
				l.SetText(fmt.Sprintf("%.2f ms", res.DelayMS))
			case 2:
				l.SetText(fmt.Sprintf("%.2f Hz", res.ModHz))
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
			l.SetText("Length")
		case 1:
			l.SetText("Delay")
		case 2:
			l.SetText("Modulation")
		}
	}

	// Hide row header column
	table.ShowHeaderColumn = false

	table.SetColumnWidth(0, 80)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 120)

	bpm.AddListener(binding.NewDataListener(func() { table.Refresh() }))

	return container.NewBorder(
		container.NewVBox(
			input,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		table,
	)
}
