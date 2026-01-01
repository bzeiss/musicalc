# Icon Usage Guide

## Two Separate Icon Files Needed

### 1. Installer Icon (Windows Only)
- **File**: `icons/appicon.ico`
- **Used by**: Inno Setup installer, Windows shortcuts
- **Format**: ICO (Windows icon format with multiple sizes)
- **Platforms**: Windows only
- **Usage**: Inno Setup script references this for installer UI

### 2. Application Runtime Icon (Cross-Platform)
- **File**: `icons/appicon.png`
- **Used by**: Running application window (title bar, taskbar, dock)
- **Format**: PNG (recommended 512x512)
- **Platforms**: Windows, Linux, macOS
- **Usage**: Loaded in Go code via `window.SetIcon()`

## Creating the Icons

### From SVG to PNG (for app runtime)
You have `icons/appicon.svg`, convert it with Inkscape to PNG 512x512. Make sure the background is transparent.

### From SVG to ICO (for Windows installer)
**Online converter (easiest):**
1. https://convertio.co/svg-ico/
2. Upload `icons/appicon.svg`
3. Select multiple sizes: 16, 32, 48, 64, 128, 256
4. Download as `appicon.ico`

## Platform-Specific Behavior

### Windows
- **Installer shortcuts**: Use `.ico` file from Inno Setup
- **Running app**: Uses `.png` loaded in code
- **Taskbar**: Shows PNG icon

### Linux
- **Desktop file**: References PNG icon path
- **Window manager**: Shows PNG icon
- **App launcher**: Shows PNG icon

### macOS
- **App bundle**: Needs `.icns` file (if packaging as .app)
- **Window**: Shows PNG icon
- **Dock**: Shows PNG icon or .icns if packaged

## Summary

| Icon Type | File | Format | Usage |
|-----------|------|--------|-------|
| Installer | `icons/appicon.ico` | ICO | Windows installer only |
| Runtime | `icons/appicon.png` | PNG | Application window (all platforms) |
| macOS Bundle | `icons/appicon.icns` | ICNS | macOS .app packaging (optional) |
