package logic

import (
	"math"
	"testing"
)

// TuningTestCase represents a test case with reference values from MusicMath
type TuningTestCase struct {
	TestName       string
	TuningName     string
	RefNoteName    string
	RefMidi        int
	RefFreq        float64
	OctaveOffset   int // 2 for C3 convention, 1 for C4 convention
	ExpectedValues []ExpectedNote
}

// ExpectedNote represents expected frequency and cents for a specific note
type ExpectedNote struct {
	NoteName      string
	Midi          int
	Frequency     float64
	Cents         float64
	FreqTolerance float64 // Allowed frequency deviation in Hz
	CentTolerance float64 // Allowed cents deviation
}

// TestTuningAccuracy tests tuning calculations against MusicMath reference values
func TestTuningAccuracy(t *testing.T) {
	testCases := []TuningTestCase{
		{
			TestName:     "Kepler Monochord No.1 with G2=500Hz",
			TuningName:   "Kepler's Monochord no.1, Harmonices Mundi (1619)",
			RefNoteName:  "G2",
			RefMidi:      55, // G2 in C3 convention: (2+2)*12 + 7 = 55
			RefFreq:      500.0,
			OctaveOffset: 2,
			ExpectedValues: []ExpectedNote{
				{NoteName: "C-2", Midi: 0, Frequency: 20.83, Cents: -2.00, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "C#-2", Midi: 1, Frequency: 21.97, Cents: -9.80, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "D-2", Midi: 2, Frequency: 23.44, Cents: +2.00, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "D#-2", Midi: 3, Frequency: 24.72, Cents: -5.90, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "E-2", Midi: 4, Frequency: 26.37, Cents: +5.90, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "F-2", Midi: 5, Frequency: 28.13, Cents: +17.60, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "F#-2", Midi: 6, Frequency: 29.30, Cents: -11.70, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "G-2", Midi: 7, Frequency: 31.25, Cents: 0.00, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "G#-2", Midi: 8, Frequency: 32.96, Cents: -7.80, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "A-2", Midi: 9, Frequency: 35.16, Cents: +3.90, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "A#-2", Midi: 10, Frequency: 37.50, Cents: +15.60, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "B-2", Midi: 11, Frequency: 39.06, Cents: -13.70, FreqTolerance: 0.01, CentTolerance: 0.5},
				// Test G2 reference itself
				{NoteName: "G2", Midi: 55, Frequency: 500.00, Cents: 0.00, FreqTolerance: 0.01, CentTolerance: 0.5},
			},
		},
		{
			TestName:     "12-tone Equal Temperament with A3=440Hz",
			TuningName:   "12 tone equal temperament",
			RefNoteName:  "A3",
			RefMidi:      69,
			RefFreq:      440.0,
			OctaveOffset: 2,
			ExpectedValues: []ExpectedNote{
				{NoteName: "A3", Midi: 69, Frequency: 440.00, Cents: 0.00, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "C3", Midi: 60, Frequency: 261.63, Cents: 0.00, FreqTolerance: 0.01, CentTolerance: 0.5},
				{NoteName: "C4", Midi: 72, Frequency: 523.25, Cents: 0.00, FreqTolerance: 0.01, CentTolerance: 0.5},
			},
		},
		{
			TestName:     "12-tone Pythagorean with A3=440Hz",
			TuningName:   "12-tone Pythagorean scale",
			RefNoteName:  "A3",
			RefMidi:      69,
			RefFreq:      440.0,
			OctaveOffset: 2,
			ExpectedValues: []ExpectedNote{
				{NoteName: "A3", Midi: 69, Frequency: 440.00, Cents: 0.00, FreqTolerance: 0.1, CentTolerance: 0.5},
				{NoteName: "C3", Midi: 60, Frequency: 260.74, Cents: -5.90, FreqTolerance: 0.1, CentTolerance: 0.5},
				{NoteName: "G3", Midi: 67, Frequency: 391.11, Cents: -3.90, FreqTolerance: 0.1, CentTolerance: 0.5},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			for _, expected := range tc.ExpectedValues {
				result := GetFrequency(expected.Midi, tc.RefFreq, tc.RefMidi, tc.TuningName)

				// Check frequency
				freqDiff := math.Abs(result.Frequency - expected.Frequency)
				if freqDiff > expected.FreqTolerance {
					t.Errorf("%s (MIDI %d): Frequency mismatch\n  Expected: %.2f Hz\n  Got:      %.2f Hz\n  Diff:     %.2f Hz (tolerance: %.2f Hz)",
						expected.NoteName, expected.Midi,
						expected.Frequency, result.Frequency,
						freqDiff, expected.FreqTolerance)
				}

				// Check cents
				centDiff := math.Abs(result.Cents - expected.Cents)
				if centDiff > expected.CentTolerance {
					t.Errorf("%s (MIDI %d): Cents mismatch\n  Expected: %+.2f cents\n  Got:      %+.2f cents\n  Diff:     %.2f cents (tolerance: %.2f cents)",
						expected.NoteName, expected.Midi,
						expected.Cents, result.Cents,
						centDiff, expected.CentTolerance)
				}

				// Log success for debugging
				if freqDiff <= expected.FreqTolerance && centDiff <= expected.CentTolerance {
					t.Logf("✓ %s (MIDI %d): %.2f Hz (%+.2f cents) - PASS",
						expected.NoteName, expected.Midi,
						result.Frequency, result.Cents)
				}
			}
		})
	}
}

// TestReferenceNoteAccuracy specifically tests that reference notes are tuned correctly
func TestReferenceNoteAccuracy(t *testing.T) {
	testCases := []struct {
		name          string
		tuningName    string
		refMidi       int
		refFreq       float64
		freqTolerance float64
	}{
		{"A4=440Hz in 12-TET", "12 tone equal temperament", 69, 440.0, 0.01},
		{"G2=500Hz in Kepler", "Kepler's Monochord no.1, Harmonices Mundi (1619)", 55, 500.0, 0.01},
		{"C4=261.63Hz in Pythagorean", "12-tone Pythagorean scale", 60, 261.63, 0.01},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetFrequency(tc.refMidi, tc.refFreq, tc.refMidi, tc.tuningName)

			freqDiff := math.Abs(result.Frequency - tc.refFreq)
			if freqDiff > tc.freqTolerance {
				t.Errorf("Reference note not tuned correctly\n  Expected: %.2f Hz\n  Got:      %.2f Hz\n  Diff:     %.4f Hz",
					tc.refFreq, result.Frequency, freqDiff)
			} else {
				t.Logf("✓ Reference note correctly tuned: %.2f Hz", result.Frequency)
			}
		})
	}
}
