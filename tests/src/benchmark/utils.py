from contextlib import contextmanager
from time import perf_counter_ns
from datetime import datetime
import random
import json

import requests


def format_syslog(hostname, appname, message):
    timestamp = datetime.now().strftime("%b %d %H:%M:%S")
    return f"{timestamp} {hostname} {appname}: {message}"


def send_log(token, log):
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


def iteration(args):
    token, data = args

    hostname = random.choice(data["hosts"])
    app = random.choice(data["apps"])
    message = random.choice(app["messages"])

    send_log(token, format_syslog(hostname, app["name"], message))


@contextmanager
def timeit(iteration_count):
    start_time = perf_counter_ns()
    yield
    end_time = perf_counter_ns()

    total_time_ns = end_time - start_time
    total_time_s = total_time_ns / 1_000_000_000
    rate = iteration_count / total_time_s

    print(f"Requests sent: {iteration_count}")
    print(f"Total time:    {total_time_s:.2f}s")
    print(f"Rate:          {rate:.2f} req/s")
