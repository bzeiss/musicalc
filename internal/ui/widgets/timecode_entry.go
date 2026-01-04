package widgets

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type TimecodeEntry struct {
	widget.Entry
	ThreeDigitFrames bool
	OnComplete       func()                 // Called when all digits are entered
	OnOperationKey   func(key fyne.KeyName) // Called when +/- keys are pressed
	rawInput         string                 // Raw input buffer (digits and dots)
}

func NewTimecodeEntry(threeDigits bool) *TimecodeEntry {
	e := &TimecodeEntry{
		ThreeDigitFrames: threeDigits,
		rawInput:         "",
	}
	e.ExtendBaseWidget(e)

	// Use monospace font for consistent character width
	e.TextStyle.Monospace = true

	// Set placeholder to show expected format
	if threeDigits {
		e.PlaceHolder = "HH:MM:SS:FFF"
	} else {
		e.PlaceHolder = "HH:MM:SS:FF"
	}

	return e
}

func (e *TimecodeEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// TypedRune handles right-justified timecode entry with dot shorthand
func (e *TimecodeEntry) TypedRune(r rune) {
	// Handle +/- as operation keys
	if r == '+' || r == '-' {
		if e.OnOperationKey != nil {
			if r == '+' {
				e.OnOperationKey(fyne.KeyEqual)
			} else {
				e.OnOperationKey(fyne.KeyMinus)
			}
		}
		return
	}

	// Handle dot (.) or comma (,) as shorthand for "00"
	// Both work for locale compatibility (comma on European keyboards)
	if r == '.' || r == ',' {
		maxDigits := 8
		if e.ThreeDigitFrames {
			maxDigits = 9
		}
		// Only add "00" if it doesn't exceed max length
		if len(e.rawInput)+2 <= maxDigits {
			e.rawInput += "00"
			e.updateDisplay()
		}
		return
	}

	// Handle digits
	if r >= '0' && r <= '9' {
		// Right-justified input: digits accumulate from right (frames → seconds → minutes → hours)
		maxDigits := 8
		if e.ThreeDigitFrames {
			maxDigits = 9
		}

		if len(e.rawInput) < maxDigits {
			e.rawInput += string(r)
			e.updateDisplay()

			// Trigger OnComplete if all digits are entered
			if len(e.rawInput) >= maxDigits && e.OnComplete != nil {
				e.OnComplete()
			}
		}
	}
}

// TypedKey handles backspace and operation keys
func (e *TimecodeEntry) TypedKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyEqual, fyne.KeyMinus:
		// Handle +/- operation keys
		if e.OnOperationKey != nil {
			e.OnOperationKey(k.Name)
		}
		return
	// Period and comma are handled in TypedRune, not here
	case fyne.KeyBackspace:
		if len(e.rawInput) > 0 {
			// Remove last character from raw input
			e.rawInput = e.rawInput[:len(e.rawInput)-1]
			e.updateDisplay()
		}
		return
	case fyne.KeyReturn, fyne.KeyEnter:
		// Confirm entry - trigger operation or completion
		if e.OnComplete != nil {
			e.OnComplete()
		}
		return
	default:
		// Ignore other navigation keys in right-justified mode
		return
	}
}

// updateDisplay converts right-justified input to formatted timecode display
func (e *TimecodeEntry) updateDisplay() {
	// Right-justified: interpret input from right to left
	// Example: "15" → 00:00:00:15, "100" → 00:00:01:00, "10000" → 00:01:00:00

	if len(e.rawInput) == 0 {
		e.SetText("")
		e.CursorColumn = 0
		return
	}

	// Pad with leading zeros to full length
	maxDigits := 8
	if e.ThreeDigitFrames {
		maxDigits = 9
	}

	// Ensure rawInput doesn't exceed maxDigits (safety check)
	if len(e.rawInput) > maxDigits {
		e.rawInput = e.rawInput[len(e.rawInput)-maxDigits:]
	}

	padded := strings.Repeat("0", maxDigits-len(e.rawInput)) + e.rawInput

	// Extract components from padded string (HHMMSSFF or HHMMSSFFF)
	var formatted string
	if e.ThreeDigitFrames {
		// HH:MM:SS:FFF
		formatted = fmt.Sprintf("%s:%s:%s:%s",
			padded[0:2], padded[2:4], padded[4:6], padded[6:9])
	} else {
		// HH:MM:SS:FF
		formatted = fmt.Sprintf("%s:%s:%s:%s",
			padded[0:2], padded[2:4], padded[4:6], padded[6:8])
	}

	e.SetText(formatted)
	e.CursorColumn = len(formatted)
}

// GetComponents extracts hours, minutes, seconds, and frames from the right-justified input
// Returns (hours, minutes, seconds, frames)
func (e *TimecodeEntry) GetComponents() (int, int, int, int) {
	if len(e.rawInput) == 0 {
		return 0, 0, 0, 0
	}

	// Pad with leading zeros
	maxDigits := 8
	if e.ThreeDigitFrames {
		maxDigits = 9
	}
	padded := strings.Repeat("0", maxDigits-len(e.rawInput)) + e.rawInput

	// Extract components
	var h, m, s, f int
	if e.ThreeDigitFrames {
		h, _ = strconv.Atoi(padded[0:2])
		m, _ = strconv.Atoi(padded[2:4])
		s, _ = strconv.Atoi(padded[4:6])
		f, _ = strconv.Atoi(padded[6:9])
	} else {
		h, _ = strconv.Atoi(padded[0:2])
		m, _ = strconv.Atoi(padded[2:4])
		s, _ = strconv.Atoi(padded[4:6])
		f, _ = strconv.Atoi(padded[6:8])
	}

	return h, m, s, f
}

// SetComponents sets the timecode from individual components
func (e *TimecodeEntry) SetComponents(hours, minutes, seconds, frames int) {
	// Convert components to right-justified raw input
	if e.ThreeDigitFrames {
		e.rawInput = fmt.Sprintf("%02d%02d%02d%03d", hours, minutes, seconds, frames)
	} else {
		e.rawInput = fmt.Sprintf("%02d%02d%02d%02d", hours, minutes, seconds, frames)
	}
	// Remove leading zeros from raw input for cleaner display
	e.rawInput = strings.TrimLeft(e.rawInput, "0")
	if e.rawInput == "" {
		e.rawInput = "0"
	}
	e.updateDisplay()
}

// SetText overrides to clear raw input buffer
func (e *TimecodeEntry) SetText(text string) {
	if text == "" {
		e.rawInput = ""
	}
	e.Entry.SetText(text)
}
