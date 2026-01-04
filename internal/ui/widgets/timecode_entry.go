package widgets

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type TimecodeEntry struct {
	widget.Entry
	ThreeDigitFrames bool
	OnComplete       func()                 // Called when all digits are entered
	OnOperationKey   func(key fyne.KeyName) // Called when +/- keys are pressed
	fields           []string               // Field values [frames, seconds, minutes, hours] - right to left
}

func NewTimecodeEntry(threeDigits bool) *TimecodeEntry {
	e := &TimecodeEntry{
		ThreeDigitFrames: threeDigits,
		fields:           []string{""},
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

// TypedRune handles right-justified timecode entry with dot as field separator
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

	// Handle dot (.) or comma (,) as field separator
	// Locks current field and moves to next field left (frames -> seconds -> minutes -> hours)
	// Both work for locale compatibility (comma on European keyboards)
	if r == '.' || r == ',' {
		// Only add new field if current field has content and we haven't reached max fields
		if len(e.fields) < 4 && e.fields[0] != "" {
			// Insert empty field at beginning, pushing current values right
			e.fields = append([]string{""}, e.fields...)
			e.updateDisplay()
		}
		return
	}

	// Handle digits
	if r >= '0' && r <= '9' {
		// Determine max digits for current field
		maxDigits := 2
		if len(e.fields) == 1 && e.ThreeDigitFrames {
			// First field (frames) can have 3 digits if enabled
			maxDigits = 3
		}

		// If current field is full, auto-advance to next field first
		if len(e.fields[0]) >= maxDigits && len(e.fields) < 4 {
			// Insert new empty field at beginning, pushing current values right
			e.fields = append([]string{""}, e.fields...)
		}

		// Add digit to the first field (current entry field)
		e.fields[0] += string(r)
		e.updateDisplay()
		return
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
	case fyne.KeyPeriod, fyne.KeyComma:
		// Handle period/comma as field separator (for mobile keyboards)
		// Desktop keyboards send these as TypedRune, but mobile sends as key events
		if len(e.fields) < 4 && e.fields[0] != "" {
			// Insert empty field at beginning, pushing current values right
			e.fields = append([]string{""}, e.fields...)
			e.updateDisplay()
		}
		return
	case fyne.KeyBackspace, fyne.KeyDelete:
		// Remove last character from first field (current entry field)
		if len(e.fields) > 0 {
			if len(e.fields[0]) > 0 {
				// Remove character from current field
				e.fields[0] = e.fields[0][:len(e.fields[0])-1]
			} else if len(e.fields) > 1 {
				// If current field is empty, remove it and merge back to previous field
				e.fields = e.fields[1:]
			}
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

// updateDisplay converts fields to formatted timecode display
func (e *TimecodeEntry) updateDisplay() {
	// Fields are stored right-to-left: [frames, seconds, minutes, hours]
	// Display as HH:MM:SS:FF

	if len(e.fields) == 0 || (len(e.fields) == 1 && e.fields[0] == "") {
		e.Entry.SetText("")
		e.CursorColumn = 0
		return
	}

	// Map fields: [0]=frames, [1]=seconds, [2]=minutes, [3]=hours
	frames := "00"
	seconds := "00"
	minutes := "00"
	hours := "00"

	if len(e.fields) >= 1 && e.fields[0] != "" {
		// fields[0] = frames (current entry)
		if e.ThreeDigitFrames {
			frames = fmt.Sprintf("%03s", e.fields[0])
			if len(frames) > 3 {
				frames = frames[len(frames)-3:]
			}
		} else {
			frames = fmt.Sprintf("%02s", e.fields[0])
			if len(frames) > 2 {
				frames = frames[len(frames)-2:]
			}
		}
	}
	if len(e.fields) >= 2 && e.fields[1] != "" {
		// fields[1] = seconds
		seconds = fmt.Sprintf("%02s", e.fields[1])
		if len(seconds) > 2 {
			seconds = seconds[len(seconds)-2:]
		}
	}
	if len(e.fields) >= 3 && e.fields[2] != "" {
		// fields[2] = minutes
		minutes = fmt.Sprintf("%02s", e.fields[2])
		if len(minutes) > 2 {
			minutes = minutes[len(minutes)-2:]
		}
	}
	if len(e.fields) >= 4 && e.fields[3] != "" {
		// fields[3] = hours
		hours = fmt.Sprintf("%02s", e.fields[3])
		if len(hours) > 2 {
			hours = hours[len(hours)-2:]
		}
	}

	formatted := fmt.Sprintf("%s:%s:%s:%s", hours, minutes, seconds, frames)
	e.Entry.SetText(formatted)
	e.CursorColumn = len(formatted)
}

// GetComponents extracts hours, minutes, seconds, and frames from fields
// Returns (hours, minutes, seconds, frames)
func (e *TimecodeEntry) GetComponents() (int, int, int, int) {
	var h, m, s, f int

	// Map fields: [0]=frames, [1]=seconds, [2]=minutes, [3]=hours
	if len(e.fields) >= 1 && e.fields[0] != "" {
		f, _ = strconv.Atoi(e.fields[0])
	}
	if len(e.fields) >= 2 && e.fields[1] != "" {
		s, _ = strconv.Atoi(e.fields[1])
	}
	if len(e.fields) >= 3 && e.fields[2] != "" {
		m, _ = strconv.Atoi(e.fields[2])
	}
	if len(e.fields) >= 4 && e.fields[3] != "" {
		h, _ = strconv.Atoi(e.fields[3])
	}

	return h, m, s, f
}

// SetComponents sets the timecode from individual components
func (e *TimecodeEntry) SetComponents(hours, minutes, seconds, frames int) {
	// Map fields: [0]=frames, [1]=seconds, [2]=minutes, [3]=hours
	e.fields = make([]string, 4)

	if e.ThreeDigitFrames {
		e.fields[0] = fmt.Sprintf("%03d", frames)
	} else {
		e.fields[0] = fmt.Sprintf("%02d", frames)
	}
	e.fields[1] = fmt.Sprintf("%02d", seconds)
	e.fields[2] = fmt.Sprintf("%02d", minutes)
	e.fields[3] = fmt.Sprintf("%02d", hours)

	e.updateDisplay()
}

// SetText overrides to clear fields
func (e *TimecodeEntry) SetText(text string) {
	if text == "" {
		e.fields = []string{""}
	}
	e.Entry.SetText(text)
}
