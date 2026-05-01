# Build Instructions

## Windows

### Prerequisites

1. **Clone the repository**
   ```powershell
   git clone https://github.com/bzeiss/musicalc.git
   cd musicalc
   ```

2. **Install LLVM and MSYS2** (provides clang plus the MinGW linker required for Fyne)
   ```powershell
   winget install -e --id MSYS2.MSYS2
   winget install -e --id LLVM.LLVM
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

5. **Verify compiler installation**
   ```powershell
   clang --version
   x86_64-w64-mingw32-gcc --version
   ```
   Both commands should output version information.

6. **Install Go dependencies**
   ```powershell
   go mod download
   ```

7. **Build the application**
   ```powershell
   .\build\scripts\build-win.ps1
   ```

8. **Run the application**
   ```powershell
   .\build\dist\musicalc.exe
   ```

## Linux

### Ubuntu

1. **Clone the repository**
   ```bash
   git clone https://github.com/bzeiss/musicalc.git
   cd musicalc
   ```

2. **Install dependencies**
   ```bash
   sudo apt-get update
   sudo apt-get install gcc libgl1-mesa-dev xorg-dev libasound2-dev pkg-config
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

6. **Build the application for Linux AMD64**
   ```bash
   mkdir -p build/dist
   export CGO_CFLAGS="-O3 -flto=auto -march=x86-64-v3"
   export CGO_LDFLAGS="-O3 -flto=auto"
   export CGO_ENABLED=1
   export CC=gcc
   export CXX=g++
   export GOOS=linux
   export GOARCH=amd64
   go build -ldflags="-s -w" -o build/dist/musicalc_linux_amd64
   ```
7. **Build the application for Linux ARM64**
   ```bash
   mkdir -p build/dist
   export CGO_CFLAGS="-O3 -flto=auto -march=armv8.4-a+crc+crypto -fomit-frame-pointer"
   export CGO_LDFLAGS="-O3 -flto=auto -Wl,--gc-sections"
   export CGO_ENABLED=1
   export CC=aarch64-linux-gnu-gcc
   export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig:/usr/share/pkgconfig
   export GOOS=linux
   export GOARCH=arm64
   go build -ldflags="-s -w" -o build/dist/musicalc_linux_arm64
   ```

8. **Build the application for Windows AMD64**
   ```bash
   mkdir -p build/dist
   export CGO_CFLAGS="-O3 -flto=auto -march=x86-64-v3 -m64"
   export CGO_LDFLAGS="-O3 -flto=auto"
   export CGO_ENABLED=1
   export CC=x86_64-w64-mingw32-gcc
   export CXX=x86_64-w64-mingw32-g++
   export GOOS=windows
   export GOARCH=amd64
   go build -ldflags="-s -w" -o build/dist/musicalc_win_amd64.exe
   ```

9. **Build the application for Windows ARM64**
   ```bash
   mkdir -p build/dist
   export CGO_CFLAGS="-O3 -mcpu=oryon_1 -fomit-frame-pointer"
   export CGO_LDFLAGS="-O3"
   export CGO_ENABLED=1
   export CC="zig cc -target aarch64-windows-gnu"
   export CXX="zig c++ -target aarch64-windows-gnu"
   export GOOS=windows
   export GOARCH=arm64
   go build -ldflags="-s -w" -o build/dist/musicalc_win_arm64.exe
   ```

10. **Run the application**
   ```bash
   ./build/dist/musicalc_xxx
   ```

## Requirements

- Go 1.26.2 or later
- GCC/MinGW (for CGO support on Windows)
- Fyne v2.7.3
