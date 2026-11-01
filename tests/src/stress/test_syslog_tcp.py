import pytest

import threading
import time

from ..settings import has_test_suite
from .._lib.stress_utils import (
    LIMITS,
    PROFILE_SECONDS,
    RECORD,
    TARGET_RECORDS_PER_SECOND,
    append_result,
    read_cpu_usage,
    read_memory_peak,
    read_stream_log_total,
    save_profile,
)


LOAD_SECONDS = 30
DRAIN_SECONDS = 5

# RFC 5424 with non-transparent framing (one message per line)
SYSLOG_MESSAGE = f"<14>1 - stress stress - - - {RECORD['message']}"

# Runs in a container on the test network. TCP pushes back when flowg is slower
# than the target rate, so the number of messages sent (printed) can be lower.
SENDER_SCRIPT = """
import socket, sys, time

host, duration, rate, message = sys.argv[1], float(sys.argv[2]), int(sys.argv[3]), sys.argv[4].encode()
frame = message + b"\\n"
sock = socket.create_connection((host, 5514))
start = time.perf_counter()
sent = 0

while (elapsed := time.perf_counter() - start) < duration:
    pending = int(elapsed * rate) - sent
    if pending > 0:
        sock.sendall(frame * pending)
        sent += pending

    time.sleep(0.001)

sock.close()
print(sent)
"""


@pytest.fixture()
def flowg_stress_environment():
    return {"FLOWG_SYSLOG_PROTOCOL": "tcp"}


@pytest.mark.skipif(not has_test_suite("stress"), reason="Stress test suite excluded")
@pytest.mark.parametrize("limits", LIMITS, ids=lambda limits: limits.name)
def test_stress_syslog_tcp(limits, flowg_stress_container, flowg_stress_config, flowg_network, docker_client, report_dir):
    profile_dir = report_dir / f"{limits.name}-syslog-tcp"
    profile_dir.mkdir()

    stored_before = read_stream_log_total()
    cpu_usage_before = read_cpu_usage(flowg_stress_container)

    sent = []
    load = threading.Thread(target=lambda: sent.append(int(docker_client.containers.run(
        "python:3.13-alpine",
        [
            "python", "-c", SENDER_SCRIPT,
            flowg_stress_container.name, str(LOAD_SECONDS), str(TARGET_RECORDS_PER_SECOND), SYSLOG_MESSAGE,
        ],
        network=flowg_network.name,
        remove=True,
    ))))
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

    print(f"{limits.name} (syslog tcp): {achieved:.0f} records/s, fits={fits}")

    append_result(report_dir / "syslog-tcp.csv", {
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
