package logic

import (
	"math"
)

type FrequencyResult struct {
	Note100          string
	Cents100         int
	Note50           string
	Cents50          int
	NearestMIDI      int
	NearestFrequency float64
}

// GetC3Frequency returns the frequency of C3 (middle C in some conventions)
func GetC3Frequency() float64 {
	// C3 is MIDI note 48, using A4=440Hz reference (MIDI 69)
	// Formula: freq = 440 * 2^((midi - 69) / 12)
	return 440.0 * math.Pow(2.0, (48.0-69.0)/12.0)
}

// GetA3Frequency returns the frequency of A3
func GetA3Frequency() float64 {
	// A3 is MIDI note 57, using A4=440Hz reference (MIDI 69)
	return 440.0 * math.Pow(2.0, (57.0-69.0)/12.0)
}

// GetC4Frequency returns the frequency of C4 (middle C)
func GetC4Frequency() float64 {
	// C4 is MIDI note 60, using A4=440Hz reference (MIDI 69)
	return 440.0 * math.Pow(2.0, (60.0-69.0)/12.0)
}

// GetA4Frequency returns the frequency of A4 (standard reference)
func GetA4Frequency() float64 {
	// A4 is MIDI note 69, the standard reference frequency
	return 440.0
}

// FrequencyToNote converts a frequency to the nearest note and cent offset
func FrequencyToNote(frequency float64) FrequencyResult {
	result := FrequencyResult{}

	if frequency <= 0 {
		return result
	}

	// Calculate semitones from A4 (440 Hz, MIDI 69)
	semitonesFromA4 := 12.0 * math.Log2(frequency/440.0)

	// Calculate MIDI note number (69 = A4)
	midiFloat := 69.0 + semitonesFromA4
	nearestMIDI := int(math.Round(midiFloat))

	// Clamp MIDI to valid range
	if nearestMIDI < 0 {
		nearestMIDI = 0
	}
	if nearestMIDI > 127 {
		nearestMIDI = 127
	}

	result.NearestMIDI = nearestMIDI

	// Calculate nearest note frequency
	result.NearestFrequency = 440.0 * math.Pow(2.0, (float64(nearestMIDI)-69.0)/12.0)

	// Get note name
	noteNames := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	noteIndex := nearestMIDI % 12
	octave := (nearestMIDI / 12) - 2 // MIDI 0 = C-2

	result.Note100 = noteNames[noteIndex]
	result.Note50 = noteNames[noteIndex]

	// Calculate cent offset for 100-cent notation
	// Cents = 1200 * log2(frequency / nearestNoteFrequency)
	centOffset := 1200.0 * math.Log2(frequency/result.NearestFrequency)
	result.Cents100 = int(math.Round(centOffset))

	// For 50-cent notation, calculate total semitones including fractional part
	totalSemitones := semitonesFromA4 + 69.0 // Total semitones from C-2

	// Use floor for negative numbers to match the tempo change logic
	semitones50 := int(math.Floor(totalSemitones))
	cents50Float := (totalSemitones - float64(semitones50)) * 100.0
	result.Cents50 = int(math.Round(cents50Float))

	// Get note name for 50-cent notation (from the floored semitone)
	noteIndex50 := semitones50 % 12
	if noteIndex50 < 0 {
		noteIndex50 += 12
	}
	octave50 := (semitones50 / 12) - 2

	result.Note50 = noteNames[noteIndex50]

	// Format note strings with octave
	if octave >= -2 && octave <= 8 {
		result.Note100 = result.Note100 + formatOctave(octave)
	}
	if octave50 >= -2 && octave50 <= 8 {
		result.Note50 = result.Note50 + formatOctave(octave50)
	}

	return result
}

func formatOctave(octave int) string {
	if octave < 0 {
		return "-" + string(rune('0'-octave))
	}
	return string(rune('0' + octave))
}
