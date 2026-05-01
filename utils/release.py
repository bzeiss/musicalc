#!/usr/bin/env python3
"""
Release helper for MusiCalc.

Git tags are the only release version source. This script validates the
current tag and runs GoReleaser without creating commits, tags, or pushes.
"""

import argparse
import re
import subprocess
import sys
from pathlib import Path


DEFAULT_CONFIGS = [
    Path("build/release/goreleaser-linux-amd64.yaml"),
    Path("build/release/goreleaser-linux-arm64.yaml"),
    Path("build/release/goreleaser-win-all.yaml"),
]

VERSION_RE = re.compile(r"^\d+\.\d+\.\d+$")


def find_project_root() -> Path:
    current = Path.cwd().resolve()
    for directory in [current] + list(current.parents):
        if (directory / "go.mod").exists():
            return directory

    script_dir = Path(__file__).parent.resolve()
    for directory in [script_dir] + list(script_dir.parents):
        if (directory / "go.mod").exists():
            return directory

    raise SystemExit("Error: could not find project root (go.mod not found)")


def run(args: list[str], cwd: Path, description: str, capture: bool = False) -> str:
    print(f"\n{description}...")
    result = subprocess.run(
        args,
        cwd=cwd,
        check=False,
        capture_output=capture,
        text=True,
    )
    if result.returncode != 0:
        if capture and result.stderr:
            print(result.stderr.strip())
        raise SystemExit(f"Error: {description} failed")
    if capture:
        return result.stdout.strip()
    return ""


def exact_version_tag(project_root: Path) -> str:
    tag = run(
        ["git", "describe", "--tags", "--exact-match"],
        project_root,
        "Reading exact Git tag",
        capture=True,
    )
    if not VERSION_RE.match(tag):
        raise SystemExit(
            f"Error: release tag '{tag}' is invalid; expected MAJOR.MINOR.PATCH, for example 0.8.7"
        )
    return tag


def check_configs(project_root: Path, configs: list[Path]) -> None:
    for config in configs:
        run(
            ["goreleaser", "check", "--config", str(config)],
            project_root,
            f"Checking {config}",
        )


def release_configs(project_root: Path, configs: list[Path], skip_publish: bool) -> None:
    for config in configs:
        command = ["goreleaser", "release", "--clean", "--config", str(config)]
        if skip_publish:
            command.append("--skip=publish")
        run(command, project_root, f"Running GoReleaser for {config}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Validate and run a tag-based MusiCalc release."
    )
    parser.add_argument(
        "--config",
        action="append",
        type=Path,
        dest="configs",
        help="GoReleaser config to use. May be specified multiple times for --check-only; release mode requires exactly one.",
    )
    parser.add_argument(
        "--check-only",
        action="store_true",
        help="Only validate the current tag and GoReleaser configs. This is the default unless --release is set.",
    )
    parser.add_argument(
        "--release",
        action="store_true",
        help="Run GoReleaser after validation. Requires exactly one --config.",
    )
    parser.add_argument(
        "--publish",
        action="store_true",
        help="Allow GoReleaser to publish. By default, publishing is skipped.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    project_root = find_project_root()
    configs = args.configs or DEFAULT_CONFIGS

    tag = exact_version_tag(project_root)
    print(f"Release version: {tag}")
    print("Git tags are user-managed; this script will not create or push tags.")

    check_configs(project_root, configs)

    if args.check_only or not args.release:
        print("\nValidation complete.")
        return

    if len(configs) != 1:
        raise SystemExit(
            "Error: release mode requires exactly one --config. Use --check-only to validate all configs."
        )

    release_configs(project_root, configs, skip_publish=not args.publish)
    print("\nRelease commands completed.")


if __name__ == "__main__":
    main()
