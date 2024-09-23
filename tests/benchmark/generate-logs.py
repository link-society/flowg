from argparse import ArgumentParser
import tomllib
import os

from datetime import datetime
from time import perf_counter_ns
import multiprocessing
import random

import requests
import json

from tqdm import tqdm


def format_syslog(hostname: str, appname: str, message: str) -> str:
    timestamp = datetime.now().strftime("%b %d %H:%M:%S")
    return f"{timestamp} {hostname} {appname}: {message}"


def send_log(token: str, log: str):
    payload = json.dumps({"record": {"message": log}}).encode("utf-8")
    resp = requests.post(
        "http://localhost:5080/api/v1/pipelines/test/logs",
        data=payload,
        headers={
            "Content-Type": "application/json",
            "Authorization": f"Bearer pat:{token}",
        },
    )
    resp.raise_for_status()


def iteration(args: tuple[str, dict]):
    token, data = args

    hostname = random.choice(data["hosts"])
    app = random.choice(data["apps"])
    message = random.choice(app["messages"])

    send_log(token, format_syslog(hostname, app["name"], message))


def main():
    parser = ArgumentParser()
    parser.add_argument("--token", required=True)
    parser.add_argument("--conftest", default="testdata.toml")

    parser.add_argument(
        "--repeat",
        type=int,
        default=int(os.getenv("BENCHMARK_ITERATIONS", 1_000_000)),
    )
    parser.add_argument(
        "--jobs",
        type=int,
        default=int(os.getenv("BENCHMARK_JOBS", multiprocessing.cpu_count())),
    )

    args = parser.parse_args()

    with open(args.conftest, "rb") as f:
        data = tomllib.load(f)

    start_time = perf_counter_ns()

    with multiprocessing.Pool(args.jobs) as pool:
        params = [(args.token, data) for _ in range(args.repeat)]
        iterable = tqdm(
            pool.imap_unordered(iteration, params),
            total=args.repeat,
            desc="Sending logs",
            unit="req",
        )

        for _ in iterable:
            pass

    end_time = perf_counter_ns()

    total_time_ns = end_time - start_time
    total_time_s = total_time_ns / 1_000_000_000
    rate = args.repeat / total_time_s

    print(f"Requests sent: {args.repeat}")
    print(f"Total time:    {total_time_s:.2f}s")
    print(f"Rate:          {rate:.2f} req/s")


if __name__ == "__main__":
    main()
