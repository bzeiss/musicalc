package widgets

import (
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

// NumericEntry is an Entry widget that shows numeric keyboard on mobile
type NumericEntry struct {
	widget.Entry
}

// NewNumericEntry creates a new numeric entry widget
func NewNumericEntry() *NumericEntry {
	e := &NumericEntry{}
	e.ExtendBaseWidget(e)
	return e
}

// Keyboard returns the numeric keyboard type for mobile devices
func (e *NumericEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
