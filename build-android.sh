#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# --- 1, 2, 3) Check Environment Variables ---
check_env_var() {
    if [ -z "${!1}" ]; then
        echo "error: $1 is not set."
        exit 1
    fi
    if [ ! -d "${!1}" ]; then
        echo "error: $1 path does not exist: ${!1}"
        exit 1
    fi
}

check_env_var "ANDROID_NDK_HOME"
check_env_var "ANDROID_HOME"
check_env_var "ANDROID_SDK_ROOT"

# --- 4, 5) Locate zipalign ---
# Priority: 1. System Path, 2. build-tools (standard), 3. cmdline-tools (user request)
if command -v zipalign >/dev/null 2>&1; then
    ZIPALIGN_PATH=$(command -v zipalign)
else
    # Search for the highest version available in build-tools (Standard location)
    # or the user-specified cmdline-tools path
    SEARCH_PATHS=(
        "${ANDROID_HOME}/build-tools"
        "${ANDROID_HOME}/cmdline-tools/latest/bin"
    )
    
    ZIPALIGN_PATH=$(find "${SEARCH_PATHS[@]}" -name zipalign -type f 2>/dev/null | sort -V | tail -n 1)
fi

if [ -z "$ZIPALIGN_PATH" ] || [ ! -x "$ZIPALIGN_PATH" ]; then
    echo "error: zipalign not found in PATH or Android SDK structure."
    exit 1
fi

echo "Using zipalign: $ZIPALIGN_PATH"

# --- Build Configuration ---
export CGO_ENABLED=1
export CGO_CFLAGS="-O3 -flto=auto -march=armv8-a+crc+crypto"
# Force 16KB Page Alignment for Android 15+ compatibility
export CGO_LDFLAGS="-O3 -flto=auto -Wl,-z,max-page-size=16384 -Wl,-z,common-page-size=16384"
export GOFLAGS="-ldflags=-s -w"

git pull && git checkout -- go.mod && git fetch --tags --force

# --- Execution ---
echo "Building APK..."
fyne package --os android/arm64 --id com.github.bzeiss --release -icon icons/appicon.png

# Ensure dist exists
mkdir -p dist

# Aligning
echo "Aligning APK to 16KB boundary..."
mv musicalc.apk dist/musicalc-unaligned.apk
"$ZIPALIGN_PATH" -f -v -P 16 4 dist/musicalc-unaligned.apk dist/musicalc.apk

# Verify alignment
"$ZIPALIGN_PATH" -c -v -P 16 4 dist/musicalc.apk

# --- 6) Smarter Alignment Check ---
echo "Verifying 16KB ($2^{14}$) alignment of internal shared libraries..."
cd dist
unzip -p musicalc.apk lib/arm64-v8a/libmusicalc.so > tmp_lib.so

# Extraction of alignment value using awk
# We check the 'align' value of the LOAD segments
ALIGNMENT=$(objdump -p tmp_lib.so | grep "LOAD" | awk '{print $NF}' | head -n 1)

if [[ "$ALIGNMENT" == *"2**14"* ]] || [[ "$ALIGNMENT" == *"0x4000"* ]]; then
    echo "SUCCESS: Library is correctly aligned to 16KB ($ALIGNMENT)."
else
    echo "FAILURE: Library alignment is $ALIGNMENT, expected 2**14 (16384)."
    rm tmp_lib.so
    exit 1
fi

rm tmp_lib.so
echo "Build and Verification Complete: dist/musicalc.apk"