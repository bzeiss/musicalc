package logic

import (
	"encoding/csv"
	"errors"
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

func trimFloat(v float64, decimals int) string {
	s := strconv.FormatFloat(v, 'f', decimals, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

func AlignmentDelayExportCSV(calc *AlignmentDelayCalculator, distUnit string, roomTemp float64, roomTempUnit string, refDist float64, refDistUnit string) (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)

	if err := w.Write([]string{
		"name",
		"ref to dist",
		"ref to dist measure",
		"delay ms",
		"delay samps",
		"samplerate",
		"room temp",
		"room temp measure",
		"ref dist",
		"ref dist measure",
	}); err != nil {
		return "", err
	}

	sampleRate := strconv.Itoa(calc.SampleRate)
	roomTempStr := trimFloat(roomTemp, 2)
	refDistStr := trimFloat(refDist, 3)

	for _, mic := range calc.Mics {
		dist := trimFloat(FromMeters(mic.DistanceMeters, distUnit), 3)
		delayMS := ""
		delaySamps := ""
		if !mic.IsBeyondReference {
			delayMS = trimFloat(mic.DelayMS, 2)
			delaySamps = strconv.Itoa(mic.DelaySamples)
		}
		if err := w.Write([]string{
			mic.Name,
			dist,
			distUnit,
			delayMS,
			delaySamps,
			sampleRate,
			roomTempStr,
			roomTempUnit,
			refDistStr,
			refDistUnit,
		}); err != nil {
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return sb.String(), nil
}

type AlignmentDelayImportRow struct {
	Name     string
	Dist     float64
	DistUnit string
}

type AlignmentDelayImportResult struct {
	SampleRate   int
	RoomTemp     float64
	RoomTempUnit string
	RefDist      float64
	RefDistUnit  string
	DistUnit     string
	Mics         []AlignmentDelayImportRow
}

func AlignmentDelayImportCSV(csvData string) (AlignmentDelayImportResult, error) {
	r := csv.NewReader(strings.NewReader(csvData))
	r.TrimLeadingSpace = true

	rows, err := r.ReadAll()
	if err != nil {
		return AlignmentDelayImportResult{}, err
	}
	if len(rows) == 0 {
		return AlignmentDelayImportResult{}, errors.New("empty csv")
	}

	start := 0
	if len(rows[0]) > 0 {
		if strings.EqualFold(strings.TrimSpace(rows[0][0]), "name") {
			start = 1
		}
	}

	var res AlignmentDelayImportResult
	for i := start; i < len(rows); i++ {
		rec := rows[i]
		if len(rec) == 0 {
			continue
		}
		if len(rec) < 10 {
			return AlignmentDelayImportResult{}, errors.New("invalid csv row: expected 10 columns")
		}

		name := strings.TrimSpace(rec[0])
		if name == "" {
			continue
		}

		dist := ParseFloat(rec[1])
		distUnit := strings.TrimSpace(rec[2])
		if distUnit == "" {
			distUnit = "m"
		}

		sampleRate := ParseSampleRate(rec[5])
		roomTemp := ParseFloat(rec[6])
		roomTempUnit := strings.TrimSpace(rec[7])
		if roomTempUnit == "" {
			roomTempUnit = "C"
		}
		refDist := ParseFloat(rec[8])
		refDistUnit := strings.TrimSpace(rec[9])
		if refDistUnit == "" {
			refDistUnit = "m"
		}

		if len(res.Mics) == 0 {
			res.SampleRate = sampleRate
			res.RoomTemp = roomTemp
			res.RoomTempUnit = roomTempUnit
			res.RefDist = refDist
			res.RefDistUnit = refDistUnit
			res.DistUnit = distUnit
		}

		res.Mics = append(res.Mics, AlignmentDelayImportRow{Name: name, Dist: dist, DistUnit: distUnit})
	}

	if len(res.Mics) == 0 {
		return AlignmentDelayImportResult{}, errors.New("no mic rows found")
	}
	if res.SampleRate <= 0 {
		res.SampleRate = 48000
	}

	return res, nil
}
