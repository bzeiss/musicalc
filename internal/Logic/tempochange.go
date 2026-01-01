package logic

import "math"

// TempoChangeResult holds all calculated values for tempo change
type TempoChangeResult struct {
	NewTempo           float64
	TimeStretchPercent float64
	TempoVariation     float64
	Semitones          int
	Cents              int
	Semitones50Cent    int // For 50-cent notation
	Cents50Cent        int // For 50-cent notation
}

// CalculateFromNewTempo calculates all values when new tempo is known
func CalculateFromNewTempo(originalTempo, newTempo float64) TempoChangeResult {
	if originalTempo <= 0 || newTempo <= 0 {
		return TempoChangeResult{}
	}

	ratio := originalTempo / newTempo
	timeStretchPercent := ratio * 100.0
	tempoVariation := ((newTempo - originalTempo) / originalTempo) * 100.0

	// Calculate semitones: semitones = 12 * log2(newTempo/originalTempo)
	semitonesFull := 12.0 * math.Log2(newTempo/originalTempo)

	// For display: semitones is the floor/ceil, cents is the remainder
	// When negative: -32.04 semitones = -33 semitones + 96 cents
	var semitones, cents int
	if semitonesFull >= 0 {
		semitones = int(math.Floor(semitonesFull))
		cents = int(math.Round((semitonesFull - float64(semitones)) * 100.0))
	} else {
		semitones = int(math.Floor(semitonesFull))
		cents = int(math.Round((semitonesFull - float64(semitones)) * 100.0))
	}

	// 50-cent notation shows raw semitone split (floor + fractional)
	// E.g., -55.69 = floor(-55.69) = -56, fractional = -55.69 - (-56) = 0.31 → +31 cents
	semitones50 := int(math.Floor(semitonesFull))
	cents50 := int(math.Round((semitonesFull - float64(semitones50)) * 100.0))

	return TempoChangeResult{
		NewTempo:           newTempo,
		TimeStretchPercent: timeStretchPercent,
		TempoVariation:     tempoVariation,
		Semitones:          semitones,
		Cents:              cents,
		Semitones50Cent:    semitones50,
		Cents50Cent:        cents50,
	}
}

// CalculateFromTimeStretch calculates all values when time stretch % is known
func CalculateFromTimeStretch(originalTempo, timeStretchPercent float64) TempoChangeResult {
	if originalTempo <= 0 || timeStretchPercent <= 0 {
		return TempoChangeResult{}
	}

	newTempo := originalTempo / (timeStretchPercent / 100.0)
	return CalculateFromNewTempo(originalTempo, newTempo)
}

// CalculateFromTranspose calculates all values when transpose semitones/cents are known
func CalculateFromTranspose(originalTempo float64, semitones, cents int) TempoChangeResult {
	if originalTempo <= 0 {
		return TempoChangeResult{}
	}

	// Convert semitones and cents to full semitones
	totalSemitones := float64(semitones) + (float64(cents) / 100.0)

	// Calculate new tempo: newTempo = originalTempo * 2^(semitones/12)
	newTempo := originalTempo * math.Pow(2.0, totalSemitones/12.0)

	return CalculateFromNewTempo(originalTempo, newTempo)
}

// CalculateFromCentsChange calculates transpose values when cents change is known
func CalculateFromCentsChange(centsChange int) (semitones, cents int) {
	// Convert total cents to semitones and remaining cents
	totalSemitones := float64(centsChange) / 100.0
	semitones = int(totalSemitones)
	cents = centsChange - (semitones * 100)

	return semitones, cents
}
