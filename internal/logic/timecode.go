package logic

import (
	"fmt"
	"math"
)

type TimecodeResult struct {
	Hours       int
	Minutes     int
	Seconds     int
	Frames      int
	TotalFrames int
	Timecode    string
}

type FPSFormat struct {
	Name      string
	FPS       float64
	DropFrame bool
}

var FPSFormats = []FPSFormat{
	{"23.976 fps", 24000.0 / 1001.0, false}, // Exact NTSC film rate
	{"24 fps", 24.0, false},
	{"25 fps", 25.0, false},
	{"29.94 fps", 30000.0 / 1001.0, false}, // Same as 29.97, exact NTSC video rate
	{"29.97 fps", 30000.0 / 1001.0, false}, // Exact NTSC video rate
	{"29.97 fps (df)", 30000.0 / 1001.0, true},
	{"30 fps", 30.0, false},
	{"50 fps", 50.0, false},
	{"59.94 fps", 60000.0 / 1001.0, false}, // Exact NTSC 60i rate
	{"60 fps", 60.0, false},
}

func GetFPSFormat(name string) FPSFormat {
	for _, format := range FPSFormats {
		if format.Name == name {
			return format
		}
	}
	return FPSFormats[6] // Default to 30 fps
}

// TimecodeToFrames converts timecode components to total frames
func TimecodeToFrames(hours, minutes, seconds, frames int, format FPSFormat) int {
	if format.DropFrame {
		// Drop frame calculation for 29.97 fps
		// Drop 2 frames every minute except every 10th minute
		totalMinutes := hours*60 + minutes
		dropFrames := 2 * (totalMinutes - totalMinutes/10)

		frameRate := int(math.Round(format.FPS))
		totalFrames := ((hours*3600 + minutes*60 + seconds) * frameRate) + frames - dropFrames
		return totalFrames
	} else {
		// Non-drop frame calculation
		frameRate := int(math.Round(format.FPS))
		totalFrames := ((hours*3600 + minutes*60 + seconds) * frameRate) + frames
		return totalFrames
	}
}

// FramesToTimecode converts total frames to timecode components
func FramesToTimecode(totalFrames int, format FPSFormat) TimecodeResult {
	result := TimecodeResult{
		TotalFrames: totalFrames,
	}

	if totalFrames < 0 {
		totalFrames = 0
	}

	frameRate := int(math.Round(format.FPS))

	if format.DropFrame {
		// Drop frame calculation for 29.97 fps
		// Add back dropped frames for calculation
		d := totalFrames / 17982 // Number of 10-minute intervals
		m := totalFrames % 17982

		if m < 2 {
			m = 2
		}

		adjustedFrames := totalFrames + 18*d + 2*((m-2)/1798)

		result.Hours = adjustedFrames / (frameRate * 3600)
		adjustedFrames %= (frameRate * 3600)
		result.Minutes = adjustedFrames / (frameRate * 60)
		adjustedFrames %= (frameRate * 60)
		result.Seconds = adjustedFrames / frameRate
		result.Frames = adjustedFrames % frameRate
	} else {
		// Non-drop frame calculation
		result.Hours = totalFrames / (frameRate * 3600)
		totalFrames %= (frameRate * 3600)
		result.Minutes = totalFrames / (frameRate * 60)
		totalFrames %= (frameRate * 60)
		result.Seconds = totalFrames / frameRate
		result.Frames = totalFrames % frameRate
	}

	result.Timecode = fmt.Sprintf("%02d:%02d:%02d:%02d", result.Hours, result.Minutes, result.Seconds, result.Frames)

	return result
}

// AddTimecodes adds two timecodes
func AddTimecodes(h1, m1, s1, f1, h2, m2, s2, f2 int, format FPSFormat) TimecodeResult {
	frames1 := TimecodeToFrames(h1, m1, s1, f1, format)
	frames2 := TimecodeToFrames(h2, m2, s2, f2, format)
	return FramesToTimecode(frames1+frames2, format)
}

// SubtractTimecodes subtracts two timecodes
func SubtractTimecodes(h1, m1, s1, f1, h2, m2, s2, f2 int, format FPSFormat) TimecodeResult {
	frames1 := TimecodeToFrames(h1, m1, s1, f1, format)
	frames2 := TimecodeToFrames(h2, m2, s2, f2, format)
	totalFrames := frames1 - frames2
	if totalFrames < 0 {
		totalFrames = 0
	}
	return FramesToTimecode(totalFrames, format)
}
