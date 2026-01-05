package logic

import (
	"math"
	"strconv"
	"strings"
)

type AlignmentMic struct {
	Name              string
	DistanceMeters    float64
	DelayMS           float64
	DelaySamples      int
	IsBeyondReference bool
}

type AlignmentDelayCalculator struct {
	TemperatureC            float64
	SampleRate              int
	ReferenceDistanceMeters float64
	Mics                    []AlignmentMic
}

func NewAlignmentDelayCalculator() *AlignmentDelayCalculator {
	return &AlignmentDelayCalculator{SampleRate: 48000}
}

func (c *AlignmentDelayCalculator) SetTemperature(value float64, unit string) {
	if strings.EqualFold(unit, "f") {
		c.TemperatureC = (value - 32) * 5 / 9
		return
	}
	c.TemperatureC = value
}

func (c *AlignmentDelayCalculator) SetSampleRateLabel(label string) {
	c.SampleRate = ParseSampleRate(label)
	if c.SampleRate <= 0 {
		c.SampleRate = 48000
	}
}

func (c *AlignmentDelayCalculator) SetReferenceDistance(value float64, unit string) {
	c.ReferenceDistanceMeters = ToMeters(value, unit)
}

func (c *AlignmentDelayCalculator) AddMic(name string, distance float64, unit string) {
	c.Mics = append(c.Mics, AlignmentMic{Name: name, DistanceMeters: ToMeters(distance, unit)})
}

func (c *AlignmentDelayCalculator) RemoveMicAt(index int) {
	if index < 0 || index >= len(c.Mics) {
		return
	}
	c.Mics = append(c.Mics[:index], c.Mics[index+1:]...)
}

func (c *AlignmentDelayCalculator) SpeedOfSoundMetersPerSecond() float64 {
	return 331.3 + (0.606 * c.TemperatureC)
}

func (c *AlignmentDelayCalculator) Recalculate() {
	for i := range c.Mics {
		delayMS, delaySamples, beyond := c.DelayForDistanceMeters(c.Mics[i].DistanceMeters)
		c.Mics[i].DelayMS = delayMS
		c.Mics[i].DelaySamples = delaySamples
		c.Mics[i].IsBeyondReference = beyond
	}
}

func (c *AlignmentDelayCalculator) DelayForDistanceMeters(targetMeters float64) (float64, int, bool) {
	distanceDiff := c.ReferenceDistanceMeters - targetMeters
	if distanceDiff < 0 {
		return 0, 0, true
	}

	speed := c.SpeedOfSoundMetersPerSecond()
	if speed <= 0 {
		return 0, 0, false
	}

	delayMS := (distanceDiff / speed) * 1000
	delaySamples := int(math.Round(delayMS * float64(c.SampleRate) / 1000))

	return delayMS, delaySamples, false
}

func ToMeters(distance float64, unit string) float64 {
	if strings.EqualFold(unit, "ft") {
		return distance * 0.3048
	}
	return distance
}

func FromMeters(distanceMeters float64, unit string) float64 {
	if strings.EqualFold(unit, "ft") {
		return distanceMeters / 0.3048
	}
	return distanceMeters
}

func ParseSampleRate(rateStr string) int {
	rateStr = strings.TrimSpace(rateStr)
	if rateStr == "" {
		return 48000
	}

	s := strings.ToLower(strings.ReplaceAll(rateStr, " ", ""))

	if strings.HasSuffix(s, "khz") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "khz"), 64)
		if err == nil {
			return int(math.Round(v * 1000))
		}
	}
	if strings.HasSuffix(s, "hz") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "hz"), 64)
		if err == nil {
			return int(math.Round(v))
		}
	}

	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return int(math.Round(v))
	}

	return 48000
}
