package ui

import (
	"fmt"
	"musicalc/internal/logic"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewTimecodeTab() fyne.CanvasObject {
	// FPS format selector
	fpsFormats := []string{}
	for _, format := range logic.FPSFormats {
		fpsFormats = append(fpsFormats, format.Name)
	}

	fpsSelect := widget.NewSelect(fpsFormats, nil)
	fpsSelect.SetSelected("30 fps")

	// First timecode inputs
	hours1Entry := widget.NewEntry()

	minutes1Entry := widget.NewEntry()

	seconds1Entry := widget.NewEntry()

	frames1Entry := widget.NewEntry()

	// Output label for first timecode (compact format)
	timecode1Label := widget.NewLabel("00:00:00:00 (0 frames @ 30)")

	// Second timecode inputs
	hours2Entry := widget.NewEntry()

	minutes2Entry := widget.NewEntry()

	seconds2Entry := widget.NewEntry()

	frames2Entry := widget.NewEntry()

	// Output label for second timecode (compact format)
	timecode2Label := widget.NewLabel("00:00:00:00 (0 frames @ 30)")

	// History display (multi-line entry for copy/paste support)
	historyText := widget.NewMultiLineEntry()
	historyText.Wrapping = fyne.TextWrapWord
	historyText.TextStyle.Monospace = true // Monospace for better alignment
	historyList := []string{}

	// Create scroll container early so operations can auto-scroll to bottom
	historyScroll := container.NewVScroll(historyText)

	// Flag to prevent circular updates
	updating := false

	// Calculate first timecode from inputs
	calculateTimecode1 := func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		h1 := int(logic.ParseFloat(hours1Entry.Text))
		m1 := int(logic.ParseFloat(minutes1Entry.Text))
		s1 := int(logic.ParseFloat(seconds1Entry.Text))
		f1 := int(logic.ParseFloat(frames1Entry.Text))

		if h1 < 0 {
			h1 = 0
		}
		if m1 < 0 {
			m1 = 0
		}
		if m1 > 99 {
			m1 = 99
		}
		if s1 < 0 {
			s1 = 0
		}
		if s1 > 99 {
			s1 = 99
		}
		if f1 < 0 {
			f1 = 0
		}
		if f1 > 99 {
			f1 = 99
		}

		format := logic.GetFPSFormat(fpsSelect.Selected)
		totalFrames := logic.TimecodeToFrames(h1, m1, s1, f1, format)
		result := logic.FramesToTimecode(totalFrames, format)

		fpsLabel := strings.Split(fpsSelect.Selected, " ")[0]
		timecode1Label.SetText(fmt.Sprintf("%s (%d frames @ %s)", result.Timecode, result.TotalFrames, fpsLabel))
	}

	// Calculate second timecode from inputs
	calculateTimecode2 := func() {
		if updating {
			return
		}

		h2 := int(logic.ParseFloat(hours2Entry.Text))
		m2 := int(logic.ParseFloat(minutes2Entry.Text))
		s2 := int(logic.ParseFloat(seconds2Entry.Text))
		f2 := int(logic.ParseFloat(frames2Entry.Text))

		if h2 < 0 {
			h2 = 0
		}
		if m2 < 0 {
			m2 = 0
		}
		if m2 > 99 {
			m2 = 99
		}
		if s2 < 0 {
			s2 = 0
		}
		if s2 > 99 {
			s2 = 99
		}
		if f2 < 0 {
			f2 = 0
		}
		if f2 > 99 {
			f2 = 99
		}

		format := logic.GetFPSFormat(fpsSelect.Selected)
		totalFrames := logic.TimecodeToFrames(h2, m2, s2, f2, format)
		result := logic.FramesToTimecode(totalFrames, format)

		fpsLabel := strings.Split(fpsSelect.Selected, " ")[0]
		timecode2Label.SetText(fmt.Sprintf("%s (%d frames @ %s)", result.Timecode, result.TotalFrames, fpsLabel))
	}

	// Track previous FPS for conversion history
	var previousFPS string
	previousFPS = "30 fps"

	// Wire up timecode change handlers
	hours1Entry.OnChanged = func(s string) { calculateTimecode1() }
	minutes1Entry.OnChanged = func(s string) { calculateTimecode1() }
	seconds1Entry.OnChanged = func(s string) { calculateTimecode1() }
	frames1Entry.OnChanged = func(s string) { calculateTimecode1() }

	hours2Entry.OnChanged = func(s string) { calculateTimecode2() }
	minutes2Entry.OnChanged = func(s string) { calculateTimecode2() }
	seconds2Entry.OnChanged = func(s string) { calculateTimecode2() }
	frames2Entry.OnChanged = func(s string) { calculateTimecode2() }

	fpsSelect.OnChanged = func(s string) {
		// Only do conversion if FPS actually changed and there's a non-zero timecode
		if previousFPS != "" && previousFPS != s {
			// Get current timecode 1 values
			h1 := int(logic.ParseFloat(hours1Entry.Text))
			m1 := int(logic.ParseFloat(minutes1Entry.Text))
			s1 := int(logic.ParseFloat(seconds1Entry.Text))
			f1 := int(logic.ParseFloat(frames1Entry.Text))

			// Only add conversion history if there's an actual timecode to convert
			if h1 > 0 || m1 > 0 || s1 > 0 || f1 > 0 {
				// Convert from previous FPS to new FPS by preserving TIMECODE NOTATION
				// MusicMath keeps H:M:S:F constant and recalculates frame count
				oldFormat := logic.GetFPSFormat(previousFPS)
				newFormat := logic.GetFPSFormat(s)

				// Get frames in old format
				oldFrames := logic.TimecodeToFrames(h1, m1, s1, f1, oldFormat)

				// Preserve the timecode notation (H:M:S:F) and recalculate frame count with new FPS
				// The H:M:S:F values stay the same, only total frames changes
				newFrames := logic.TimecodeToFrames(h1, m1, s1, f1, newFormat)

				oldTC := fmt.Sprintf("%02d:%02d:%02d:%02d", h1, m1, s1, f1)
				newTC := oldTC // Timecode notation stays the same
				oldFpsLabel := strings.Split(previousFPS, " ")[0]
				newFpsLabel := strings.Split(s, " ")[0]

				conversionEntry := fmt.Sprintf("  %s (%df) @%s\n= %s (%df) @%s",
					oldTC, oldFrames, oldFpsLabel, newTC, newFrames, newFpsLabel)
				historyList = append(historyList, conversionEntry)
				historyText.SetText(strings.Join(historyList, "\n\n"))
				historyText.Refresh()
				historyScroll.ScrollToBottom()

				// Timecode H:M:S:F values don't change, only recalculate display with new FPS
				// No need to update entry fields since H:M:S:F stay the same
			}
		}
		previousFPS = s
		calculateTimecode1()
		calculateTimecode2()
	}

	// Add operation
	addButton := widget.NewButton("Add", func() {
		h1 := int(logic.ParseFloat(hours1Entry.Text))
		m1 := int(logic.ParseFloat(minutes1Entry.Text))
		s1 := int(logic.ParseFloat(seconds1Entry.Text))
		f1 := int(logic.ParseFloat(frames1Entry.Text))

		h2 := int(logic.ParseFloat(hours2Entry.Text))
		m2 := int(logic.ParseFloat(minutes2Entry.Text))
		s2 := int(logic.ParseFloat(seconds2Entry.Text))
		f2 := int(logic.ParseFloat(frames2Entry.Text))

		// Clamp values
		if h1 < 0 {
			h1 = 0
		}
		if m1 < 0 {
			m1 = 0
		}
		if m1 > 99 {
			m1 = 99
		}
		if s1 < 0 {
			s1 = 0
		}
		if s1 > 99 {
			s1 = 99
		}
		if f1 < 0 {
			f1 = 0
		}
		if f1 > 99 {
			f1 = 99
		}

		if h2 < 0 {
			h2 = 0
		}
		if m2 < 0 {
			m2 = 0
		}
		if m2 > 99 {
			m2 = 99
		}
		if s2 < 0 {
			s2 = 0
		}
		if s2 > 99 {
			s2 = 99
		}
		if f2 < 0 {
			f2 = 0
		}
		if f2 > 99 {
			f2 = 99
		}

		format := logic.GetFPSFormat(fpsSelect.Selected)
		result := logic.AddTimecodes(h1, m1, s1, f1, h2, m2, s2, f2, format)

		// Add to history
		tc1 := fmt.Sprintf("%02d:%02d:%02d:%02d", h1, m1, s1, f1)
		tc2 := fmt.Sprintf("%02d:%02d:%02d:%02d", h2, m2, s2, f2)
		frames1 := logic.TimecodeToFrames(h1, m1, s1, f1, format)
		frames2 := logic.TimecodeToFrames(h2, m2, s2, f2, format)

		fpsLabel := strings.Split(fpsSelect.Selected, " ")[0]
		historyEntry := fmt.Sprintf("  %s (%df)\n+ %s (%df)\n= %s (%df) @%s",
			tc1, frames1, tc2, frames2, result.Timecode, result.TotalFrames, fpsLabel)
		historyList = append(historyList, historyEntry)
		historyText.SetText(strings.Join(historyList, "\n\n"))
		historyText.Refresh()
		historyScroll.ScrollToBottom()

		// Update first timecode with result and reset second timecode
		updating = true
		hours1Entry.SetText(strconv.Itoa(result.Hours))
		minutes1Entry.SetText(strconv.Itoa(result.Minutes))
		seconds1Entry.SetText(strconv.Itoa(result.Seconds))
		frames1Entry.SetText(strconv.Itoa(result.Frames))
		hours2Entry.SetText("")
		minutes2Entry.SetText("")
		seconds2Entry.SetText("")
		frames2Entry.SetText("")
		updating = false
		calculateTimecode1()
		calculateTimecode2()
	})

	// Subtract operation
	subtractButton := widget.NewButton("Subtract", func() {
		h1 := int(logic.ParseFloat(hours1Entry.Text))
		m1 := int(logic.ParseFloat(minutes1Entry.Text))
		s1 := int(logic.ParseFloat(seconds1Entry.Text))
		f1 := int(logic.ParseFloat(frames1Entry.Text))

		h2 := int(logic.ParseFloat(hours2Entry.Text))
		m2 := int(logic.ParseFloat(minutes2Entry.Text))
		s2 := int(logic.ParseFloat(seconds2Entry.Text))
		f2 := int(logic.ParseFloat(frames2Entry.Text))

		// Clamp values
		if h1 < 0 {
			h1 = 0
		}
		if m1 < 0 {
			m1 = 0
		}
		if m1 > 99 {
			m1 = 99
		}
		if s1 < 0 {
			s1 = 0
		}
		if s1 > 99 {
			s1 = 99
		}
		if f1 < 0 {
			f1 = 0
		}
		if f1 > 99 {
			f1 = 99
		}

		if h2 < 0 {
			h2 = 0
		}
		if m2 < 0 {
			m2 = 0
		}
		if m2 > 99 {
			m2 = 99
		}
		if s2 < 0 {
			s2 = 0
		}
		if s2 > 99 {
			s2 = 99
		}
		if f2 < 0 {
			f2 = 0
		}
		if f2 > 99 {
			f2 = 99
		}

		format := logic.GetFPSFormat(fpsSelect.Selected)
		result := logic.SubtractTimecodes(h1, m1, s1, f1, h2, m2, s2, f2, format)

		// Add to history
		tc1 := fmt.Sprintf("%02d:%02d:%02d:%02d", h1, m1, s1, f1)
		tc2 := fmt.Sprintf("%02d:%02d:%02d:%02d", h2, m2, s2, f2)
		frames1 := logic.TimecodeToFrames(h1, m1, s1, f1, format)
		frames2 := logic.TimecodeToFrames(h2, m2, s2, f2, format)

		fpsLabel := strings.Split(fpsSelect.Selected, " ")[0]
		historyEntry := fmt.Sprintf("  %s (%df)\n- %s (%df)\n= %s (%df) @%s",
			tc1, frames1, tc2, frames2, result.Timecode, result.TotalFrames, fpsLabel)
		historyList = append(historyList, historyEntry)
		historyText.SetText(strings.Join(historyList, "\n\n"))
		historyText.Refresh()
		historyScroll.ScrollToBottom()

		// Update first timecode with result and reset second timecode
		updating = true
		hours1Entry.SetText(strconv.Itoa(result.Hours))
		minutes1Entry.SetText(strconv.Itoa(result.Minutes))
		seconds1Entry.SetText(strconv.Itoa(result.Seconds))
		frames1Entry.SetText(strconv.Itoa(result.Frames))
		hours2Entry.SetText("")
		minutes2Entry.SetText("")
		seconds2Entry.SetText("")
		frames2Entry.SetText("")
		updating = false
		calculateTimecode1()
		calculateTimecode2()
	})

	// Reset operation
	resetButton := widget.NewButton("Reset", func() {
		updating = true
		hours1Entry.SetText("")
		minutes1Entry.SetText("")
		seconds1Entry.SetText("")
		frames1Entry.SetText("")
		hours2Entry.SetText("")
		minutes2Entry.SetText("")
		seconds2Entry.SetText("")
		frames2Entry.SetText("")
		historyList = []string{}
		historyText.SetText("")
		updating = false
		calculateTimecode1()
		calculateTimecode2()
	})

	// Clear History operation
	clearHistoryButton := widget.NewButton("Clear History", func() {
		historyList = []string{}
		historyText.SetText("")
	})

	// Initialize
	calculateTimecode1()
	calculateTimecode2()

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Timecode 1:"),
			container.NewGridWithColumns(4,
				widget.NewLabel("Hours:"),
				hours1Entry,
				widget.NewLabel("Minutes:"),
				minutes1Entry,
			),
			container.NewGridWithColumns(4,
				widget.NewLabel("Seconds:"),
				seconds1Entry,
				widget.NewLabel("Frames:"),
				frames1Entry,
			),
			timecode1Label,
			widget.NewSeparator(),
			widget.NewLabel("Timecode 2:"),
			container.NewGridWithColumns(4,
				widget.NewLabel("Hours:"),
				hours2Entry,
				widget.NewLabel("Minutes:"),
				minutes2Entry,
			),
			container.NewGridWithColumns(4,
				widget.NewLabel("Seconds:"),
				seconds2Entry,
				widget.NewLabel("Frames:"),
				frames2Entry,
			),
			timecode2Label,
			widget.NewSeparator(),
			container.NewGridWithColumns(2,
				widget.NewLabel("FPS Format:"),
				fpsSelect,
			),
			widget.NewSeparator(),
			container.NewGridWithColumns(4,
				addButton,
				subtractButton,
				resetButton,
				clearHistoryButton,
			),
			widget.NewSeparator(),
			widget.NewLabel("History:"),
		), // top
		nil,           // bottom
		nil,           // left
		nil,           // right
		historyScroll, // center - will expand to fill remaining space
	)
}
