# Creating Windows Installer

## Prerequisites

1. **Inno Setup**: Download and install from https://jrsoftware.org/isdl.php

## Building the Installer

1. **Build the application first:**
   ```powershell
   .\build\scripts\build-win.ps1
   ```

2. **Open Inno Setup Compiler**

3. **Compile the installer:**
   ```powershell
   .\build\scripts\build-installer.ps1
   ```

   For a release installer, first check out the exact version tag and then run:
   ```powershell
   .\build\scripts\build-installer.ps1 -Release
   ```

4. **Output:**
   - Installer will be created in `build/dist/installer/`

## Customization

Edit `build/installer/musicalc.iss` to change:
- `MyAppPublisher` - Your name/company
- `MyAppURL` - Your website/repository URL
- `AppId` - Unique GUID (keep as-is or generate new one)

## Testing

1. Run the generated installer from `build/dist/installer/`
2. Install to default location or custom directory
3. Verify desktop shortcut (if selected)
4. Test the application runs correctly
5. Uninstall via Windows Settings > Apps

## Distribution

The installer file is a single executable that can be distributed to users. It includes:
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

2. **Install Package Dependencies**
   See [Build Instructions](README-BUILD.md).

3. **Building packages with GoReleaser**

   For Linux AMD64:

   ```bash
   goreleaser check --config build/release/goreleaser-linux-amd64.yaml
   # for snapshot
   goreleaser release --snapshot --clean --config build/release/goreleaser-linux-amd64.yaml
   # for release
   goreleaser release --clean --config build/release/goreleaser-linux-amd64.yaml --skip=publish
   ```

   For Linux ARM64:

   ```bash
   goreleaser check --config build/release/goreleaser-linux-arm64.yaml
   # for snapshot
   goreleaser release --snapshot --clean --config build/release/goreleaser-linux-arm64.yaml
   # for release
   goreleaser release --clean --config build/release/goreleaser-linux-arm64.yaml --skip=publish
   ```

   For Windows AMD64/ARM64:
   ```bash
   goreleaser check --config build/release/goreleaser-win-all.yaml
   # for snapshot
   goreleaser release --snapshot --clean --config build/release/goreleaser-win-all.yaml
   # for release
   goreleaser release --clean --config build/release/goreleaser-win-all.yaml --skip=publish
   ```

   The Windows ARM64 build uses Zig for cross-compilation. The GoReleaser config redirects Zig cache files to `build/.cache/zig-global` and `build/.cache/zig-local`.

   Snapshot builds create packages in the `build/dist/` directory without requiring an exact Git tag. Release builds must be run from an exact version tag.

4. **Create a release (requires Git tag)**
   ```bash
   # Create and push tags manually; release tooling does not create or push tags.
   git tag -a 0.1.0 -m "Release 0.1.0"
   git push origin 0.1.0

   # Validate the current exact tag and all GoReleaser configs.
   python utils/release.py --check-only

   # Run one release config. Publishing is skipped unless --publish is supplied.
   python utils/release.py --release --config build/release/goreleaser-linux-amd64.yaml
   ```

**Generated artifacts**:
- `.tar.gz` archive with binary, install script, icon, and desktop entry
- `.deb` package for Debian/Ubuntu
- `.rpm` package for Fedora/RHEL
