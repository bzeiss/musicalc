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
	OnComplete       func() // Called when all digits are entered
}

func NewTimecodeEntry(threeDigits bool) *TimecodeEntry {
	e := &TimecodeEntry{ThreeDigitFrames: threeDigits}
	e.ExtendBaseWidget(e)
	e.updatePlaceholder()
	return e
}

func (e *TimecodeEntry) updatePlaceholder() {
	if e.ThreeDigitFrames {
		e.PlaceHolder = "00:00:00:000"
	} else {
		e.PlaceHolder = "00:00:00:00"
	}
}

func (e *TimecodeEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

// TypedRune handles inserting numbers at the current cursor position
func (e *TimecodeEntry) TypedRune(r rune) {
	if r < '0' || r > '9' {
		return
	}

	// Get raw digits (without colons)
	raw := strings.ReplaceAll(e.Text, ":", "")
	limit := 8
	if e.ThreeDigitFrames {
		limit = 9
	}

	if len(raw) < limit {
		// Calculate position in raw string by counting non-colon chars before cursor
		rawPos := 0
		for i := 0; i < e.CursorColumn && i < len(e.Text); i++ {
			if e.Text[i] != ':' {
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
	case fyne.KeyBackspace:
		pos := e.CursorColumn
		if pos == 0 {
			return
		}

		text := []rune(e.Text)
		// If deleting a colon, delete the number before it instead
		shift := 1
		if text[pos-1] == ':' {
			shift = 2
		}

		if pos-shift >= 0 {
			newText := append(text[:pos-shift], text[pos:]...)
			formatted := e.format(strings.ReplaceAll(string(newText), ":", ""))
			e.SetText(formatted)
			e.CursorColumn = pos - shift
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
	var result strings.Builder
	digitCount := 0
	for _, char := range raw {
		if digitCount > 0 && (digitCount == 2 || digitCount == 4 || digitCount == 6) {
			result.WriteRune(':')
		}
		result.WriteRune(char)
		digitCount++
	}
	return result.String()
}

// GetComponents extracts hours, minutes, seconds, and frames from the formatted timecode
// Returns (hours, minutes, seconds, frames)
func (e *TimecodeEntry) GetComponents() (int, int, int, int) {
	parts := strings.Split(e.Text, ":")
	if len(parts) < 4 {
		return 0, 0, 0, 0
	}

	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	s, _ := strconv.Atoi(parts[2])
	f, _ := strconv.Atoi(parts[3])

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
