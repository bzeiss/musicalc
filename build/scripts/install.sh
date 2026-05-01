#!/bin/bash
# Musicalc Linux Installer

# Move binary
sudo cp musicalc /usr/local/bin/

# Setup Icons
sudo mkdir -p /usr/local/share/icons/hicolor/512x512/apps/
sudo cp icon.png /usr/local/share/icons/hicolor/512x512/apps/musicalc.png

# Setup Desktop Entry
sudo mkdir -p /usr/local/share/applications/
sudo cp musicalc.desktop /usr/local/share/applications/

echo "----------------------------------------"
echo "MusiCalc has been installed!"
echo "You can now find it in your Start Menu."
echo "----------------------------------------"