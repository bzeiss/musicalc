package ui

import (
	"fmt"
	"math"
	"musicalc/internal/logic"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

var (
	speakerInitialized bool
	speakerMutex       sync.Mutex
)

func NewFrequencyToNoteTab() fyne.CanvasObject {
	// Input field
	frequencyEntry := widget.NewEntry()
	frequencyEntry.SetPlaceHolder("Enter frequency (Hz)")
	frequencyEntry.SetText("440.00")

	// Middle C convention selector
	middleCSelect := widget.NewSelect([]string{"C3", "C4"}, nil)
	middleCSelect.SetSelected("C3")

	// Output labels
	note100Label := widget.NewLabel("A4")
	cents100Label := widget.NewLabel("0")
	note50Label := widget.NewLabel("A4")
	cents50Label := widget.NewLabel("0")

	// Flag to prevent circular updates
	updating := false

	// Calculate and update all fields
	calculateFromFrequency := func() {
		if updating {
			return
		}
		updating = true
		defer func() { updating = false }()

		freq := logic.ParseFloat(frequencyEntry.Text)
		if freq > 0 {
			// Determine octave offset based on Middle C setting
			octaveOffset := 1 // C4 convention (default)
			if middleCSelect.Selected == "C3" {
				octaveOffset = 2 // C3 convention
			}

			result := logic.FrequencyToNote(freq, octaveOffset)

			note100Label.SetText(result.Note100)

			// Format cents with sign
			cents100Sign := ""
			if result.Cents100 > 0 {
				cents100Sign = "+"
			}
			cents100Label.SetText(fmt.Sprintf("%s%d", cents100Sign, result.Cents100))

			note50Label.SetText(result.Note50)

			cents50Sign := ""
			if result.Cents50 > 0 {
				cents50Sign = "+"
			}
			cents50Label.SetText(fmt.Sprintf("%s%d", cents50Sign, result.Cents50))
		}
	}

	// Wire up change handlers
	frequencyEntry.OnChanged = func(s string) {
		calculateFromFrequency()
	}

	middleCSelect.OnChanged = func(s string) {
		calculateFromFrequency()
	}

	// Quick-select buttons
	c3Button := widget.NewButton("C3", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCSelect.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("C", 3, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	a3Button := widget.NewButton("A3", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCSelect.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("A", 3, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	c4Button := widget.NewButton("C4", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCSelect.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("C", 4, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	a4Button := widget.NewButton("A4", func() {
		octaveOffset := 1 // C4 convention (default)
		if middleCSelect.Selected == "C3" {
			octaveOffset = 2 // C3 convention
		}
		freq := logic.GetFrequencyForNote("A", 4, octaveOffset)
		frequencyEntry.SetText(fmt.Sprintf("%.2f", freq))
	})

	// Play button - generates sine wave tone
	playButton := widget.NewButton("▶ Play", func() {
		freq := logic.ParseFloat(frequencyEntry.Text)
		if freq > 0 && freq < 20000 {
			go playSineWave(freq, 3*time.Second)
		}
	})

	// Initialize with default frequency (A4 = 440.00 Hz)
	calculateFromFrequency()

	return container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Frequency (Hz):"),
			frequencyEntry,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Middle C:"),
			middleCSelect,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Quick-Select:"),
			container.NewHBox(c3Button, a3Button, c4Button, a4Button),
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Note:"),
			note100Label,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Cents:"),
			cents100Label,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Play:"),
			playButton,
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewLabel("Note (50 cents):"),
			note50Label,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Cents (50 cents):"),
			cents50Label,
		),
	)
}

// initSpeaker initializes the audio speaker if not already done
func initSpeaker() error {
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

// playSineWave generates and plays a sine wave at the given frequency for the specified duration
func playSineWave(frequency float64, duration time.Duration) {
	// Initialize speaker if needed
	if err := initSpeaker(); err != nil {
		return
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
}
