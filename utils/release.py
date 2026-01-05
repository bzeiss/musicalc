#!/usr/bin/env python3
"""
Release script for MusiCalc
Interactive release workflow with version management and git automation
"""

import sys
import os
import re
import subprocess
from pathlib import Path


def find_project_root():
    """Find project root by looking for go.mod file"""
    current = Path.cwd().resolve()
    
    # Check current directory and all parents
    for directory in [current] + list(current.parents):
        if (directory / "go.mod").exists():
            return directory
    
    # If not found, check if we're in a subdirectory with ../go.mod
    script_dir = Path(__file__).parent.resolve()
    for directory in [script_dir] + list(script_dir.parents):
        if (directory / "go.mod").exists():
            return directory
    
    print("✗ Error: Could not find project root (go.mod not found)")
    sys.exit(1)


def read_current_version(project_root):
    """Read current version from VERSION file"""
    version_file = project_root / "VERSION"
    if version_file.exists():
        return version_file.read_text().strip()
    return None


def update_version_file(project_root, new_version):
    """Update VERSION file"""
    version_file = project_root / "VERSION"
    version_file.write_text(new_version + "\n")
    print(f"✓ Updated VERSION: {new_version}")


def update_inno_setup(project_root, new_version):
    """Update version in musicalc.iss Inno Setup script"""
    iss_file = project_root / "musicalc.iss"
    
    if not iss_file.exists():
        print(f"✗ Error: {iss_file} not found")
        return False
    
    content = iss_file.read_text(encoding='utf-8')
    
    # Update #define MyAppVersion
    pattern = r'(#define MyAppVersion\s+")[^"]+(")' 
    
    # Check if pattern exists in file
    if not re.search(pattern, content):
        print(f"✗ Error: Could not find version pattern in {iss_file}")
        return False
    
    replacement = r'\g<1>' + new_version + r'\g<2>'
    new_content = re.sub(pattern, replacement, content)
    
    # Only write if content changed
    if new_content != content:
        iss_file.write_text(new_content, encoding='utf-8')
        print(f"✓ Updated musicalc.iss: {new_version}")
    else:
        print(f"✓ musicalc.iss already at version: {new_version}")
    
    return True


def update_desktop_files(project_root, new_version):
    """Update version in .desktop files (StartupWMClass field)"""
    desktop_files = [
        project_root / "musicalc.desktop",
        project_root / "musicalc-pkg.desktop"
    ]
    
    success = True
    for desktop_file in desktop_files:
        if not desktop_file.exists():
            print(f"⚠ Warning: {desktop_file.name} not found, skipping")
            continue
        
        content = desktop_file.read_text(encoding='utf-8')
        
        # Update StartupWMClass=MusiCalc v0.8.4
        pattern = r'(StartupWMClass=MusiCalc v)[0-9.]+'
        
        # Check if pattern exists in file
        if not re.search(pattern, content):
            print(f"✗ Error: Could not find StartupWMClass pattern in {desktop_file.name}")
            success = False
            continue
        
        replacement = r'\g<1>' + new_version
        new_content = re.sub(pattern, replacement, content)
        
        # Only write if content changed
        if new_content != content:
            desktop_file.write_text(new_content, encoding='utf-8')
            print(f"✓ Updated {desktop_file.name}: {new_version}")
        else:
            print(f"✓ {desktop_file.name} already at version: {new_version}")
    
    return success


def validate_version(version_str):
    """Validate version string format (e.g., 0.8.3)"""
    pattern = r'^\d+\.\d+\.\d+$'
    return re.match(pattern, version_str) is not None


def run_command(cmd, cwd, description):
    """Run a shell command and return success status"""
    print(f"\n→ {description}...")
    try:
        result = subprocess.run(
            cmd,
            cwd=cwd,
            shell=True,
            check=True,
            capture_output=True,
            text=True
        )
        if result.stdout:
            print(result.stdout.strip())
        print(f"✓ {description} completed")
        return True
    except subprocess.CalledProcessError as e:
        print(f"✗ Error: {description} failed")
        if e.stderr:
            print(e.stderr.strip())
        return False


def main():
    # Find project root
    project_root = find_project_root()
    print(f"Project root: {project_root}")
    print()
    
    # Read current version
    current_version = read_current_version(project_root)
    if not current_version:
        print("✗ Error: VERSION file not found or empty")
        sys.exit(1)
    
    print(f"Current version: {current_version}")
    print()
    
    # Ask if user wants to keep or change version
    response = input(f"Keep version {current_version}? [Y/n]: ").strip().lower()
    
    new_version = current_version
    if response == 'n':
        while True:
            new_version = input("Enter new version (format: xx.yy.zz): ").strip()
            if validate_version(new_version):
                break
            print("✗ Invalid version format. Please use format like 0.8.4")
        
        # Update VERSION file
        update_version_file(project_root, new_version)
    else:
        print(f"Using version: {new_version}")
    
    print()
    
    # Run go mod tidy
    if not run_command("go mod tidy", project_root, "Running go mod tidy"):
        sys.exit(1)
    
    # Update musicalc.iss
    print()
    if not update_inno_setup(project_root, new_version):
        sys.exit(1)
    
    # Update .desktop files
    print()
    if not update_desktop_files(project_root, new_version):
        sys.exit(1)
    
    # Ask for commit message
    print()
    commit_message = input("Enter release commit message: ").strip()
    if not commit_message:
        commit_message = f"Release v{new_version}"
    
    # Git workflow
    print()
    print("═" * 50)
    print("Git Workflow")
    print("═" * 50)
    
    # Stage files
    if not run_command("git add VERSION musicalc.iss musicalc.desktop musicalc-pkg.desktop go.mod go.sum", project_root, "Staging files"):
        sys.exit(1)
    
    # Check if there are changes to commit
    try:
        result = subprocess.run(
            "git diff --cached --quiet",
            cwd=project_root,
            shell=True
        )
        has_changes = result.returncode != 0
    except subprocess.CalledProcessError:
        has_changes = True
    
    # Commit and push if there are changes
    if has_changes:
        if not run_command(f'git commit -m "{commit_message}"', project_root, "Committing changes"):
            sys.exit(1)
        
        # Push commits
        if not run_command("git push", project_root, "Pushing commits to origin"):
            sys.exit(1)
    else:
        print("\n→ No changes to commit (version unchanged)")
        print("✓ Skipping commit step")
    
    # Create and push tag
    tag_name = f"{new_version}"
    if not run_command(f'git tag -a {tag_name} -m "Release {tag_name}"', project_root, f"Creating tag {tag_name}"):
        sys.exit(1)
    
    if not run_command(f"git push origin {tag_name}", project_root, f"Pushing tag {tag_name}"):
        sys.exit(1)
    
    # Success summary
    print()
    print("═" * 50)
    print("✓ Release Complete!")
    print("═" * 50)
    print(f"Version: {new_version}")
    print(f"Tag: {tag_name}")
    print(f"Commit: {commit_message}")
    print()
    print("Next steps:")
    print("  1. Build the application: go build -ldflags=\"-s -w -H=windowsgui\" -o musicalc.exe")
    print("  2. Test the application")
    print("  3. Create installer (if using Inno Setup)")


if __name__ == "__main__":
    main()
