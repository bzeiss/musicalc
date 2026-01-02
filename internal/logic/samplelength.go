package logic

type SamplerResult struct {
	Samples int
	MS      float64
}

func GetSampleData(rate, bpm, beats float64) SamplerResult {
	if bpm <= 0 || rate <= 0 {
		return SamplerResult{}
	}
	ms := (60.0 / bpm) * beats * 1000.0
	return SamplerResult{Samples: int((ms / 1000.0) * rate), MS: ms}
}
