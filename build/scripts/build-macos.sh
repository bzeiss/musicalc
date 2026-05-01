#!/usr/bin/env bash
set -euo pipefail

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${project_root}"

version="${MUSICALC_VERSION:-}"
if [[ -z "${version}" ]]; then
  if version="$(git describe --tags --exact-match 2>/dev/null)"; then
    :
  else
    short_commit="$(git rev-parse --short HEAD 2>/dev/null || true)"
    if [[ -n "${short_commit}" ]]; then
      version="dev-snapshot-${short_commit}"
    else
      version="dev"
    fi
  fi
fi

bundle_short_version="${version}"
semver_prefix='^([0-9]+[.][0-9]+[.][0-9]+)'
if [[ "${bundle_short_version}" =~ ${semver_prefix} ]]; then
  bundle_short_version="${BASH_REMATCH[1]}"
else
  bundle_short_version="0.0.0"
fi

dist_dir="build/dist/macos"
work_dir="build/.cache/macos"
arm64_binary="${work_dir}/musicalc_arm64"
amd64_binary="${work_dir}/musicalc_amd64"
app_dir="${dist_dir}/MusiCalc.app"
contents_dir="${app_dir}/Contents"
macos_dir="${contents_dir}/MacOS"
resources_dir="${contents_dir}/Resources"
iconset_dir="${work_dir}/appicon.iconset"
dmg_root="${work_dir}/dmg-root"
dmg_path="${dist_dir}/musicalc_${version}_universal.dmg"
zip_path="${dist_dir}/MusiCalc_${version}_universal.app.zip"

rm -rf "${dist_dir}" "${work_dir}"
mkdir -p "${dist_dir}" "${work_dir}" "${macos_dir}" "${resources_dir}" "${iconset_dir}" "${dmg_root}"

ldflags="-s -w -X main.version=${version}"

echo "Building darwin/arm64"
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="clang -arch arm64" CXX="clang++ -arch arm64" \
  go build -trimpath -ldflags="${ldflags}" -o "${arm64_binary}" .

echo "Building darwin/amd64"
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC="clang -arch x86_64" CXX="clang++ -arch x86_64" \
  go build -trimpath -ldflags="${ldflags}" -o "${amd64_binary}" .

echo "Creating universal binary"
lipo -create -output "${macos_dir}/MusiCalc" "${arm64_binary}" "${amd64_binary}"
chmod 0755 "${macos_dir}/MusiCalc"

echo "Creating app icon"
sips -z 16 16 icons/appicon.png --out "${iconset_dir}/icon_16x16.png" >/dev/null
sips -z 32 32 icons/appicon.png --out "${iconset_dir}/icon_16x16@2x.png" >/dev/null
sips -z 32 32 icons/appicon.png --out "${iconset_dir}/icon_32x32.png" >/dev/null
sips -z 64 64 icons/appicon.png --out "${iconset_dir}/icon_32x32@2x.png" >/dev/null
sips -z 128 128 icons/appicon.png --out "${iconset_dir}/icon_128x128.png" >/dev/null
sips -z 256 256 icons/appicon.png --out "${iconset_dir}/icon_128x128@2x.png" >/dev/null
sips -z 256 256 icons/appicon.png --out "${iconset_dir}/icon_256x256.png" >/dev/null
sips -z 512 512 icons/appicon.png --out "${iconset_dir}/icon_256x256@2x.png" >/dev/null
sips -z 512 512 icons/appicon.png --out "${iconset_dir}/icon_512x512.png" >/dev/null
sips -z 1024 1024 icons/appicon.png --out "${iconset_dir}/icon_512x512@2x.png" >/dev/null
iconutil -c icns "${iconset_dir}" -o "${resources_dir}/appicon.icns"

cat > "${contents_dir}/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleDevelopmentRegion</key>
  <string>en</string>
  <key>CFBundleDisplayName</key>
  <string>MusiCalc</string>
  <key>CFBundleExecutable</key>
  <string>MusiCalc</string>
  <key>CFBundleIconFile</key>
  <string>appicon</string>
  <key>CFBundleIdentifier</key>
  <string>com.musicalc</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundleName</key>
  <string>MusiCalc</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleShortVersionString</key>
  <string>${bundle_short_version}</string>
  <key>CFBundleVersion</key>
  <string>${bundle_short_version}</string>
  <key>LSMinimumSystemVersion</key>
  <string>11.0</string>
  <key>NSHighResolutionCapable</key>
  <true/>
</dict>
</plist>
EOF

plutil -lint "${contents_dir}/Info.plist"
lipo -info "${macos_dir}/MusiCalc"

echo "Creating app zip"
ditto -c -k --keepParent "${app_dir}" "${zip_path}"

echo "Creating DMG"
cp -R "${app_dir}" "${dmg_root}/"
ln -s /Applications "${dmg_root}/Applications"
hdiutil create -volname "MusiCalc" -srcfolder "${dmg_root}" -ov -format UDZO "${dmg_path}"
hdiutil verify "${dmg_path}"

echo "Created ${app_dir}"
echo "Created ${zip_path}"
echo "Created ${dmg_path}"
