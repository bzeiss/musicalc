//go:build android || ios

package audio

import (
	"math"
	"time"

	"golang.org/x/mobile/exp/audio/al"
)

var (
	alInitialized bool
	source        al.Source
	buffer        al.Buffer
)

// InitAudio initializes OpenAL for mobile platforms
// Uses golang.org/x/mobile/exp/audio/al which is compatible with Fyne's x/mobile
func InitAudio() error {
	if alInitialized {
		return nil
	}

	if err := al.OpenDevice(); err != nil {
		return err
	}

	source = al.GenSources(1)[0]
	alInitialized = true
	return nil
}

// PlayFrequency plays a sine wave at the specified frequency using OpenAL
// This uses the same golang.org/x/mobile that Fyne uses, avoiding conflicts
func PlayFrequency(frequency float64, duration time.Duration) error {
	if err := InitAudio(); err != nil {
		return err
	}

	// Generate PCM samples
	sampleRate := 44100
	numSamples := int(float64(sampleRate) * duration.Seconds())
	samples := generateSineWavePCM(frequency, sampleRate, numSamples)

	// Delete old buffer if exists
	if buffer != 0 {
		al.DeleteBuffers(buffer)
	}

	// Create new buffer and fill with PCM data
	buffer = al.GenBuffers(1)[0]
	al.BufferData(buffer, al.FormatMono16, samples, int32(sampleRate))

	// Stop any previous playback
	al.SourceStop(source)

	// Attach buffer to source and play
	al.Sourcei(source, al.Buffer, int32(buffer))
	al.SourcePlay(source)

	return nil
}

// generateSineWavePCM creates 16-bit mono PCM samples for a sine wave
func generateSineWavePCM(frequency float64, sampleRate, numSamples int) []byte {
	pcm := make([]byte, numSamples*2) // 2 bytes per 16-bit sample

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := math.Sin(2 * math.Pi * frequency * t)

		// Convert to 16-bit PCM
		value := int16(sample * 32767)

		// Little-endian encoding
		pcm[i*2] = byte(value & 0xff)
		pcm[i*2+1] = byte((value >> 8) & 0xff)
	}

	return pcm
}
