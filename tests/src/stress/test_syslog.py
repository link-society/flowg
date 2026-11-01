import pytest

import socket
import threading
import time

import requests

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


LOAD_SECONDS = 30
DRAIN_SECONDS = 5

# RFC 5424, the record ends up in the "message" field like with the HTTP API
SYSLOG_MESSAGE = f"<14>1 - stress stress - - - {RECORD['message']}".encode()


@pytest.mark.skipif(not has_test_suite("stress"), reason="Stress test suite excluded")
@pytest.mark.parametrize("limits", LIMITS, ids=lambda limits: limits.name)
def test_stress_syslog(limits, flowg_stress_container, flowg_stress_config, docker_client, report_dir):
    profile_dir = report_dir / f"{limits.name}-syslog"
    profile_dir.mkdir()

    stored_before = read_stream_log_total()
    cpu_usage_before = read_cpu_usage(flowg_stress_container)

    sent = []
    load = threading.Thread(target=lambda: sent.append(send_syslog(LOAD_SECONDS, TARGET_RECORDS_PER_SECOND)))
    load.start()
    save_profile(f"profile?seconds={PROFILE_SECONDS}", profile_dir / "cpu.pb.gz")
    load.join()

    # records still queued in flowg when the load stops are counted too
    time.sleep(DRAIN_SECONDS)

    flowg_stress_container.reload()
    state = flowg_stress_container.attrs["State"]

    if state["Running"]:
        cpu_seconds = round((read_cpu_usage(flowg_stress_container) - cpu_usage_before) / 1e9, 1)
        stored = read_stream_log_total() - stored_before
    else:
        cpu_seconds = None
        stored = 0

    save_profile("heap", profile_dir / "heap.pb.gz")

    achieved = stored / LOAD_SECONDS
    reached = achieved >= TARGET_RECORDS_PER_SECOND * 0.99
    fits = state["Running"] and not state["OOMKilled"] and reached
    memory_peak = read_memory_peak(docker_client, flowg_stress_container) if state["Running"] else None

    print(f"{limits.name} (syslog): {achieved:.0f} records/s, fits={fits}")

    append_result(report_dir / "results-syslog.csv", {
        "limits": limits.name,
        "memtable_size_mb": limits.memtable_size,
        "target_records_per_s": TARGET_RECORDS_PER_SECOND,
        "sent_records_per_s": round(sent[0] / LOAD_SECONDS),
        "achieved_records_per_s": round(achieved),
        "fits": fits,
        "oom_killed": state["OOMKilled"],
        "memory_peak_mb": memory_peak,
        "cpu_seconds": cpu_seconds,
    })


def send_syslog(duration, rate):
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    start = time.perf_counter()
    sent = 0

    while (elapsed := time.perf_counter() - start) < duration:
        while sent < int(elapsed * rate):
            sock.sendto(SYSLOG_MESSAGE, ("127.0.0.1", 5514))
            sent += 1

        time.sleep(0.001)

    return sent


def read_stream_log_total():
    resp = requests.get("http://localhost:9113/metrics", timeout=10)
    resp.raise_for_status()

    for line in resp.text.splitlines():
        if line.startswith('flowg_stream_log_total{stream="stress"}'):
            return float(line.split()[-1])

    return 0
