# Build Instructions

## Windows

### Prerequisites

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

## Linux

### Debian/Ubuntu

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

### Fedora

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

## macOS

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
