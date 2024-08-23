#!/usr/bin/env python3

from argparse import ArgumentParser
import tomllib

from datetime import datetime
from time import sleep, time
import random

from urllib.request import Request, urlopen
import json


def format_syslog(hostname: str, appname: str, message: str) -> str:
    timestamp = datetime.now().strftime("%b %d %H:%M:%S")
    return f"{timestamp} {hostname} {appname}: {message}"


def send_log(log: str):
    payload = json.dumps({"record": {"message": log}}).encode("utf-8")
    req = Request(
        "http://localhost:5080/api/v1/pipelines/test/logs",
        data=payload,
        headers={"Content-Type": "application/json"},
    )

    with urlopen(req) as resp:
        assert resp.status == 200


def main():
    parser = ArgumentParser()
    parser.add_argument("--conftest", default="conftest.toml")
    parser.add_argument("--req-per-sec", type=int, default=100)
    parser.add_argument("--req-count", type=int, default=1_000_000)
    args = parser.parse_args()

    with open(args.conftest, "rb") as f:
        data = tomllib.load(f)

    req_count = args.req_count
    while req_count > 0:
        perc = (args.req_count - req_count) / args.req_count * 100
        print(f"Progress: {perc:.2f}%", end="")

        start = time()

        batch_count = 0
        while batch_count < args.req_per_sec:
            hostname = random.choice(data["hosts"])
            app = random.choice(data["apps"])
            message = random.choice(app["messages"])

            send_log(format_syslog(hostname, app["name"], message))

            batch_count += 1

        req_count = max(0, req_count - batch_count)

        end = time()
        print(f" - {batch_count} logs in {end - start:.2f}s")

        sleep(1)


if __name__ == "__main__":
    main()
