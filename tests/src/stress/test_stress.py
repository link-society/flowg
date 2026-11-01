import pytest

from dataclasses import dataclass
import csv
import json
import subprocess

import requests

from ..settings import has_test_suite


LOAD_DURATION = "30s"
PROFILE_SECONDS = 20
CONNECTIONS = 16

RECORDS_PER_REQUEST = 10
TARGET_RECORDS_PER_SECOND = 10_000


@dataclass(frozen=True)
class Limits:
    name: str
    nano_cpus: int
    memory: str
    memtable_size: int


# Container sizes, the smallest being the single-node target on Kubernetes
LIMITS = [
    Limits("128m-256MB", 128_000_000, "256m", 8),
    Limits("256m-512MB", 256_000_000, "512m", 16),
    Limits("512m-1024MB", 512_000_000, "1024m", 32),
    Limits("1000m-2048MB", 1_000_000_000, "2048m", 64),
]

RECORD = {
    "message": json.dumps({
        "level": "info",
        "service": "api",
        "msg": "request completed",
        "duration_ms": "42",
    }),
}


@pytest.mark.skipif(not has_test_suite("stress"), reason="Stress test suite excluded")
@pytest.mark.parametrize("limits", LIMITS, ids=lambda limits: limits.name)
def test_stress(limits, flowg_stress_container, flowg_stress_token, flowg_stress_config, report_dir):
    profile_dir = report_dir / limits.name
    profile_dir.mkdir()

    body = json.dumps({"records": [RECORD] * RECORDS_PER_REQUEST})

    load = subprocess.Popen(
        [
            "oha",
            "--no-tui",
            "--output-format", "json",
            "-z", LOAD_DURATION,
            "--wait-ongoing-requests-after-deadline",
            "-q", str(TARGET_RECORDS_PER_SECOND // RECORDS_PER_REQUEST),
            "-c", str(CONNECTIONS),
            "-m", "POST",
            "-H", f"Authorization: Bearer {flowg_stress_token}",
            "-H", "Content-Type: application/json",
            "-d", body,
            "http://localhost:5080/api/v1/pipelines/transform/logs/struct",
        ],
        stdout=subprocess.PIPE,
    )

    save_profile(f"profile?seconds={PROFILE_SECONDS}", profile_dir / "cpu.pb.gz")
    stdout, _ = load.communicate()
    save_profile("heap", profile_dir / "heap.pb.gz")

    assert load.returncode == 0
    report = json.loads(stdout)

    flowg_stress_container.reload()
    state = flowg_stress_container.attrs["State"]

    statuses = report["statusCodeDistribution"]
    accepted = statuses.get("200", 0)
    failed = sum(statuses.values()) - accepted + sum(report["errorDistribution"].values())
    achieved = accepted * RECORDS_PER_REQUEST / report["summary"]["total"]
    # oha's total duration includes the wait for in-flight requests, allow 1% for it
    reached = achieved >= TARGET_RECORDS_PER_SECOND * 0.99
    fits = state["Running"] and not state["OOMKilled"] and failed == 0 and reached

    print(f"{limits.name}: {achieved:.0f} records/s, fits={fits}")

    results = report_dir / "results.csv"
    write_header = not results.exists()

    with results.open("a", newline="") as f:
        writer = csv.writer(f)

        if write_header:
            writer.writerow([
                "limits",
                "memtable_size_mb",
                "target_records_per_s",
                "achieved_records_per_s",
                "fits",
                "oom_killed",
                "failed_requests",
            ])

        writer.writerow([limits.name, limits.memtable_size, TARGET_RECORDS_PER_SECOND, round(achieved), fits, state["OOMKilled"], failed])


def save_profile(query, path):
    try:
        resp = requests.get(f"http://localhost:9113/debug/pprof/{query}", timeout=PROFILE_SECONDS + 10)

    except requests.ConnectionError:
        return  # flowg was killed, the results say so

    resp.raise_for_status()
    path.write_bytes(resp.content)
