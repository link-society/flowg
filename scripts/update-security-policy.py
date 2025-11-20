#!/usr/bin/env python3

from argparse import ArgumentParser
from packaging.version import Version

import subprocess


REPOSITORY = "https://github.com/link-society/flowg"
MAJOR_SUPPORT_WINDOW = 3

SECURITY_POLICY_HEADER = f"""
# Security Policy

## Reporting a Vulnerability

Please open a `Security Bug Report` on the
[bug tracker]({REPOSITORY}/issues/new/choose).

## Supported Versions

| Version | Supported |
| --- | --- |
"""


def get_tags() -> list[Version]:
    output = subprocess.check_output(
        "git tag --list --sort=-v:refname",
        text=True,
        shell=True,
    )

    lines = [
        line.strip()
        for line in output.splitlines()
    ]
    tags = [
        Version(tag[1:])
        for tag in lines
        if tag.startswith("v")
    ]

    return tags


def is_supported(
    t: Version,
    max_major: int,
    latest_by_major: dict[int, Version],
) -> bool:
    has_stable = max_major >= 1
    min_supported_major = max(1, max_major - (MAJOR_SUPPORT_WINDOW - 1))
    supported = False

    if any([
        not has_stable and t.major == 0,
        has_stable and t.major >= min_supported_major,
    ]):
        supported = (t == latest_by_major[t.major])

    return supported


def get_supported_versions(tags: list[Version]) -> list[tuple[Version, bool]]:
    assert len(tags) > 0

    latest_by_major: dict[int, Version] = {}

    for t in tags:
        m = t.major
        if m not in latest_by_major or t > latest_by_major[m]:
            latest_by_major[m] = t

    max_major = max(t.major for t in tags)

    return [
        (t, is_supported(t, max_major, latest_by_major))
        for t in tags
    ]


def main():
    parser = ArgumentParser()
    parser.add_argument("--next-release", nargs=1, required=True)
    args = parser.parse_args()

    next_release = Version(args.next_release[0])
    tags = get_tags()
    tags.insert(0, next_release)

    supported_versions = get_supported_versions(tags)

    print(SECURITY_POLICY_HEADER.strip())
    for t, supported in supported_versions:
        icon = ":white_check_mark:" if supported else ":x:"
        print(f"| [{t}]({REPOSITORY}/releases/tag/v{t}) | {icon} |")


if __name__ == "__main__":
    main()
