import pytest

import json
import subprocess

from ..settings import has_test_suite
from .._lib.stress_utils import (
    LIMITS,
    PROFILE_SECONDS,
    RECORD,
    TARGET_RECORDS_PER_SECOND,
    append_result,
    read_cpu_usage,
    read_memory_peak,
    save_profile,
)


LOAD_DURATION = "30s"
CONNECTIONS = 16
RECORDS_PER_REQUEST = 10


@pytest.mark.skipif(not has_test_suite("stress"), reason="Stress test suite excluded")
@pytest.mark.parametrize("limits", LIMITS, ids=lambda limits: limits.name)
def test_stress(limits, flowg_stress_container, flowg_stress_token, flowg_stress_config, docker_client, report_dir):
    profile_dir = report_dir / limits.name
    profile_dir.mkdir()

    body = json.dumps({"records": [RECORD] * RECORDS_PER_REQUEST})
    cpu_usage_before = read_cpu_usage(flowg_stress_container)

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

    flowg_stress_container.reload()
    state = flowg_stress_container.attrs["State"]

    if state["Running"]:
        cpu_seconds = round((read_cpu_usage(flowg_stress_container) - cpu_usage_before) / 1e9, 1)
    else:
        cpu_seconds = None

    save_profile("heap", profile_dir / "heap.pb.gz")

    assert load.returncode == 0
    report = json.loads(stdout)

    statuses = report["statusCodeDistribution"]
    accepted = statuses.get("200", 0)
    failed = sum(statuses.values()) - accepted + sum(report["errorDistribution"].values())
    achieved = accepted * RECORDS_PER_REQUEST / report["summary"]["total"]
    # oha's total duration includes the wait for in-flight requests, allow 1% for it
    reached = achieved >= TARGET_RECORDS_PER_SECOND * 0.99
    fits = state["Running"] and not state["OOMKilled"] and failed == 0 and reached
    memory_peak = read_memory_peak(docker_client, flowg_stress_container) if state["Running"] else None

    print(f"{limits.name}: {achieved:.0f} records/s, fits={fits}")

    append_result(report_dir / "http.csv", {
        "limits": limits.name,
        "memtable_size_mb": limits.memtable_size,
        "target_records_per_s": TARGET_RECORDS_PER_SECOND,
        "achieved_records_per_s": round(achieved),
        "fits": fits,
        "oom_killed": state["OOMKilled"],
        "failed_requests": failed,
        "memory_peak_mb": memory_peak,
        "cpu_seconds": cpu_seconds,
    })
