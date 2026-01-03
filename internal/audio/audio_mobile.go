//go:build android || ios

package audio

import (
	"time"
)

// InitAudio is a no-op stub for mobile platforms
// Audio playback is disabled on Android/iOS to avoid linker conflicts
func InitAudio() error {
	// No-op: audio is not supported on mobile to avoid Fyne/Oto conflicts
	return nil
}

// PlayFrequency is a no-op stub for mobile platforms
// Audio playback is disabled on Android/iOS to avoid linker conflicts
func PlayFrequency(frequency float64, duration time.Duration) error {
	// No-op: audio is not supported on mobile to avoid Fyne/Oto conflicts
	return nil
}
