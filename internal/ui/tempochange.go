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

func NewTempoChangeTab() fyne.CanvasObject {
	// Tempo bindings
	originalTempo := binding.NewString()
	_ = originalTempo.Set("140.00")

	// Input/output fields
	originalTempoEntry := widget.NewEntry()
	originalTempoEntry.SetText("140.00")

	newTempoEntry := widget.NewEntry()
	newTempoEntry.SetText("22.00")

	timeStretchEntry := widget.NewEntry()
	timeStretchEntry.SetText("636.36")

	// Read-only output fields (using Labels)
	tempoDeltaLabel := widget.NewLabel("-84.29 %")

	// Editable transpose inputs (standard notation)
	semitonesEntry := widget.NewEntry()
	semitonesEntry.SetText("-33")
	centsEntry := widget.NewEntry()
	centsEntry.SetText("96")

	// Read-only 50-cent notation outputs
	semitones50Label := widget.NewLabel("-32")
	cents50Label := widget.NewLabel("-4")

	// Flag to prevent circular updates
	updating := false

	// Calculate from new tempo
	calcFromNewTempo := func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		origTempo := logic.ParseFloat(originalTempoEntry.Text)
		newTempo := logic.ParseFloat(newTempoEntry.Text)

		if origTempo > 0 && newTempo > 0 {
			res := logic.CalculateFromNewTempo(origTempo, newTempo)
			timeStretchEntry.SetText(fmt.Sprintf("%.2f", res.TimeStretchPercent))

			sign := ""
			if res.TempoVariation > 0 {
				sign = "+"
			}
			tempoDeltaLabel.SetText(fmt.Sprintf("%s%.2f %%", sign, res.TempoVariation))

			semiSign := ""
			if res.Semitones > 0 {
				semiSign = "+"
			}
			semitonesEntry.SetText(fmt.Sprintf("%s%d", semiSign, res.Semitones))

			centsSign := ""
			if res.Cents > 0 {
				centsSign = "+"
			}
			centsEntry.SetText(fmt.Sprintf("%s%d", centsSign, res.Cents))

			semi50Sign := ""
			if res.Semitones50Cent > 0 {
				semi50Sign = "+"
			}
			semitones50Label.SetText(fmt.Sprintf("%s%d", semi50Sign, res.Semitones50Cent))

			cents50Sign := ""
			if res.Cents50Cent > 0 {
				cents50Sign = "+"
			}
			cents50Label.SetText(fmt.Sprintf("%s%d", cents50Sign, res.Cents50Cent))
		}
	}

	// Calculate from time stretch %
	calcFromTimeStretch := func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		origTempo := logic.ParseFloat(originalTempoEntry.Text)
		timeStretch := logic.ParseFloat(timeStretchEntry.Text)

		if origTempo > 0 && timeStretch > 0 {
			res := logic.CalculateFromTimeStretch(origTempo, timeStretch)
			newTempoEntry.SetText(fmt.Sprintf("%.2f", res.NewTempo))

			sign := ""
			if res.TempoVariation > 0 {
				sign = "+"
			}
			tempoDeltaLabel.SetText(fmt.Sprintf("%s%.2f %%", sign, res.TempoVariation))

			semiSign := ""
			if res.Semitones > 0 {
				semiSign = "+"
			}
			semitonesEntry.SetText(fmt.Sprintf("%s%d", semiSign, res.Semitones))

			centsSign := ""
			if res.Cents > 0 {
				centsSign = "+"
			}
			centsEntry.SetText(fmt.Sprintf("%s%d", centsSign, res.Cents))

			semi50Sign := ""
			if res.Semitones50Cent > 0 {
				semi50Sign = "+"
			}
			semitones50Label.SetText(fmt.Sprintf("%s%d", semi50Sign, res.Semitones50Cent))

			cents50Sign := ""
			if res.Cents50Cent > 0 {
				cents50Sign = "+"
			}
			cents50Label.SetText(fmt.Sprintf("%s%d", cents50Sign, res.Cents50Cent))
		}
	}

	// Calculate from transpose (semitones/cents inputs)
	calcFromTranspose := func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		origTempo := logic.ParseFloat(originalTempoEntry.Text)

		// Allow incomplete input (e.g., just "-" or "+") without updating
		semiText := semitonesEntry.Text
		centsText := centsEntry.Text
		if semiText == "" || semiText == "-" || semiText == "+" ||
			centsText == "" || centsText == "-" || centsText == "+" {
			return
		}

		semitones := int(logic.ParseFloat(semiText))
		cents := int(logic.ParseFloat(centsText))

		// Store original values to check if normalization changed them
		origSemitones := semitones
		origCents := cents

		// Normalize cents to [-50, 50] range by wrapping to adjacent semitones
		// E.g., 51 cents => +1 semitone, -49 cents
		for cents > 50 {
			cents -= 100
			semitones++
		}
		for cents < -50 {
			cents += 100
			semitones--
		}

		// Dynamic semitone bounds based on tempo limits (5.0 - 1000.0 BPM)
		// Calculate max/min semitones (as floats for precise boundary checking)
		if origTempo > 0 {
			const minTempo = 5.0
			const maxTempo = 999.99 // Slightly under 1000 to avoid rounding errors

			// Calculate exact semitone boundaries
			maxSemitonesFloat := 12.0 * math.Log2(maxTempo/origTempo)
			minSemitonesFloat := 12.0 * math.Log2(minTempo/origTempo)

			// Calculate total semitones including cents
			totalSemitones := float64(semitones) + float64(cents)/100.0

			// Clamp to tempo boundaries
			if totalSemitones > maxSemitonesFloat {
				// Clamp to max: use floor for semitones to match 50-cent notation
				semitones = int(math.Floor(maxSemitonesFloat))
				cents = int(math.Round((maxSemitonesFloat - float64(semitones)) * 100.0))
			}
			if totalSemitones < minSemitonesFloat {
				// Clamp to min: use floor for semitones to match 50-cent notation
				semitones = int(math.Floor(minSemitonesFloat))
				cents = int(math.Round((minSemitonesFloat - float64(semitones)) * 100.0))
			}
		}

		// Only update the UI if values changed after normalization
		if semitones != origSemitones {
			semiSign := ""
			if semitones > 0 {
				semiSign = "+"
			}
			semitonesEntry.SetText(fmt.Sprintf("%s%d", semiSign, semitones))
		}

		if cents != origCents {
			centSign := ""
			if cents > 0 {
				centSign = "+"
			}
			centsEntry.SetText(fmt.Sprintf("%s%d", centSign, cents))
		}

		if origTempo > 0 {
			res := logic.CalculateFromTranspose(origTempo, semitones, cents)
			newTempoEntry.SetText(fmt.Sprintf("%.2f", res.NewTempo))
			timeStretchEntry.SetText(fmt.Sprintf("%.2f", res.TimeStretchPercent))

			sign := ""
			if res.TempoVariation > 0 {
				sign = "+"
			}
			tempoDeltaLabel.SetText(fmt.Sprintf("%s%.2f %%", sign, res.TempoVariation))

			semi50Sign := ""
			if res.Semitones50Cent > 0 {
				semi50Sign = "+"
			}
			semitones50Label.SetText(fmt.Sprintf("%s%d", semi50Sign, res.Semitones50Cent))

			cents50Sign := ""
			if res.Cents50Cent > 0 {
				cents50Sign = "+"
			}
			cents50Label.SetText(fmt.Sprintf("%s%d", cents50Sign, res.Cents50Cent))
		}
	}

	// Reset function (defined after calculation functions so it can call them)
	resetToDefaults := func() {
		updating = true
		_ = originalTempo.Set("140.00")
		originalTempoEntry.SetText("140.00")
		newTempoEntry.SetText("22.00")
		updating = false
		// Trigger recalculation of all dependent fields
		calcFromNewTempo()
	}

	// Wire up change handlers (only for input fields)
	originalTempoEntry.OnChanged = func(s string) {
		_ = originalTempo.Set(s)
		calcFromNewTempo()
	}
	newTempoEntry.OnChanged = func(s string) { calcFromNewTempo() }
	timeStretchEntry.OnChanged = func(s string) { calcFromTimeStretch() }
	semitonesEntry.OnChanged = func(s string) { calcFromTranspose() }
	centsEntry.OnChanged = func(s string) { calcFromTranspose() }

	// Swap button to swap original and new tempo
	swapBtn := widget.NewButton("⇄", func() {
		origText := originalTempoEntry.Text
		newText := newTempoEntry.Text
		originalTempoEntry.SetText(newText)
		newTempoEntry.SetText(origText)
	})

	// Reset button
	resetBtn := widget.NewButton("↻", func() {
		resetToDefaults()
	})

	return container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Tempo:"),
			originalTempoEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("New tempo:"),
			newTempoEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Swap Tempo / New Tempo:"),
			swapBtn,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Reset:"),
			resetBtn,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Time stretching %:"),
			timeStretchEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Tempo Delta:"),
			tempoDeltaLabel,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Transpose semis:"),
			semitonesEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Semis for 50 cents notation:"),
			semitones50Label,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Transpose cents:"),
			centsEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Cents for 50 cents notation:"),
			cents50Label,
		),
	)
}
