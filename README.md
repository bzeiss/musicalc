# MusiCalc

A cross-platform music calculator application built with Go and Fyne. Provides real-time calculators for tempo, pitch, sampling, and time-stretching operations commonly used in music production and audio engineering.

This project is heavily inspired by [MusicMath](https://dev.laurentcolson.com/musicmath.html), an excellent tool currently exclusive to the Apple ecosystem. This version aims to bring that same utility to Windows and Linux, with a simpler look and feel. For macOS users looking for a polished commercial solution, MusicMath is highly recommended.

<img width="664" height="781" alt="image" src="https://github.com/user-attachments/assets/26dd8ca2-a173-4f0e-91e0-1a62b615870f" />

## Disclaimer

This application is provided for informational purposes only and is not guaranteed to be accurate. Calculations may contain errors due to software bugs, rounding, or floating-point limitations. Use this software entirely at your own risk; the author assumes no responsibility for any financial loss, data loss, or damages resulting from its use. Always manually verify critical results.

## Features

### ⏱️ Timecode Calculator
- Support for multiple frame rates: 23.976, 24, 25, 29.94, 29.97, 29.97 drop frame, 30, 50, 59.94, 60 fps
- Precise NTSC frame rate handling with exact fractional values
- Dual timecode input fields with real-time calculation
- Add and subtract timecode operations
- Frame count preservation when switching frame rates (H:M:S:F notation stays constant)
- Compact display format showing timecode, frame count, and FPS
- Unlimited calculation history with copy/paste support
- Essential for video editing, post-production, and audio-for-video work

### 🎵 Tempo to Delay Calculator
- Calculate delay times (ms) and modulation frequencies (Hz) for various note divisions
- Supports standard note values: whole notes to 1/64 notes
- Includes dotted and triplet variations
- Real-time BPM input with instant recalculation
- Perfect for setting up delay effects, LFOs, and rhythmic modulation

### 🎹 Note to Frequency Calculator
- Complete MIDI note range: C-2 to B8 (MIDI 0-131)
- Custom reference tuning support (default A3 = 440 Hz)
- Displays frequencies for all chromatic notes
- Dual MIDI convention display (C4=60 standard / C3=60 alternative)
- Real-time frequency calculation with adjustable reference pitch
- Ideal for tuning synthesizers, creating custom scales, and frequency analysis

### 🎼 Frequency to Note Calculator
- Convert any frequency to the closest musical note
- Dual notation display: 100-cent and 50-cent systems
- Shows pitch deviation in cents for fine-tuning accuracy
- Quick-select buttons for common reference frequencies
- Adjustable middle C convention (C3/C4)
- Custom reference tuning support (default A4 = 440 Hz)
- Perfect for analyzing recordings, tuning acoustic instruments, and spectrum analysis

### ⏱️ Sample Length Calculator
- Bidirectional calculation: change any field and others update automatically
- Calculate sample count from tempo, beats, and sample rate
- Calculate tempo from sample length and beats
- Calculate duration (ms) from any combination of parameters
- Support for multiple sample rates (8kHz to 192kHz)
- Adjustable beat divisions for loop calculations
- Essential for sampler programming and loop creation

### 🎚️ Tempo Change Calculator
- Calculate pitch changes from tempo adjustments
- Time stretching percentage calculations
- Transpose semitones and cents (standard 100-cent notation)
- 50-cent notation display for granular pitch control
- Tempo delta percentage calculation
- Dynamic clamping to valid tempo range (5-1000 BPM)
- Bidirectional: adjust tempo, time stretch, or transpose values
- Swap tempo values with one click
- Perfect for time-stretching audio, pitch correction, and sampler tuning

## TODOs

- Support for different temperaments in Note To Frequency Calculator
- Tap tempo (maybe)

## Documentation

- **[User Guide](README-USERGUIDE.md)** - How to use each calculator
- **[Build Instructions](README-BUILD.md)** - How to build from source on Windows, Linux, and macOS
- **[Creating Installers](README-INSTALLER.md)** - Creating Windows installers and Linux packages
- **[Icon Resources](README-ICONS.md)** - Icon specifications and resources

## Requirements

- Go 1.24.5 or later
- GCC/MinGW (for CGO support on Windows)
- Fyne v2.7.1

## License

MIT License - Copyright (c) 2026 B. Zeiss

See [LICENSE](LICENSE) for details.
