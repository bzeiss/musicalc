#!/usr/bin/env python3
"""
Generate tabicons.go by bundling SVG icon files using fyne bundle.
This script automates the process of creating Fyne resource bundles from SVG files.
"""

import os
import subprocess
import sys
from pathlib import Path

# Configuration constants
TARGET_PACKAGE = "ui"
TARGET_FILENAME = "tabicons.go"
OUTPUT_DIRECTORY = "internal/ui"


def get_project_root():
    """
    Find the project root directory relative to this script's location.
    Assumes script is in utils/ and project root is 1 level up.
    """
    script_dir = Path(__file__).resolve().parent
    project_root = script_dir.parent
    return project_root


def main():
    """Main execution function."""
    # Get project root
    project_root = get_project_root()
    print(f"Project root: {project_root}")
    
    # Construct paths
    svg_directory = project_root / OUTPUT_DIRECTORY
    output_file = svg_directory / TARGET_FILENAME
    
    # Verify SVG directory exists
    if not svg_directory.exists():
        print(f"Error: SVG directory does not exist: {svg_directory}")
        sys.exit(1)
    
    print(f"SVG directory: {svg_directory}")
    
    # Collect all SVG files from the directory
    svg_files = sorted([f.name for f in svg_directory.glob("*.svg")])
    
    if not svg_files:
        print(f"Error: No SVG files found in {svg_directory}")
        sys.exit(1)
    
    print(f"Found {len(svg_files)} SVG files: {', '.join(svg_files)}\n")
    
    # Remove existing output file if it exists
    if output_file.exists():
        output_file.unlink()
        print(f"Removed existing {TARGET_FILENAME}")
    
    # Change to the SVG directory
    os.chdir(svg_directory)
    print(f"Changed directory to: {svg_directory}\n")
    
    # Bundle SVG files
    for index, svg_file in enumerate(svg_files):
        is_first = (index == 0)
        
        # Build fyne bundle command
        cmd = ["fyne", "bundle", "-package", TARGET_PACKAGE]
        
        if not is_first:
            cmd.append("-append")
        
        cmd.extend(["-o", TARGET_FILENAME, svg_file])
        
        # Execute command
        print(f"Bundling {svg_file}{'...' if is_first else ' (append)...'}")
        try:
            result = subprocess.run(cmd, check=True, capture_output=True, text=True)
            if result.stdout:
                print(result.stdout)
        except subprocess.CalledProcessError as e:
            print(f"Error bundling {svg_file}: {e.stderr}")
            sys.exit(1)
    
    print(f"\n✓ Successfully generated {TARGET_FILENAME}")
    print(f"  Location: {output_file}")
    print(f"  Bundled {len(svg_files)} SVG icons")


if __name__ == "__main__":
    main()
