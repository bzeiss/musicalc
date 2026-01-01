package logic

import "math"

func GetFrequency(midiNote int, refA4 float64) float64 {
	if refA4 <= 0 {
		refA4 = 440.0
	}
	// Standard MIDI frequency formula: f = 440 * 2^((m-69)/12)
	return refA4 * math.Pow(2, float64(midiNote-69)/12.0)
}
