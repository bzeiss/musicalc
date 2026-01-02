# Build Instructions

## Windows

### Prerequisites

1. **Clone the repository**
   ```powershell
   git clone https://github.com/bzeiss/musicalc.git
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
   git clone https://github.com/bzeiss/musicalc.git
   cd musicalc
   ```

2. **Install dependencies**
   ```bash
   sudo apt-get update
   sudo apt-get install gcc libgl1-mesa-dev xorg-dev
   ```
3. **Linux ARM64 Cross-compilation support (requires a more recent Ubuntu Server version for building)**
   ```bash
   sudo dpkg --add-architecture arm64
   sudo apt-get update
   sudo apt install gcc-aarch64-linux-gnu libasound2-dev:arm64 libgl1-mesa-dev:arm64 libx11-dev:arm64 libxrandr-dev:arm64 libxxf86vm-dev:arm64 libxi-dev:arm64 libxcursor-dev:arm64 libxinerama-dev:arm64
   ```

3. **Windows AMD64 Cross-compilation support (requires a more recent Ubuntu Server version for building)**
   ```bash
   sudo apt-get update
   sudo apt install gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64
   ```

4. **Windows ARM64 Cross-compilation support (requires a more recent Ubuntu Server version for building)**
   - Download Zig from: https://ziglang.org/download/
   - Put zig somewere into your environment PATH
   - test by calling "zig" and "zig c++"

5. **Install Go dependencies**
   ```bash
   go mod download
   ```

6. **Build the application**
   ```bash
   go build -ldflags="-s -w" -o musicalc
   ```

7. **Run the application**
   ```bash
   ./musicalc
   ```

## Requirements

- Go 1.24.5 or later
- GCC/MinGW (for CGO support on Windows)
- Fyne v2.7.1
