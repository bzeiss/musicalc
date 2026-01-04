package logic

import (
	"bytes"
	"math"
	"sync"

	sclres "musicalc/internal/logic/scl"

	scala "github.com/chinenual/go-scala"
)

// NoteFrequency holds the frequency and cents deviation for a note
type NoteFrequency struct {
	Frequency float64 // Frequency in Hz
	Cents     float64 // Cents deviation from 12-TET (positive = sharp, negative = flat)
}

// tuningCache holds the currently loaded tuning to avoid reloading on every call
var (
	tuningCacheMutex sync.RWMutex
	cachedTuningName string
	cachedTuning     scala.Tuning
	cachedRefA4      float64
)

// GetFrequency calculates the frequency for a given MIDI note using the specified tuning.
// Parameters:
//   - midiNote: MIDI note number (0-127)
//   - refA4: Reference frequency for A4 (default 440 Hz if <= 0)
//   - tuningName (optional): Name of tuning from scl.AvailableScales. Uses default if empty.
//
// Returns NoteFrequency with frequency in Hz and cents deviation from 12-TET.
func GetFrequency(midiNote int, refA4 float64, tuningName ...string) NoteFrequency {
	if refA4 <= 0 {
		refA4 = 440.0
	}

	// Determine which tuning to use
	var selectedTuning string
	if len(tuningName) > 0 && tuningName[0] != "" {
		selectedTuning = tuningName[0]
	} else {
		selectedTuning = sclres.DefaultScaleName
	}

	// Load tuning (from cache or fresh)
	tuning := loadTuning(selectedTuning, refA4)
	if tuning == nil {
		// Fallback to 12-TET if tuning load fails
		freq := refA4 * math.Pow(2, float64(midiNote-69)/12.0)
		return NoteFrequency{Frequency: freq, Cents: 0.0}
	}

	// Use the tuning to get the frequency directly
	// The tuning is calibrated to our refA4 frequency
	freq := tuning.FrequencyForMidiNote(midiNote)

	// Calculate 12-TET frequency for comparison
	tetFreq := refA4 * math.Pow(2, float64(midiNote-69)/12.0)

	// Calculate cents deviation from 12-TET
	// cents = 1200 * log2(freq / tetFreq)
	var cents float64
	if tetFreq > 0 && freq > 0 {
		cents = 1200.0 * math.Log2(freq/tetFreq)
	}

	return NoteFrequency{Frequency: freq, Cents: cents}
}

// loadTuning loads and caches the specified tuning by name and reference frequency
func loadTuning(tuningName string, refA4 float64) scala.Tuning {
	// Check cache first
	tuningCacheMutex.RLock()
	if cachedTuningName == tuningName && cachedTuning != nil && cachedRefA4 == refA4 {
		tuning := cachedTuning
		tuningCacheMutex.RUnlock()
		return tuning
	}
	tuningCacheMutex.RUnlock()

	// Load tuning from resources
	tuningCacheMutex.Lock()
	defer tuningCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if cachedTuningName == tuningName && cachedTuning != nil && cachedRefA4 == refA4 {
		return cachedTuning
	}

	// Get the scale info from available scales
	scaleInfo, exists := sclres.AvailableScales[tuningName]
	if !exists {
		return nil
	}

	// Parse the SCL file content using go-scala
	reader := bytes.NewReader(scaleInfo.Resource.Content())
	scale, err := scala.ScaleFromSCLStream(reader)
	if err != nil {
		return nil
	}

	// Create keyboard mapping tuned to our reference A4
	// Standard keyboard mapping with A69 (MIDI 69) tuned to refA4
	kbm, err := scala.KeyboardMappingTuneA69To(refA4)
	if err != nil {
		return nil
	}

	// Create tuning from scale and keyboard mapping
	tuning, err := scala.TuningFromSCLAndKBM(scale, kbm)
	if err != nil {
		return nil
	}

	// Cache the loaded tuning
	cachedTuningName = tuningName
	cachedTuning = tuning
	cachedRefA4 = refA4

	return tuning
}
