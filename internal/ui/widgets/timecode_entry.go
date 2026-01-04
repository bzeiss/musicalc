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
}

func NewTimecodeEntry(threeDigits bool) *TimecodeEntry {
	e := &TimecodeEntry{ThreeDigitFrames: threeDigits}
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

// TypedRune handles inserting numbers at the current cursor position
func (e *TimecodeEntry) TypedRune(r rune) {
	// Handle +/- operation keys
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

	if r < '0' || r > '9' {
		return
	}

	// Get raw digits (without colons and spaces)
	raw := strings.ReplaceAll(strings.ReplaceAll(e.Text, ":", ""), " ", "")
	limit := 8
	if e.ThreeDigitFrames {
		limit = 9
	}

	if len(raw) < limit {
		// Calculate position in raw string by counting digits (not colons or spaces) before cursor
		rawPos := 0
		for i := 0; i < e.CursorColumn && i < len(e.Text); i++ {
			if e.Text[i] != ':' && e.Text[i] != ' ' {
				rawPos++
			}
		}

		// Insert digit into raw string
		rawRunes := []rune(raw)
		newRaw := append(rawRunes[:rawPos], append([]rune{r}, rawRunes[rawPos:]...)...)

		// Format with colons
		formatted := e.format(string(newRaw))
		e.SetText(formatted)

		// Calculate new cursor position (rawPos+1 digit, plus colons before it)
		newCursorPos := 0
		digitCount := 0
		for i := 0; i < len(formatted); i++ {
			if formatted[i] != ':' {
				digitCount++
				if digitCount > rawPos+1 {
					break
				}
			}
			newCursorPos++
		}

		e.CursorColumn = newCursorPos
		e.Refresh()

		// Trigger OnComplete if all digits are entered
		if len(newRaw) == limit && e.OnComplete != nil {
			e.OnComplete()
		}
	}
}

// TypedKey allows standard navigation keys to work normally
func (e *TimecodeEntry) TypedKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyEqual, fyne.KeyMinus:
		// Handle +/- operation keys
		if e.OnOperationKey != nil {
			e.OnOperationKey(k.Name)
		}
		return
	case fyne.KeyBackspace:
		pos := e.CursorColumn
		if pos == 0 {
			return
		}

		// Get raw digits only
		raw := strings.ReplaceAll(strings.ReplaceAll(e.Text, ":", ""), " ", "")
		if len(raw) == 0 {
			return
		}

		// Calculate position in raw string
		rawPos := 0
		for i := 0; i < pos && i < len(e.Text); i++ {
			if e.Text[i] != ':' && e.Text[i] != ' ' {
				rawPos++
			}
		}

		if rawPos > 0 {
			// Remove digit before cursor in raw string
			rawRunes := []rune(raw)
			newRaw := append(rawRunes[:rawPos-1], rawRunes[rawPos:]...)

			// Format and update
			formatted := e.format(string(newRaw))
			e.SetText(formatted)

			// Calculate new cursor position
			newCursorPos := 0
			digitCount := 0
			for i := 0; i < len(formatted); i++ {
				if formatted[i] != ':' && formatted[i] != ' ' {
					digitCount++
					if digitCount >= rawPos {
						break
					}
				}
				newCursorPos++
			}
			e.CursorColumn = newCursorPos
			e.Refresh()
		}
	case fyne.KeyLeft, fyne.KeyRight, fyne.KeyUp, fyne.KeyDown:
		// Let the standard widget handle navigation
		e.Entry.TypedKey(k)
	default:
		e.Entry.TypedKey(k)
	}
}

func (e *TimecodeEntry) format(raw string) string {
	// Always maintain HH:MM:SS:FF structure with spaces for empty positions
	limit := 8
	if e.ThreeDigitFrames {
		limit = 9
	}

	// Pad raw string with spaces to full length
	padded := raw
	for len(padded) < limit {
		padded += " "
	}

	// Format as HH:MM:SS:FF (or HH:MM:SS:FFF)
	if e.ThreeDigitFrames {
		return fmt.Sprintf("%c%c:%c%c:%c%c:%c%c%c",
			padded[0], padded[1], padded[2], padded[3],
			padded[4], padded[5], padded[6], padded[7], padded[8])
	}
	return fmt.Sprintf("%c%c:%c%c:%c%c:%c%c",
		padded[0], padded[1], padded[2], padded[3],
		padded[4], padded[5], padded[6], padded[7])
}

// GetComponents extracts hours, minutes, seconds, and frames from the formatted timecode
// Returns (hours, minutes, seconds, frames)
func (e *TimecodeEntry) GetComponents() (int, int, int, int) {
	parts := strings.Split(e.Text, ":")
	if len(parts) < 4 {
		return 0, 0, 0, 0
	}

	// Replace spaces with zeros for parsing
	h, _ := strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(parts[0]), " ", "0"))
	m, _ := strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(parts[1]), " ", "0"))
	s, _ := strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(parts[2]), " ", "0"))
	f, _ := strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(parts[3]), " ", "0"))

	return h, m, s, f
}

// SetComponents sets the timecode from individual components
func (e *TimecodeEntry) SetComponents(hours, minutes, seconds, frames int) {
	if e.ThreeDigitFrames {
		e.SetText(fmt.Sprintf("%02d:%02d:%02d:%03d", hours, minutes, seconds, frames))
	} else {
		e.SetText(fmt.Sprintf("%02d:%02d:%02d:%02d", hours, minutes, seconds, frames))
	}
}
