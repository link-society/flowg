import pytest

from datetime import datetime
from shutil import rmtree
from pathlib import Path
import json
import os

import requests


@pytest.fixture(scope="module")
def cache_dir():
    cache_dir = Path.cwd() / "cache" / "benchmark"
    rmtree(cache_dir, ignore_errors=True)
    cache_dir.mkdir(parents=True)
    yield cache_dir
    rmtree(cache_dir, ignore_errors=True)


@pytest.fixture()
def iteration_count():
    return int(os.getenv("BENCHMARK_ITERATIONS", 100_000))


@pytest.fixture(scope="module")
def flowg_config(flowg_admin_token, spec_dir):
    pipelines = ["test"]

    for pipeline_name in pipelines:
        print(f"Creating pipeline: {pipeline_name}")

        with open(spec_dir / "config" / "pipelines" / f"{pipeline_name}.json") as f:
            pipeline = json.load(f)

            resp = requests.put(
                f"http://localhost:5080/api/v1/pipelines/{pipeline_name}",
                headers={"Authorization": f"Bearer {flowg_admin_token}"},
                json={"flow": pipeline},
            )
            resp.raise_for_status()

    transformers = ["apache", "json", "logfmt", "syslog"]

    for transformer_name in transformers:
        print(f"Creating transformer: {transformer_name}")

        with open(spec_dir / "config" / "transformers" / f"{transformer_name}.vrl") as f:
            transformer = f.read()

            resp = requests.put(
                f"http://localhost:5080/api/v1/transformers/{transformer_name}",
                headers={"Authorization": f"Bearer {flowg_admin_token}"},
                json={"script": transformer},
            )
            resp.raise_for_status()


@pytest.fixture()
def testdata(cache_dir):
    def format_syslog(hostname, appname, message):
        timestamp = datetime.now().strftime("%b %d %H:%M:%S")
        return f"{timestamp} {hostname} {appname}: {message}"

    hosts = [
        "localhost",
        "vm1.example.com",
        "vm2.example.com",
        "1.2.3.4",
    ]

    apps = [
        {
            "name": "iam01",
            "messages": [
                'level=info msg="User login successful"',
                'level=error msg="User login failed" error="invalid password"',
                'level=warn msg="File not found" file="example.txt"',
                'level=debug msg="Processing request" duration=250ms',
                'level=info msg="User logout successful"',
                'level=error msg="User logout failed" error="session expired"',
            ],
        },
        {
            "name": "iam02",
            "messages": [
                '{"level":"info","msg":"User login successful"}',
                '{"level":"error","msg":"User login failed","error":"invalid password"}',
                '{"level":"warn","msg":"File not found","file":"example.txt"}',
                '{"level":"debug","msg":"Processing request","duration":"250ms"}',
                '{"level":"info","msg":"User logout successful"}',
                '{"level":"error","msg":"User logout failed","error":"session expired"}',
            ],
        },
        {
            "name": "db",
            "messages": [
                'level=info msg="Database connection established" db="mysql" host="localhost" port=3306',
                'level=error msg="Database connection failed" db="mysql" host="localhost" port=3306 error="connection refused"',
                'level=warn msg="Database connection lost" db="mysql" host="localhost" port=3306',
                'level=debug msg="Database query executed" db="mysql" host="localhost" port=3306 query="SELECT * FROM users"',
                'level=info msg="Database connection closed" db="mysql" host="localhost" port=3306',
            ],
        },
        {
            "name": "apache2",
            "messages": [
                '192.168.1.1 - - [23/Aug/2024:14:55:31 +0000] "GET /index.html HTTP/1.1" 200 1234',
                '192.168.1.2 - - [23/Aug/2024:14:56:12 +0000] "POST /login HTTP/1.1" 302 546',
                '192.168.1.3 - - [23/Aug/2024:14:57:45 +0000] "GET /about-us HTTP/1.1" 404 321',
                '192.168.1.4 - - [23/Aug/2024:14:58:02 +0000] "GET /contact HTTP/1.1" 200 789',
                '192.168.1.5 - - [23/Aug/2024:14:58:56 +0000] "GET /nonexistentpage HTTP/1.1" 404 217',
                '192.168.1.6 - - [23/Aug/2024:14:59:32 +0000] "POST /api/data HTTP/1.1" 500 654',
                '192.168.1.7 - - [23/Aug/2024:15:00:01 +0000] "GET /blog HTTP/1.1" 301 123',
                '192.168.1.8 - - [23/Aug/2024:15:01:12 +0000] "PUT /update HTTP/1.1" 204 0',
                '192.168.1.9 - - [23/Aug/2024:15:01:57 +0000] "DELETE /delete-item HTTP/1.1" 403 342',
                '192.168.1.10 - - [23/Aug/2024:15:02:34 +0000] "GET /dashboard HTTP/1.1" 200 1567',
            ],
        },
    ]

    testdata_path = cache_dir / "testdata.txt"

    with testdata_path.open("w") as writer:
        for host in hosts:
            for app in apps:
                for message in app["messages"]:
                    content = format_syslog(host, app["name"], message)
                    payload = json.dumps({"records": [{"message": content}]})
                    print(payload, file=writer)

    yield testdata_path
