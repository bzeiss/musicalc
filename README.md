# MusiCalc

A cross-platform music calculator application built with Go and Fyne. Provides real-time calculators for tempo, pitch, sampling, and time-stretching operations commonly used in music production and audio engineering.

This project is heavily inspired by [MusicMath](https://dev.laurentcolson.com/musicmath.html), an excellent tool currently exclusive to the Apple ecosystem. This version aims to bring that same utility to all platforms, reimagined with a unique look and feel. To support the original creator, I will not be providing pre-compiled macOS binaries; I encourage macOS users to purchase the original MusicMath instead. However, the source remains open for those who prefer to build it manually.

## Features

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

- Timecode Calculator
- Tap tempo (maybe)
- Frequency to note

## User Guide

### Tempo to Delay

1. **Set your project tempo** in the BPM field at the top
2. **View delay times** for each note division in the table
3. **Use the values** to configure:
   - Delay effect times for rhythmic echoes
   - LFO rates for synchronized modulation
   - Gate/sequencer timing

**Example**: At 120 BPM, a 1/4 note = 500 ms, perfect for a quarter-note delay

### Note to Frequency

1. **Set reference frequency** (default: 440 Hz for A3)
2. **Select reference note** if using non-standard tuning
3. **Browse the frequency table** to find:
   - Exact frequencies for synthesizer tuning
   - MIDI note numbers for programming
   - Pitch relationships between notes

**Example**: To tune a 808 kick to your track's key, find the root note frequency and adjust the kick's pitch to match

### Sample Length

1. **Choose your workflow**:
   - Enter **Tempo** to calculate sample length for a specific BPM
   - Enter **Length in samples** to find what tempo matches your loop
   - Enter **Length in ms** to calculate both tempo and sample count

2. **Set parameters**:
   - **Sample Rate**: Select your project's sample rate (typically 44100 or 48000)
   - **Beats**: Set loop length (4 = 4-beat loop, 8 = 8-beat loop, etc.)

3. **Read calculated values** instantly in all other fields

**Example Use Cases**:
- You have a 88200-sample loop at 44.1kHz → Calculator shows it's 2000ms and fits 120 BPM at 4 beats
- You want a 4-bar loop at 140 BPM → Calculator shows you need 151,543 samples at 44.1kHz
- You need to know what tempo your found sample is → Enter its sample count and get the BPM

### Tempo Change

1. **Set original tempo** of your audio material
2. **Choose your operation**:
   
   **Option A - Change Tempo Directly**:
   - Enter **New Tempo** to see pitch shift required
   - View **Time Stretch %** (e.g., 200% = double speed)
   - See **Tempo Delta** (percentage change)
   
   **Option B - Change Time Stretch**:
   - Enter **Time Stretching %** to see resulting tempo
   - View pitch changes in semitones and cents
   
   **Option C - Transpose**:
   - Enter **Transpose Semis** and **Cents** for desired pitch shift
   - View resulting tempo and time stretch amount

3. **Use the Swap button** to quickly reverse tempo/new tempo values
4. **Use Reset** to return to default values (Tempo: 140, New Tempo: 22)

**Understanding the Output**:
- **Transpose Semis/Cents**: Standard pitch notation (100 cents = 1 semitone)
- **50 Cents Notation**: Alternative pitch display (50 cents = 1 semitone) used by some samplers
- **Tempo Delta**: Percentage change from original tempo (+/- %)

**Example Use Cases**:
- Speed up a 140 BPM loop to 170 BPM → See it requires +3.45 semitones pitch shift
- You need to pitch a sample up 7 semitones → See it will play at 150% speed
- Match a sample's tempo to your project without changing pitch → Use time-stretching calculations

## Build Instructions

### Windows

#### Prerequisites

1. **Clone the repository**
   ```powershell
   git clone https://github.com/yourusername/musicalc.git
   cd musicalc
   ```

2. **Install MSYS2** (provides MinGW GCC compiler required for Fyne)
   ```powershell
   winget install -e --id MSYS2.MSYS2
   ```

3. **Install GCC via MSYS2**
   
   Open "MSYS2 MSYS" from Start Menu and run:
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```

4. **Add MinGW to PATH**
   
   Add `C:\msys64\mingw64\bin` to your system PATH environment variable
   
   - Press `Win+R` and type `sysdm.cpl`, press Enter
   - Go to "Advanced" tab → "Environment Variables"
   - Edit "Path" under System variables
   - Add new entry: `C:\msys64\mingw64\bin`
   - Click OK and restart your terminal

5. **Verify GCC installation**
   ```powershell
   gcc --version
   ```
   Should output GCC version information

6. **Install Go dependencies**
   ```powershell
   go mod download
   ```

7. **Build the application**
   ```powershell
   go build -ldflags="-s -w" -o musicalc.exe
   ```
or for a production build without the console:
   ```powershell
   go build -ldflags="-s -w -H=windowsgui" -o musicalc.exe
   ```

8. **Run the application**
   ```powershell
   .\musicalc.exe
   ```

### Linux

#### Debian/Ubuntu

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/musicalc.git
   cd musicalc
   ```

2. **Install dependencies**
   ```bash
   sudo apt-get update
   sudo apt-get install gcc libgl1-mesa-dev xorg-dev
   ```

3. **Install Go dependencies**
   ```bash
   go mod download
   ```

4. **Build the application**
   ```bash
   go build -ldflags="-s -w" -o musicalc
   ```

5. **Run the application**
   ```bash
   ./musicalc
   ```

#### Fedora

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/musicalc.git
   cd musicalc
   ```

2. **Install dependencies**
   ```bash
   sudo dnf install gcc mesa-libGL-devel libXcursor-devel libXrandr-devel \
                    libXinerama-devel libXi-devel libXxf86vm-devel
   ```

3. **Install Go dependencies**
   ```bash
   go mod download
   ```

4. **Build the application**
   ```bash
   go build -ldflags="-s -w" -o musicalc
   ```

5. **Run the application**
   ```bash
   ./musicalc
   ```

### macOS

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/musicalc.git
   cd musicalc
   ```

2. **Install Xcode Command Line Tools** (if not already installed)
   ```bash
   xcode-select --install
   ```

3. **Install Go dependencies**
   ```bash
   go mod download
   ```

4. **Build the application**
   ```bash
   go build -ldflags="-s -w" -o musicalc
   ```

5. **Run the application**
   ```bash
   ./musicalc
   ```

   **Optional**: Create an app bundle
   ```bash
   # Install fyne command
   go install fyne.io/fyne/v2/cmd/fyne@latest
   
   # Package as macOS app
   fyne package -os darwin -icon Icon.png
   ```

## Requirements

- Go 1.24.5 or later
- GCC/MinGW (for CGO support on Windows)
- Fyne v2.7.1

## License

MIT License - Copyright (c) 2026 B. Zeiss

See [LICENSE](LICENSE) for details.
