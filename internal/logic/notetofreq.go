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
	cachedRefFreq    float64
	cachedRefMidi    int
)

// GetFrequency calculates the frequency for a given MIDI note using the specified tuning.
// Parameters:
//   - midiNote: MIDI note number (0-127)
//   - refFreq: Reference frequency in Hz (default 440 Hz if <= 0)
//   - refMidi: Reference MIDI note for refFreq (default 69 = A4 if <= 0)
//   - tuningName (optional): Name of tuning from scl.AvailableScales. Uses default if empty.
//
// Returns NoteFrequency with frequency in Hz and cents deviation from 12-TET.
func GetFrequency(midiNote int, refFreq float64, refMidi int, tuningName ...string) NoteFrequency {
	if refFreq <= 0 {
		refFreq = 440.0
	}

	if refMidi <= 0 || refMidi > 127 {
		refMidi = 69 // Default to A4
	}

	// Determine which tuning to use
	var selectedTuning string
	if len(tuningName) > 0 && tuningName[0] != "" {
		selectedTuning = tuningName[0]
	} else {
		selectedTuning = sclres.DefaultScaleName
	}

	// Load tuning (from cache or fresh)
	tuning := loadTuning(selectedTuning, refFreq, refMidi)
	if tuning == nil {
		// Fallback to 12-TET if tuning load fails
		freq := refFreq * math.Pow(2, float64(midiNote-refMidi)/12.0)
		return NoteFrequency{Frequency: freq, Cents: 0.0}
	}

	// Use the tuning to get the frequency directly
	// The tuning is calibrated to our reference frequency at refMidi
	freq := tuning.FrequencyForMidiNote(midiNote)

	// Calculate 12-TET frequency for comparison
	tetFreq := refFreq * math.Pow(2, float64(midiNote-refMidi)/12.0)

	// Calculate cents deviation from 12-TET
	// cents = 1200 * log2(freq / tetFreq)
	var cents float64
	if tetFreq > 0 && freq > 0 {
		cents = 1200.0 * math.Log2(freq/tetFreq)
	}

	return NoteFrequency{Frequency: freq, Cents: cents}
}

// loadTuning loads and caches the specified tuning by name, reference frequency, and reference MIDI note
func loadTuning(tuningName string, refFreq float64, refMidi int) scala.Tuning {
	// Check cache first
	tuningCacheMutex.RLock()
	if cachedTuningName == tuningName && cachedTuning != nil && cachedRefFreq == refFreq && cachedRefMidi == refMidi {
		tuning := cachedTuning
		tuningCacheMutex.RUnlock()
		return tuning
	}
	tuningCacheMutex.RUnlock()

	// Load tuning from resources
	tuningCacheMutex.Lock()
	defer tuningCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if cachedTuningName == tuningName && cachedTuning != nil && cachedRefFreq == refFreq && cachedRefMidi == refMidi {
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

	// Use StartScaleOnAndTuneNoteTo with scale starting on the reference note itself
	// This matches how MusicMath interprets the scale mapping
	kbm, err := scala.KeyboardMappingStartScaleOnAndTuneNoteTo(refMidi, refMidi, refFreq)
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
	cachedRefFreq = refFreq
	cachedRefMidi = refMidi

	return tuning
}
