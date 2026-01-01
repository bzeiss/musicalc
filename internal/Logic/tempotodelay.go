package logic

type TempoResult struct {
	DelayMS float64
	ModHz   float64
}

func GetTempoData(bpm float64, multiplier float64) TempoResult {
	if bpm <= 0 { bpm = 120.0 }
	d := (60000.0 / bpm) * multiplier
	return TempoResult{DelayMS: d, ModHz: 1000.0 / d}
}