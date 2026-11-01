from dataclasses import dataclass
import csv
import json

import requests


PROFILE_SECONDS = 20
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
    # default memtable size, to compare with the tuned one above
    Limits("512m-1024MB-memtable64", 512_000_000, "1024m", 64),
]

RECORD = {
    "message": json.dumps({
        "level": "info",
        "service": "api",
        "msg": "request completed",
        "duration_ms": "42",
    }),
}


def append_result(path, row):
    write_header = not path.exists()

    with path.open("a", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=row.keys())

        if write_header:
            writer.writeheader()

        writer.writerow(row)


def read_cpu_usage(container):
    stats = container.stats(stream=False, one_shot=True)
    return stats["cpu_stats"]["cpu_usage"]["total_usage"]


def read_memory_peak(docker_client, container):
    peak = read_from_sidecar(docker_client, container, "/proc/1/root/sys/fs/cgroup/memory.peak")
    return int(peak) // 2**20


def read_from_sidecar(docker_client, container, path):
    # the flowg image has no shell, read from a sidecar sharing its PID namespace, flowg being PID 1
    return docker_client.containers.run(
        "busybox",
        ["cat", path],
        pid_mode=f"container:{container.id}",
        cap_add=["SYS_PTRACE"],
        remove=True,
    ).decode()


def save_profile(query, path):
    try:
        resp = requests.get(f"http://localhost:9113/debug/pprof/{query}", timeout=PROFILE_SECONDS + 10)

    except requests.ConnectionError:
        return  # flowg was killed, the results say so

    resp.raise_for_status()
    path.write_bytes(resp.content)
