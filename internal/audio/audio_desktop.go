//go:build !android && !ios

package audio

import (
	"math"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

var (
	speakerInitialized bool
	speakerMutex       sync.Mutex
)

// InitAudio initializes the audio system for desktop platforms
func InitAudio() error {
	speakerMutex.Lock()
	defer speakerMutex.Unlock()

	if speakerInitialized {
		return nil
	}

	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return err
	}

	speakerInitialized = true
	return nil
}

// PlayFrequency plays a sine wave at the specified frequency for the given duration
func PlayFrequency(frequency float64, duration time.Duration) error {
	// Initialize speaker if needed
	if err := InitAudio(); err != nil {
		return err
	}

	// Stop any currently playing sounds to prevent artifacts
	speaker.Clear()

	sr := beep.SampleRate(44100)

	// Create sine wave generator
	sine := newSineWave(sr, frequency)

	// Take only the duration we want
	limited := beep.Take(sr.N(duration), sine)

	// Play the sound
	speaker.Play(limited)

	return nil
}

// sineWaveStreamer generates a sine wave at the specified frequency
type sineWaveStreamer struct {
	frequency  float64
	sampleRate beep.SampleRate
	position   float64
}

func newSineWave(sampleRate beep.SampleRate, frequency float64) *sineWaveStreamer {
	return &sineWaveStreamer{
		frequency:  frequency,
		sampleRate: sampleRate,
		position:   0,
	}
}

func (s *sineWaveStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		// Generate sine wave sample
		sample := math.Sin(2 * math.Pi * s.frequency * s.position / float64(s.sampleRate))
		samples[i][0] = sample // Left channel
		samples[i][1] = sample // Right channel

		s.position++
	}
	return len(samples), true
}

func (s *sineWaveStreamer) Err() error {
	return nil
}
