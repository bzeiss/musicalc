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

// GetFrequencyForNote returns the frequency for a given note name and octave based on octave convention
// noteName: e.g., "C", "A", "C#"
// octave: e.g., 3, 4
// octaveOffset: 1 for C4 convention, 2 for C3 convention
func GetFrequencyForNote(noteName string, octave int, octaveOffset int) float64 {
	noteNames := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	// Find note index
	noteIndex := -1
	for i, n := range noteNames {
		if n == noteName {
			noteIndex = i
			break
		}
	}

	if noteIndex == -1 {
		return 0.0
	}

	// Calculate MIDI note number based on octave convention
	// MIDI = (octave + octaveOffset) * 12 + noteIndex
	midiNote := (octave+octaveOffset)*12 + noteIndex

	// Calculate frequency from MIDI note
	// freq = 440 * 2^((midi - 69) / 12)
	return 440.0 * math.Pow(2.0, (float64(midiNote)-69.0)/12.0)
}

// FrequencyToNote converts a frequency to the nearest note and cent offset
// octaveOffset: 1 for C4 convention (MIDI 60 = C4), 2 for C3 convention (MIDI 60 = C3)
func FrequencyToNote(frequency float64, octaveOffset int) FrequencyResult {
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
	octave := (nearestMIDI / 12) - octaveOffset

	// 100-cent notation: floor semitones, show 0-99 cents
	// Calculate total semitones including fractional part
	totalSemitones := semitonesFromA4 + 69.0 // Total semitones from C-2

	// Use floor for semitones
	semitones100 := int(math.Floor(totalSemitones))
	cents100Float := (totalSemitones - float64(semitones100)) * 100.0
	result.Cents100 = int(math.Round(cents100Float))

	// Get note name for 100-cent notation (from the floored semitone)
	noteIndex100 := semitones100 % 12
	if noteIndex100 < 0 {
		noteIndex100 += 12
	}
	octave100 := (semitones100 / 12) - octaveOffset

	result.Note100 = noteNames[noteIndex100]

	// 50-cent notation: round to nearest semitone, show ±50 cents
	// Cents = 1200 * log2(frequency / nearestNoteFrequency)
	centOffset := 1200.0 * math.Log2(frequency/result.NearestFrequency)
	result.Cents50 = int(math.Round(centOffset))

	result.Note50 = noteNames[noteIndex]

	// Format note strings with octave
	if octave100 >= -1 && octave100 <= 9 {
		result.Note100 = result.Note100 + formatOctave(octave100)
	}
	if octave >= -1 && octave <= 9 {
		result.Note50 = result.Note50 + formatOctave(octave)
	}

	return result
}

func formatOctave(octave int) string {
	if octave < 0 {
		return "-" + string(rune('0'-octave))
	}
	return string(rune('0' + octave))
}
