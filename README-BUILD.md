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
   $env:CGO_FLAGS="-O3 -flto=auto -march=x86-64-v3"
   $env:CGO_LDFLAGS="-O3 -flto=auto"
   $env:CGO_ENABLED=1
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
   export CGO_CFLAGS="-O3 -flto=auto -march=x86-64-v3"
   export CGO_LDFLAGS="-O3 -flto=auto"
   export CGO_ENABLED=1
   export CC=gcc
   export CXX=g++
   export GOOS=linux
   export GOARCH=amd64
   go build -ldflags="-s -w" -o musicalc_linux_amd64
   ```
7. **Build the application for Linux ARM64**
   ```bash
   export CGO_CFLAGS="-O3 -flto=auto -march=armv8.4-a+crc+crypto -fomit-frame-pointer"
   export CGO_LDFLAGS="-O3 -flto=auto -Wl,--gc-sections"
   export CGO_ENABLED=1
   export CC=aarch64-linux-gnu-gcc
   export PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig:/usr/share/pkgconfig
   export GOOS=linux
   export GOARCH=arm64
   go build -ldflags="-s -w" -o musicalc_linux_arm64
   ```

8. **Build the application for Windows AMD64**
   ```bash
   export CGO_CFLAGS="-O3 -flto=auto -march=x86-64-v3 -m64"
   export CGO_LDFLAGS="-O3 -flto=auto"
   export CGO_ENABLED=1
   export CC=x86_64-w64-mingw32-gcc
   export CXX=x86_64-w64-mingw32-g++
   export GOOS=windows
   export GOARCH=amd64
   go build -ldflags="-s -w" -o musicalc_win_amd64
   ```

9. **Build the application for Windows ARM64**
   ```bash
   export CGO_CFLAGS="-O3 -mcpu=oryon_1 -fomit-frame-pointer"
   export CGO_LDFLAGS="-O3"
   export CGO_ENABLED=1
   export CC="zig cc -target aarch64-windows-gnu"
   export CXX="zig c++ -target aarch64-windows-gnu"
   export GOOS=windows
   export GOARCH=arm64
   go build -ldflags="-s -w" -o musicalc_win_arm64
   ```

10. **Build the application for Android ARM64**
   ```bash
   export ANDROID_NDK_HOME=/path/to/your/android-ndk
   export ANDROID_HOME=/path/to/your/android-sdk
   export ANDROID_SDK_ROOT=$ANDROID_HOME
   export PATH=$PATH:${ANDROID_HOME}/cmdline-tools/latest/bin
   ./build-android.sh
   ```

   For the Android SDK and NDK, you can use the Android Studio SDK Manager to install them.
   ```bash
   sdkmanager --licenses # accept all licenses
   sdkmanager "platform-tools" "build-tools;36.1.0" "platforms;android-36" # choose the most recent version
   ```

11. **Run the application**
   ```bash
   ./musicalc_xxx
   ```

## Requirements

- Go 1.24.5 or later
- GCC/MinGW (for CGO support on Windows)
- Fyne v2.7.1

- GCC/MinGW (for CGO support on Windows)
- Fyne v2.7.1
