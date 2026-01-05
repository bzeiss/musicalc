# MusiCalc

A cross-platform music calculator application built with Go and Fyne. Provides real-time calculators for tempo, pitch, sampling, and time-stretching operations commonly used in music production and audio engineering.

This project is heavily inspired by [MusicMath](https://dev.laurentcolson.com/musicmath.html), an excellent tool currently exclusive to the Apple ecosystem. This version aims to bring that same utility to Windows, Linux and Android, with a simpler look and feel. For macOS users looking for a polished commercial solution, MusicMath is highly recommended. MusicMath has long been a reference tool for music-related calculations, and this project draws inspiration from that work.

<table style="width: 100%;" border="0">
  <tr>
    <td>
      <img width="450" height="682" alt="image" src="https://github.com/user-attachments/assets/8b6ecfaf-d865-4699-9926-d56f9806eb15" />
    </td>
    <td>
      <img width="450" height="679" alt="image" src="https://github.com/user-attachments/assets/0f95847c-aa18-463b-a47b-81f80bf3f957" />
    </td>
    <td>
      <img width="450" height="679" alt="image" src="https://github.com/user-attachments/assets/2db75878-294e-425c-9ef3-1086adba30a7" />
    </td>
  </tr>
</table>

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

- **Multi-Mic Time Alignment Matrix:** This feature calculates the exact millisecond delay needed to align multiple "close-mics" with a distant "reference-mic" based on their physical distance and the current studio temperature. By defining a zero-reference point (such as a room mic), engineers can instantly generate a list of delay offsets for an entire drum kit or orchestral section to ensure all transients line up perfectly in the DAW, resulting in a tighter and more focused sound.
- **Phase-Safe Distance & 3-to-1 Rule Helper:** This tool calculates physical "Sweet Spots" and "Death Zones" for microphone placement relative to the wavelength of a specific fundamental frequency, such as a kick drum's 60Hz thump. By mapping these phase relationships to physical distances, it helps engineers avoid destructive interference and includes a dedicated 3-to-1 rule calculator to ensure that bleed between multiple microphones remains phase-coherent and musically pleasing.
- **Acoustic Room Mode Calculator:** This utility predicts the specific resonant frequencies (standing waves) of a rectangular recording or mixing space by analyzing its length, width, and height. It identifies axial, tangential, and oblique modes to help engineers anticipate "bass build-up" or "frequency nulls," making it an indispensable tool for placing acoustic treatment and finding the most accurate listening position within a room.
- **Air Absorption & Humidity Compensator:** Designed for large-scale recording sessions and live sound reinforcement, this calculator determines how much high-frequency energy is naturally lost as sound travels through the air based on temperature and relative humidity. It provides the exact decibel boost required at specific frequencies (like 10kHz) to recover the "brilliance" lost over long distances, ensuring that distant microphones maintain the same clarity as close-up sources.
- **Tap tempo:** Detect BPM from a series of taps (maybe, is built in practically everywhere)

## Documentation

- **[User Guide](README-USERGUIDE.md)** - How to use each calculator
- **[Build Instructions](README-BUILD.md)** - How to build from source on Windows, Linux, and macOS
- **[Creating Installers](README-INSTALLER.md)** - Creating Windows installers and Linux packages
- **[Icon Resources](README-ICONS.md)** - Icon specifications and resources
- **[Test Instructions](README-TESTS.md)** - How to run tests

## Requirements

- Go 1.24.5 or later
- GCC/MinGW (for CGO support on Windows)
- Fyne v2.7.1

## License

MIT License - Copyright (c) 2026 B. Zeiss

See [LICENSE](LICENSE) for details.
