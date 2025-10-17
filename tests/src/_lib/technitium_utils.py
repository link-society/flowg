import pytest
import os
import time
import requests
from pathlib import Path

from contextlib import contextmanager

from . import docker_utils


TECHNITIUM_IMAGE = "technitium/dns-server:latest"
POWERDNS_IMAGE = "powerdns/pdns-auth-44"

@contextmanager
def container(docker_client, *, name, network, report_dir):
    print("Pulling the Power DNS Server image")
    try:
        docker_client.images.pull(POWERDNS_IMAGE)

    except Exception as err:
        pytest.fail(f"{err}", pytrace=False)

    try:
        print(f"Creating Container: {name}")
        container = docker_client.containers.run(
            image=POWERDNS_IMAGE,
            name=name,
            network=network.name,
            hostname=name,
            ports={
                "53/tcp": 5300,
                '5380/tcp': 5380
            },
            detach=True,
        )

    except Exception as err:
        pytest.fail(f"{err}", pytrace=True)

    try:
        print(f"Waiting for healthcheck: {name}")
        # wait_for_healthcheck(container)

    except RuntimeError as err:
        print(f"Container '{name}' healthcheck failed. Logs:")
        print(container.logs().decode('utf-8'))
        docker_utils.teardown_container(container, report_dir)
        pytest.fail(f"{err}", pytrace=False)

    yield

    docker_utils.teardown_container(container, report_dir)


def wait_for_healthcheck(container):
    attempts = 10
    interval = 1

    while True:
        try:
            resp = requests.get("http://localhost:5380/")
            resp.raise_for_status()

            print(f"Health check for {container.name} successful!")
            return

        except Exception as err:
            if attempts == 0:
                raise TimeoutError(f"{container.name} not healthy: {err}")

            else:
                attempts -= 1
                time.sleep(interval)