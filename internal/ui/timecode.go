package ui

import (
	"fmt"
	"musicalc/internal/logic"
	"musicalc/internal/ui/widgets"
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

	// First timecode input (single field)
	timecode1Entry := widgets.NewTimecodeEntry(false)

	// Output label for first timecode (frame count only)
	timecode1Label := widget.NewLabel("(0f @ 30)")

	// Second timecode input (single field)
	timecode2Entry := widgets.NewTimecodeEntry(false)

	// Output label for second timecode (frame count only)
	timecode2Label := widget.NewLabel("(0f @ 30)")

	// Auto-focus Timecode 2 when Timecode 1 is complete
	timecode1Entry.OnComplete = func() {
		fyne.CurrentApp().Driver().CanvasForObject(timecode2Entry).Focus(timecode2Entry)
	}

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

		h1, m1, s1, f1 := timecode1Entry.GetComponents()

		format := logic.GetFPSFormat(fpsSelect.Selected)
		totalFrames := logic.TimecodeToFrames(h1, m1, s1, f1, format)
		result := logic.FramesToTimecode(totalFrames, format)

		fpsLabel := strings.Split(fpsSelect.Selected, " ")[0]
		timecode1Label.SetText(fmt.Sprintf("(%df @ %s)", result.TotalFrames, fpsLabel))
	}

	// Calculate second timecode from inputs
	calculateTimecode2 := func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		h2, m2, s2, f2 := timecode2Entry.GetComponents()

		format := logic.GetFPSFormat(fpsSelect.Selected)
		totalFrames := logic.TimecodeToFrames(h2, m2, s2, f2, format)
		result := logic.FramesToTimecode(totalFrames, format)

		fpsLabel := strings.Split(fpsSelect.Selected, " ")[0]
		timecode2Label.SetText(fmt.Sprintf("(%df @ %s)", result.TotalFrames, fpsLabel))
	}

	// Track previous FPS for conversion history
	var previousFPS string
	previousFPS = "30 fps"

	// Wire up change handlers
	timecode1Entry.OnChanged = func(s string) {
		calculateTimecode1()
	}

	timecode2Entry.OnChanged = func(s string) {
		calculateTimecode2()
	}

	fpsSelect.OnChanged = func(s string) {
		// Only do conversion if FPS actually changed and there's a non-zero timecode
		if previousFPS != "" && previousFPS != s {
			// Get current timecode 1 values
			h1, m1, s1, f1 := timecode1Entry.GetComponents()

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
	addButton := widget.NewButton("+", func() {
		h1, m1, s1, f1 := timecode1Entry.GetComponents()
		h2, m2, s2, f2 := timecode2Entry.GetComponents()

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
		timecode1Entry.SetComponents(result.Hours, result.Minutes, result.Seconds, result.Frames)
		timecode2Entry.SetText("")
		updating = false
		calculateTimecode1()
		calculateTimecode2()

		// Focus Timecode 2 for next input
		fyne.CurrentApp().Driver().CanvasForObject(timecode2Entry).Focus(timecode2Entry)
	})
	addButton.Importance = widget.HighImportance

	// Subtract operation
	subtractButton := widget.NewButton("-", func() {
		h1, m1, s1, f1 := timecode1Entry.GetComponents()
		h2, m2, s2, f2 := timecode2Entry.GetComponents()

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
		timecode1Entry.SetComponents(result.Hours, result.Minutes, result.Seconds, result.Frames)
		timecode2Entry.SetText("")
		updating = false
		calculateTimecode1()
		calculateTimecode2()

		// Focus Timecode 2 for next input
		fyne.CurrentApp().Driver().CanvasForObject(timecode2Entry).Focus(timecode2Entry)
	})
	subtractButton.Importance = widget.HighImportance

	// Keyboard shortcuts for +/- operations
	handleOperationKey := func(key fyne.KeyName) {
		switch key {
		case fyne.KeyEqual: // + key (Shift+= on US keyboards)
			addButton.OnTapped()
		case fyne.KeyMinus:
			subtractButton.OnTapped()
		}
	}
	timecode1Entry.OnOperationKey = handleOperationKey
	timecode2Entry.OnOperationKey = handleOperationKey

	// Reset operation
	resetButton := widget.NewButton("↻ Reset", func() {
		updating = true
		timecode1Entry.SetText("")
		timecode2Entry.SetText("")
		historyList = []string{}
		historyText.SetText("")
		updating = false
		calculateTimecode1()
		calculateTimecode2()

		// Focus Timecode 1 for next input
		fyne.CurrentApp().Driver().CanvasForObject(timecode1Entry).Focus(timecode1Entry)
	})

	// Clear History operation
	clearHistoryButton := widget.NewButton("Clear Hist.", func() {
		historyList = []string{}
		historyText.SetText("")
	})

	// Initialize
	calculateTimecode1()
	calculateTimecode2()

	// Create the tab container
	tabContainer := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Timecode Calculator", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			container.NewGridWithColumns(2,
				timecode1Entry,
				timecode1Label,
			),
			container.NewGridWithColumns(2,
				timecode2Entry,
				timecode2Label,
			),
			container.NewGridWithColumns(2,
				addButton,
				subtractButton,
			),
			widget.NewSeparator(),
			container.NewGridWithColumns(3,
				fpsSelect,
				clearHistoryButton,
				resetButton,
			),
			widget.NewSeparator(),
		), // top
		nil,           // bottom
		nil,           // left
		nil,           // right
		historyScroll, // center - will expand to fill remaining space
	)

	return tabContainer
}
