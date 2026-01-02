# Creating Windows Installer

## Prerequisites

1. **Inno Setup**: Download and install from https://jrsoftware.org/isdl.php

## Building the Installer

1. **Build the application first:**
   ```powershell
   go build -ldflags="-s -w" -o musicalc.exe
   ```

2. **Open Inno Setup Compiler**

3. **Open the script:**
   - File → Open → Select `musicalc.iss`

4. **Compile:**
   - Build → Compile (or press F9)

5. **Output:**
   - Installer will be created in `installer/MusiCalc-Setup-1.0.0.exe`

## Customization

Edit `musicalc.iss` to change:
- `MyAppVersion` - Application version number
- `MyAppPublisher` - Your name/company
- `MyAppURL` - Your website/repository URL
- `AppId` - Unique GUID (keep as-is or generate new one)

## Testing

1. Run the generated installer: `installer/MusiCalc-Setup-1.0.0.exe`
2. Install to default location or custom directory
3. Verify desktop shortcut (if selected)
4. Test the application runs correctly
5. Uninstall via Windows Settings → Apps

## Distribution

The installer file (`MusiCalc-Setup-1.0.0.exe`) is a single executable that can be distributed to users. It includes:
- Application executable
- App icon
- Start menu shortcuts
- Optional desktop shortcut
- Uninstaller

# Creating Linux Packages and Tarballs

GoReleaser automates the creation of `.deb`, `.rpm`, and `.tar.gz` packages for distribution.

1. **Install GoReleaser**
   ```bash
   go install github.com/goreleaser/goreleaser/v2@latest
   ```

2. **Validate configuration**
   ```bash
   goreleaser check
   ```

3. **Build snapshot packages (for testing)**
   ```bash
   goreleaser release --snapshot --clean
   ```

   for ARM64:

   ```bash
   goreleaser release --snapshot --clean --config .goreleaser-arm64.yaml
   ```

   This creates packages in the `dist/` directory without requiring a Git tag.

4. **Create a release (requires Git tag)**
   ```bash
   # Create and push a version tag
   git tag -a v0.1.0 -m "Release version 0.1.0"
   git push origin v0.1.0
   
   # Build and publish release
   goreleaser release --clean --skip=publish
   ```

   for ARM64:

   ```bash
   goreleaser release --clean --config .goreleaser-arm64.yaml
   ```

**Generated artifacts**:
- `.tar.gz` archive with binary, install script, icon, and desktop entry
- `.deb` package for Debian/Ubuntu
- `.rpm` package for Fedora/RHEL
